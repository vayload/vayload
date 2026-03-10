/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package providers

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type oauth2Facade struct {
	providers sync.Map
	config    *config.Config
}

func NewOAuth2Facade(config *config.Config) *oauth2Facade {
	facade := &oauth2Facade{
		providers: sync.Map{},
		config:    config,
	}

	if config.OAuth.Google.ClientID != "" && config.OAuth.Google.ClientSecret != "" {
		googleOAuth2 := NewGoogleOAuth2Strategy(
			config.OAuth.Google.ClientID,
			config.OAuth.Google.ClientSecret,
			config.OAuth.RedirectBase+"/v1/auth/google/callback",
		)

		facade.providers.Store(domain.OAuth2Google, googleOAuth2)
	}

	if config.OAuth.Facebook.ClientID != "" && config.OAuth.Facebook.ClientSecret != "" {
		facebookOAuth2 := NewFacebookOAuth2Strategy(
			config.OAuth.Facebook.ClientID,
			config.OAuth.Facebook.ClientSecret,
			config.OAuth.RedirectBase+"/v1/auth/facebook/callback",
		)

		facade.providers.Store(domain.OAuth2Facebook, facebookOAuth2)
	}

	return facade
}

func (f *oauth2Facade) Select(provider domain.OAuth2Provider) (domain.OAuth2Strategy, error) {
	strategy, ok := f.providers.Load(provider)
	if !ok {
		return nil, fmt.Errorf("OAuth2 provider not found: %s", provider)
	}

	return strategy.(domain.OAuth2Strategy), nil
}

func (f *oauth2Facade) GetAuthRedirectURL(provider domain.OAuth2Provider, state *domain.OAuth2State) (string, error) {
	strategy, err := f.Select(provider)
	if err != nil {
		return "", err
	}

	return strategy.GetAuthURL(state), nil
}

func (f *oauth2Facade) ExchangeCode(provider domain.OAuth2Provider, code string) (*domain.OAuthUser, error) {
	strategy, err := f.Select(provider)
	if err != nil {
		return nil, err
	}

	return strategy.ExchangeCode(code)
}

func (f *oauth2Facade) ValidateToken(provider domain.OAuth2Provider, token string) (*domain.OAuthUser, error) {
	strategy, err := f.Select(provider)
	if err != nil {
		return nil, err
	}

	return strategy.ValidateToken(token)
}

type GoogleOAuth2Strategy struct {
	config *oauth2.Config
}

func NewGoogleOAuth2Strategy(clientID, clientSecret, redirectURL string) *GoogleOAuth2Strategy {
	return &GoogleOAuth2Strategy{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *GoogleOAuth2Strategy) GetAuthURL(state *domain.OAuth2State) string {
	data, _ := json.Marshal(state)
	stateEncoded := base64.URLEncoding.EncodeToString(data)

	return g.config.AuthCodeURL(stateEncoded)
}

func (g *GoogleOAuth2Strategy) ExchangeCode(code string) (*domain.OAuthUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google token exchange failed: %w", err)
	}

	// Get user info from Google
	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}
	defer resp.Body.Close()

	var userData map[string]any
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Google user info response: %w", err)
	}

	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, fmt.Errorf("failed to parse Google user info: %w", err)
	}

	return &domain.OAuthUser{
		SID:           toString(userData["sub"]),
		Email:         toString(userData["email"]),
		EmailVerified: toBool(userData["email_verified"]),
		Name:          toString(userData["name"]),
		FirstName:     toString(userData["given_name"]),
		LastName:      toString(userData["family_name"]),
		AvatarURL:     toString(userData["picture"]),
	}, nil
}

func (g *GoogleOAuth2Strategy) ValidateToken(token string) (*domain.OAuthUser, error) {
	return nil, fmt.Errorf("use ExchangeCode for OAuth2 flow")
}

type FacebookOAuth2Strategy struct {
	config *oauth2.Config
}

func NewFacebookOAuth2Strategy(clientID, clientSecret, redirectURL string) *FacebookOAuth2Strategy {
	return &FacebookOAuth2Strategy{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "public_profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/v17.0/dialog/oauth",
				TokenURL: "https://graph.facebook.com/v17.0/oauth/access_token",
			},
		},
	}
}

func (f *FacebookOAuth2Strategy) GetAuthURL(state *domain.OAuth2State) string {
	data, _ := json.Marshal(state)
	stateEncoded := base64.URLEncoding.EncodeToString(data)

	return f.config.AuthCodeURL(stateEncoded)
}

func (f *FacebookOAuth2Strategy) ExchangeCode(code string) (*domain.OAuthUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := f.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("facebook token exchange failed: %w", err)
	}

	// Get user info from Facebook
	client := f.config.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, fmt.Errorf("failed to get Facebook user info: %w", err)
	}
	defer resp.Body.Close()

	var userData map[string]any
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Facebook user info response: %w", err)
	}

	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, fmt.Errorf("failed to parse Facebook user info: %w", err)
	}

	// Extract picture URL if available
	pictureURL := ""
	if picture, ok := userData["picture"].(map[string]any); ok {
		if data, ok := picture["data"].(map[string]any); ok {
			pictureURL = toString(data["url"])
		}
	}

	return &domain.OAuthUser{
		SID:           toString(userData["id"]),
		Email:         toString(userData["email"]),
		EmailVerified: true, // Facebook emails are typically verified
		Name:          toString(userData["name"]),
		FirstName:     "",
		LastName:      "",
		AvatarURL:     pictureURL,
	}, nil
}

func (f *FacebookOAuth2Strategy) ValidateToken(token string) (*domain.OAuthUser, error) {
	// This method can be used for direct token validation if needed
	// For OAuth2 flow, use ExchangeCode instead
	return nil, fmt.Errorf("use ExchangeCode for OAuth2 flow")
}

type AppleOAuth2Strategy struct {
	clientID    string
	teamID      string
	keyID       string
	keyPEM      string
	redirectURL string
}

func NewAppleOAuth2Strategy(clientID, teamID, keyID, keyPEM, redirectURL string) *AppleOAuth2Strategy {
	return &AppleOAuth2Strategy{
		clientID:    clientID,
		teamID:      teamID,
		keyID:       keyID,
		keyPEM:      keyPEM,
		redirectURL: redirectURL,
	}
}

func (a *AppleOAuth2Strategy) GetAuthURL(state *domain.OAuth2State) string {
	data, _ := json.Marshal(state)
	stateEncoded := base64.URLEncoding.EncodeToString(data)

	params := url.Values{
		"client_id":     {a.clientID},
		"redirect_uri":  {a.redirectURL},
		"response_type": {"code"},
		"response_mode": {"form_post"},
		"scope":         {"name email"},
		"state":         {stateEncoded},
	}

	return fmt.Sprintf("https://appleid.apple.com/auth/authorize?%s", params.Encode())
}

func (a *AppleOAuth2Strategy) ExchangeCode(code string) (*domain.OAuthUser, error) {
	// Generate Apple client secret
	clientSecret, err := a.generateClientSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Apple client secret: %w", err)
	}

	// Exchange code for tokens
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", a.redirectURL)
	form.Set("client_id", a.clientID)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", "https://appleid.apple.com/auth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create Apple token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apple token request failed: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp map[string]any
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Apple token response: %w", err)
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse Apple token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse ID token to get user info
	idToken := toString(tokenResp["id_token"])
	if idToken == "" {
		return nil, fmt.Errorf("no id_token received from Apple")
	}

	// Parse JWT claims (skip signature verification for now)
	claims := jwt.MapClaims{}
	_, _, err = jwt.NewParser().ParseUnverified(idToken, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Apple ID token: %w", err)
	}

	return &domain.OAuthUser{
		SID:           toString(claims["sub"]),
		Email:         toString(claims["email"]),
		EmailVerified: toBool(claims["email_verified"]),
		Name:          toString(claims["name"]),
		FirstName:     "",
		LastName:      "",
		AvatarURL:     "",
	}, nil
}

// generateClientSecret creates a JWT for Apple client secret
func (a *AppleOAuth2Strategy) generateClientSecret() (string, error) {
	block, _ := pem.Decode([]byte(a.keyPEM))
	if block == nil {
		return "", fmt.Errorf("invalid Apple key PEM")
	}

	var key any
	var err error

	// Try different key formats
	key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		key, err = x509.ParseECPrivateKey(block.Bytes)
	}
	if err != nil {
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	}
	if err != nil {
		return "", fmt.Errorf("failed to parse Apple private key: %w", err)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": a.teamID,
		"iat": now.Unix(),
		"exp": now.Add(time.Hour * 6).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": a.clientID,
	}

	var token *jwt.Token

	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		token = jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	case *rsa.PrivateKey:
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	default:
		return "", fmt.Errorf("unsupported key type: %T", k)
	}

	token.Header["kid"] = a.keyID
	return token.SignedString(key)
}

// ValidateToken implements OAuthStrategy interface for backward compatibility
func (a *AppleOAuth2Strategy) ValidateToken(token string) (*domain.OAuthUser, error) {
	// This method can be used for direct token validation if needed
	// For OAuth2 flow, use ExchangeCode instead
	return nil, fmt.Errorf("use ExchangeCode for OAuth2 flow")
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toBool(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
