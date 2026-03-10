/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Security
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package security

import (
	"crypto"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/shared/snowflake"
	"github.com/vayload/vayload/pkg/logger"
)

type JwtConfig struct {
	// Used for HMAC signing (refresh tokens)
	SecretKey string
	// Used for RSA signing (access tokens)
	PrivateKeyBytes []byte
	// Used for RSA verification (access tokens)
	PublicKeyBytes []byte
	// Expiration times in hours for access tokens
	ExpireAccessToken int64
	// Expiration time in hours for refresh tokens
	ExpireRefreshToken int64
}

type JwtManager struct {
	// Secret key for HMAC signing
	SecretKey string
	// RSA keys for signing and verification
	privateKey crypto.PrivateKey
	// Public key for verifying JWT tokens
	publicKey crypto.PublicKey
	// Expiration times in minutes for access tokens
	ExpireAccessToken int64
	// Expiration time in hours for refresh tokens
	ExpireRefreshToken int64
}

const JWT_MANAGER_KEY = "jwtManager"

func NewJwtManager(config JwtConfig) *JwtManager {
	// Parse Private Key (expect PKCS#8)
	privKey, err := x509.ParsePKCS8PrivateKey(config.PrivateKeyBytes)
	if err != nil {
		logger.F(err, logger.Fields{"context": "NewJwtManager", "action": "parse private key"})
	}
	edPriv, ok := privKey.(ed25519.PrivateKey)
	if !ok {
		logger.F(err, logger.Fields{"context": "NewJwtManager", "action": "parse private key"})
	}

	// Parse Public Key (expect PKIX)
	pubKey, err := x509.ParsePKIXPublicKey(config.PublicKeyBytes)
	if err != nil {
		logger.F(err, logger.Fields{"context": "NewJwtManager", "action": "parse public key"})
	}
	edPub, ok := pubKey.(ed25519.PublicKey)
	if !ok {
		logger.F(err, logger.Fields{"context": "NewJwtManager", "action": "parse public key"})
	}

	return &JwtManager{
		privateKey:         edPriv,
		publicKey:          edPub,
		ExpireAccessToken:  config.ExpireAccessToken,
		ExpireRefreshToken: config.ExpireRefreshToken,
		SecretKey:          config.SecretKey,
	}
}

func (tk *JwtManager) GenerateJwtTokenWithRefresh(user *domain.AuthUser) (domain.SignedTokenWithRefresh, error) {
	token, err := tk.GenerateJwtToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken := tk.CreateRefreshToken(fmt.Sprintf("%d_%d_%s", user.ID, user.ClientId, user.Email))

	return &jwtToken{
		AccessToken:      token.GetAccessToken(),
		RefreshToken:     refreshToken,
		ExpiresAt:        token.GetExpiresAt(),
		IssuedAt:         time.Now().UTC(),
		Payload:          token.GetPayload(),
		Meta:             user.Meta,
		ExpiresIn:        token.GetExpiresIn(),
		RefreshExpiresAt: time.Now().UTC().Add(time.Duration(tk.ExpireRefreshToken) * time.Hour),
		RefreshExpiresIn: tk.ExpireRefreshToken * 3600, // Convert hours to seconds
	}, nil
}

func (tk *JwtManager) GenerateJwtToken(user *domain.AuthUser) (domain.SignedToken, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(tk.ExpireAccessToken) * time.Minute)
	payload := jwt.MapClaims{
		"sub":        user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"client_id":  user.ClientId,
		"country_id": user.CountryId,
		"iat":        time.Now().UTC().Unix(),
		"exp":        expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, payload)
	signedToken, err := token.SignedString(tk.privateKey)
	if err != nil {
		return nil, err
	}

	return &jwtToken{
		AccessToken: signedToken,
		ExpiresAt:   expiresAt,
		IssuedAt:    time.Now().UTC(),
		Payload:     payload,
		Meta:        user.Meta,
		ExpiresIn:   tk.ExpireAccessToken * 60, // Convert minutes to seconds
	}, nil
}

func (tk *JwtManager) ValidateToken(accessToken string) (*domain.AuthUser, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		return tk.publicKey, nil
	})
	if err != nil {
		logger.E(err, logger.Fields{"context": "ValidateToken", "action": "parse token"})
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return parseClaims(claims)
}

func (tk *JwtManager) CreateRefreshToken(payload string) string {
	expire := time.Now().UTC().Add(time.Duration(tk.ExpireRefreshToken) * time.Hour).Unix()
	signature := tk.signHMAC(payload)
	signed := fmt.Sprintf("%s|%d|%s", payload, expire, signature)

	return base64.RawURLEncoding.EncodeToString([]byte(signed))

}

func (tk *JwtManager) ValidateRefreshToken(tokenString string) (*domain.AuthUser, error) {
	tokenBytes, err := base64.RawURLEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(tokenBytes), "|")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	data, expireStr, sig := parts[0], parts[1], parts[2]
	expectedSig := tk.signHMAC(data)

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return nil, errors.New("invalid token signature")
	}

	exp, err := strconv.Atoi(expireStr)
	if err != nil || time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	partsData := strings.Split(data, "_")
	if len(partsData) != 3 {
		return nil, errors.New("invalid token payload")
	}

	userID, err := strconv.ParseInt(partsData[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	role := partsData[1]

	email := partsData[2]

	return &domain.AuthUser{
		ID:    snowflake.ID(userID),
		Email: email,
		Role:  domain.UserRole(role),
	}, nil
}

func (tk *JwtManager) signHMAC(payload string) string {
	h := hmac.New(sha256.New, []byte(tk.SecretKey))
	h.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func parseClaims(claims jwt.MapClaims) (*domain.AuthUser, error) {
	id, err := toInt64(claims["sub"])
	if err != nil {
		return nil, fmt.Errorf("invalid id claim: %w", err)
	}

	role, exists := claims["role"].(string)
	if !exists {
		return nil, fmt.Errorf("invalid role claim: %w", err)
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email claim")
	}

	clientId, err := toInt64(claims["client_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid client_id claim: %w", err)
	}

	countryId, err := toInt64(claims["country_id"])
	if err != nil {
		logger.W("invalid country_id claim", logger.Fields{"error": err})
	}

	var meta map[string]any
	if metaValue, exists := claims["meta"]; exists && metaValue != nil {
		var err error
		meta, err = parseMeta(metaValue)
		if err != nil {
			logger.W("invalid meta claim", logger.Fields{"error": err})
			// Don't return error, just use nil meta
			meta = nil
		}
	}

	countryID := snowflake.ID(countryId)
	return &domain.AuthUser{
		ID:        snowflake.ID(id),
		Email:     email,
		Role:      domain.UserRole(role),
		ClientId:  int64(clientId),
		CountryId: &countryID,
		Meta:      meta,
	}, nil
}

func toInt64(value any) (int64, error) {
	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unexpected type %T", v)
	}
}

func parseMeta(value any) (map[string]any, error) {
	if value == nil {
		return nil, nil
	}

	switch v := value.(type) {
	case map[string]any:
		return v, nil
	default:
		return nil, fmt.Errorf("unexpected type %T", v)
	}
}

type jwtToken struct {
	AccessToken      string
	RefreshToken     string
	ExpiresAt        time.Time
	IssuedAt         time.Time
	Payload          any
	Meta             map[string]any
	ExpiresIn        int64
	RefreshExpiresAt time.Time
	RefreshExpiresIn int64
}

func (t *jwtToken) GetPayload() any {
	return t.Payload
}

func (t *jwtToken) GetMeta() map[string]any {
	return t.Meta
}

func (t *jwtToken) GetAccessToken() string {
	return t.AccessToken
}

func (t *jwtToken) GetRefreshToken() string {
	return t.RefreshToken
}

func (t *jwtToken) GetExpiresAt() time.Time {
	return t.ExpiresAt
}

// GetExpiresIn returns the expiration time in seconds.
func (t *jwtToken) GetExpiresIn() int64 {
	return t.ExpiresIn
}

func (t *jwtToken) GetRefreshTokenExpiresAt() time.Time {
	return t.RefreshExpiresAt
}

// GetRefreshTokenExpiresIn returns the refresh token expiration time in seconds.
func (t *jwtToken) GetRefreshTokenExpiresIn() int64 {
	return t.RefreshExpiresIn
}
