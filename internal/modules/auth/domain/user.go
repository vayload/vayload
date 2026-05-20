/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import (
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

// User represents the core entity in our domain, matching the 'users' table in sqlite.sql.
type User struct {
	ID           snowflake.ID   `json:"id,string"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	Phone        *string        `json:"phone,omitempty"`
	PasswordHash *string        `json:"-"`
	FirstName    *string        `json:"first_name,omitempty"`
	LastName     *string        `json:"last_name,omitempty"`
	AvatarURL    *string        `json:"avatar_url,omitempty"`
	IsSuperAdmin bool           `json:"is_super_admin"`
	IsSSOUser    bool           `json:"is_sso_user"`
	Role         string         `json:"role,omitempty"`
	CountryID    snowflake.ID   `json:"country_id,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	Settings     map[string]any `json:"settings,omitempty"`
	Attributes   map[string]any `json:"attributes,omitempty"`

	// Verification and Security
	EmailConfirmedAt *time.Time `json:"email_confirmed_at,omitempty"`
	PhoneConfirmedAt *time.Time `json:"phone_confirmed_at,omitempty"`
	ConfirmedAt      *time.Time `json:"confirmed_at,omitempty"`

	ConfirmationToken *string    `json:"-"`
	RecoveryToken     *string    `json:"-"`
	EmailChangeToken  *string    `json:"-"`
	PhoneChangeToken  *string    `json:"-"`
	OTPCode           *string    `json:"-"`
	EmailChange       *string    `json:"-"`
	PhoneChange       *string    `json:"-"`
	BannedUntil       *time.Time `json:"banned_until,omitempty"`

	ConfirmationSentAt *time.Time `json:"-"`
	RecoverySentAt     *time.Time `json:"-"`
	EmailChangeSentAt  *time.Time `json:"-"`
	PhoneChangeSentAt  *time.Time `json:"-"`
	OTPSentAt          *time.Time `json:"-"`
	LastSignInAt       *time.Time `json:"last_sign_in_at,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
}

func NewUser(username string, email string, password *string) *User {
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}
}

func (u *User) SetPassword(password *string) {
	u.PasswordHash = password
}

func (u *User) UnsetPassword() {
	u.PasswordHash = nil
}

// IsBanned returns true if the account is currently banned.
func (u *User) IsBanned() bool {
	if u.BannedUntil == nil {
		return false
	}
	return u.BannedUntil.After(time.Now().UTC())
}

// IsEmailVerified returns true if the email has been verified.
func (u *User) IsEmailVerified() bool {
	return u.EmailConfirmedAt != nil
}

// IsPhoneVerified returns true if the phone has been verified.
func (u *User) IsPhoneVerified() bool {
	return u.PhoneConfirmedAt != nil
}

// IsAccountActive returns true if the account is active (not deleted and not banned).
func (u *User) IsAccountActive() bool {
	return u.DeletedAt == nil && !u.IsBanned()
}

// SignedToken represents a signed JWT token.
type SignedToken interface {
	GetAccessToken() string
	GetExpiresAt() time.Time
	GetExpiresIn() int64
	GetPayload() any
}

// SignedTokenWithRefresh represents a signed JWT token with a refresh token.
type SignedTokenWithRefresh interface {
	SignedToken
	GetRefreshToken() string
	GetRefreshTokenExpiresAt() time.Time
	GetRefreshTokenExpiresIn() int64
	GetMeta() map[string]any
}

// OAuthSession represents an active authentication session.
type OAuthSession struct {
	AccessToken      string         `json:"access_token"`
	TokenType        string         `json:"token_type"`
	RefreshToken     string         `json:"refresh_token,omitempty"`
	ExpiresIn        int64          `json:"expires_in"`
	ExpiresAt        time.Time      `json:"expires_at"`
	User             *User          `json:"user,omitempty"`
	Meta             map[string]any `json:"meta,omitempty"`
	ExpiresRefreshAt time.Time      `json:"expires_refresh_at,omitempty"`
	ExpiresRefreshIn int64          `json:"expires_refresh_in,omitempty"`
}

func (o *OAuthSession) ToJson() map[string]any {
	data := map[string]any{
		"access_token": o.AccessToken,
		"token_type":   o.TokenType,
		"expires_in":   o.ExpiresIn,
		"expires_at":   o.ExpiresAt.Unix(),
	}

	if o.RefreshToken != "" {
		data["refresh_token"] = o.RefreshToken
	}
	if o.User != nil {
		data["user"] = o.User
	}
	if o.Meta != nil {
		data["meta"] = o.Meta
	}
	return data
}

func NewOAuthSession(token SignedTokenWithRefresh, user *User, meta map[string]any) *OAuthSession {
	return &OAuthSession{
		AccessToken:      token.GetAccessToken(),
		TokenType:        "Bearer",
		RefreshToken:     token.GetRefreshToken(),
		ExpiresIn:        token.GetExpiresIn(),
		ExpiresAt:        token.GetExpiresAt(),
		User:             user,
		Meta:             meta,
		ExpiresRefreshAt: token.GetRefreshTokenExpiresAt(),
		ExpiresRefreshIn: token.GetRefreshTokenExpiresIn(),
	}
}

// AuthContext represents the context of an authentication request.
type AuthContext struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Method    string `json:"method"`
	SessionID string `json:"sid,omitempty"`
}

// IdentifierType represents the type of user identifier used for login.
type IdentifierType string

const (
	IdentifierTypeEmail    IdentifierType = "email"
	IdentifierTypePhone    IdentifierType = "phone"
	IdentifierTypeUsername IdentifierType = "username"
	IdentifierTypeID       IdentifierType = "id"
)

type OtpProvider struct {
	Email    []string `json:"email"`
	SMS      []string `json:"sms"`
	WhatsApp []string `json:"whatsapp"`
}

func (i IdentifierType) Valid() bool {
	switch i {
	case IdentifierTypeEmail, IdentifierTypePhone, IdentifierTypeUsername, IdentifierTypeID:
		return true
	}
	return false
}
