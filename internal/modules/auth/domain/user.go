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

type UserRole string

const (
	AdminRole     UserRole = "admin"
	PatientRole   UserRole = "patient"
	CaregiverRole UserRole = "caregiver"
	DoctorRole    UserRole = "doctor"
)

// User represents the core entity in our domain.
type User struct {
	ID           snowflake.ID   `json:"id,string"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	Phone        *string        `json:"phone"`                // Optional phone number
	Password     *string        `json:"-"`                    // Password is hashed, not returned in responses
	Role         UserRole       `json:"role"`                 // User role (admin, patient, etc.)
	LastName     *string        `json:"last_name"`            // Optional last name
	Meta         map[string]any `json:"meta,omitempty"`       // For additional user metadata
	ClientId     *int64         `json:"client_id,omitempty"`  // For policies and authorization relationships
	ProfileId    *int64         `json:"profile_id,omitempty"` // For user profile relationships (depends of clientId)
	CountryID    *snowflake.ID  `json:"country_id,omitempty"` // For user-country relationship
	AuthType     string         `json:"-"`                    // The type of authentication method (e.g., password, oauth, google, etc.)
	AvatarURL    *string        `json:"avatar_url,omitempty"` // Optional avatar URL
	DigevoCoreId *int64         `json:"-"`                    // Digevo core ID for VAS users
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`

	// Account status fields (not exposed in JSON)
	EmailVerified bool       `json:"-"` // Email verification status
	PhoneVerified bool       `json:"-"` // Phone verification status
	BannedUntil   *time.Time `json:"-"` // Ban expiration date
	DeletedAt     *time.Time `json:"-"` // Soft delete timestamp
	OTPCode       string     `json:"-"` // OTP code for authentication
}

func (user *User) SetPassword(password *string) {
	user.Password = password
}

// IsDeleted returns true if the account has been deleted/disabled by admin
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsBanned returns true if the account is currently banned
func (u *User) IsBanned() bool {
	if u.BannedUntil == nil {
		return false
	}

	return u.BannedUntil.After(time.Now().UTC())
}

// GetBannedUntilFormatted returns the ban expiration date formatted as RFC3339
func (u *User) GetBannedUntilFormatted() string {
	if u.BannedUntil == nil {
		return ""
	}

	return u.BannedUntil.Format(time.RFC3339)
}

// IsEmailVerified returns true if the email has been verified
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}

// IsPhoneVerified returns true if the phone has been verified
func (u *User) IsPhoneVerified() bool {
	return u.PhoneVerified
}

// IsAccountActive returns true if the account is active (not deleted and not banned)
func (u *User) IsAccountActive() bool {
	return !u.IsDeleted() && !u.IsBanned()
}

type UserAuthCodes struct {
	Code      string    `json:"code"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserAuth struct {
	ID                    snowflake.ID
	Role                  UserRole
	Email                 string
	Phone                 *string
	PasswordHash          string
	EmailConfirmed        bool
	PhoneConfirmed        bool
	ConfirmationToken     string
	RecoveryToken         string
	OTPCode               string
	BannedUntil           *time.Time
	DeletedAt             *time.Time
	ReauthenticationToken string
	AuthType              string
	AuthSID               string
	IsSuperAdmin          bool
	IsSSOUser             bool
	LastSignInAt          *time.Time
}

// IsDeleted returns true if the account has been deleted/disabled by admin
func (u *UserAuth) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsBanned returns true if the account is currently banned
func (u *UserAuth) IsBanned() bool {
	if u.BannedUntil == nil {
		return false
	}

	return u.BannedUntil.After(time.Now().UTC())
}

// GetBannedUntilFormatted returns the ban expiration date formatted as RFC3339
func (u *UserAuth) GetBannedUntilFormatted() string {
	if u.BannedUntil == nil {
		return ""
	}

	return u.BannedUntil.Format(time.RFC3339)
}

// IsEmailVerified returns true if the email has been verified
func (u *UserAuth) IsEmailVerified() bool {
	return u.EmailConfirmed
}

// IsPhoneVerified returns true if the phone has been verified
func (u *UserAuth) IsPhoneVerified() bool {
	return u.PhoneConfirmed
}

// IsAccountActive returns true if the account is active (not deleted and not banned)
func (u *UserAuth) IsAccountActive() bool {
	return !u.IsDeleted() && !u.IsBanned()
}

type UserNotificationSettings struct {
	Email    bool
	Push     bool
	SMS      bool
	WhatsApp bool
}
type UserSettings struct {
	Language      string
	Notifications UserNotificationSettings
}

type UserIdentity struct {
	UserID       snowflake.ID `json:"user_id"`
	Email        string       `json:"email"`         // Email is optional, used for login
	Phone        *string      `json:"phone"`         // Phone is optional, used for login
	CountryID    *int64       `json:"country_id"`    // Optional country ID for OTP providers
	DocumentID   *string      `json:"document_id"`   // Unique identifier for the user (DNI, passport, etc.)
	DocumentType *string      `json:"document_type"` // Type of document (DNI, passport, etc.)
}

type UserMeta struct {
	ConfirmationToken string `json:"confirmation_token,omitempty"` // Token for email confirmation
	VerificationCode  string `json:"verification_code,omitempty"`  // Verification code for email/phone
	EmailVerified     bool   `json:"email_verified,omitempty"`     // Email verification status
	PhoneVerified     bool   `json:"phone_verified,omitempty"`     // Phone verification status
	VasProvider       string `json:"vas_provider,omitempty"`       // VAS provider identifier
	DigevoCoreId      int64  `json:"digevo_core_id,omitempty"`     // Digevo core ID
	SubscriptionId    string `json:"subscription_id,omitempty"`    // Digevo subscription ID
}

type AuthUser struct {
	ID        snowflake.ID   `json:"id,omitempty"`
	Email     string         `json:"email"`
	Role      UserRole       `json:"role"`
	ClientId  int64          `json:"client_id,omitempty"`  // For user policies relationship
	CountryId *snowflake.ID  `json:"country_id,omitempty"` // For user-country relationship
	Meta      map[string]any `json:"meta,omitempty"`       // For additional user metadata
}

type UserMetadata struct {
	User
	OtpCode           string `json:"otp_code,omitempty"`           // OTP code for 2FA
	VerificationCode  string `json:"verification_code,omitempty"`  // Verification code for email/phone
	ConfirmationToken string `json:"confirmation_token,omitempty"` // Token for email confirmation
}

type OAuthSession struct {
	AccessToken      string
	TokenType        string
	RefreshToken     string
	ExpiresIn        int64
	ExpiresAt        time.Time
	User             *User
	Policies         *UserPolicy // Optional, for user policies
	Meta             map[string]any
	ExpiresRefreshAt time.Time // Optional, for refresh token expiration
	ExpiresRefreshIn int64     // Optional, for refresh token expiration in seconds
}

func (o *OAuthSession) GetAccessToken() string {
	return o.AccessToken
}

func (o *OAuthSession) GetRefreshToken() string {
	return o.RefreshToken
}

func (o *OAuthSession) ToJson() map[string]any {
	data := map[string]any{
		"access_token":  o.AccessToken,
		"token_type":    o.TokenType,
		"refresh_token": o.RefreshToken,
		"expires_in":    o.ExpiresIn,
		"expires_at":    o.ExpiresAt.Unix(),
	}

	if o.User != nil {
		data["user"] = o.User
	}
	if o.Meta != nil {
		data["meta"] = o.Meta
	}
	if o.Policies != nil {
		data["policies"] = *o.Policies
	}

	return data
}

func (o *OAuthSession) FromSignedJwt(token SignedTokenWithRefresh, data ...any) *OAuthSession {
	user := &User{}
	meta := make(map[string]any)
	if len(data) > 0 {
		if u, ok := data[0].(*User); ok {
			user = u
		}
	}
	if len(data) > 1 {
		if m, ok := data[1].(map[string]any); ok {
			meta = m
		}
	}

	return &OAuthSession{
		AccessToken:      token.GetAccessToken(),
		TokenType:        "Bearer",
		RefreshToken:     token.GetRefreshToken(),
		ExpiresIn:        token.GetExpiresIn(),
		ExpiresAt:        token.GetExpiresAt(),
		User:             user,
		Meta:             meta,
		ExpiresRefreshAt: token.GetRefreshTokenExpiresAt(),
		ExpiresRefreshIn: token.GetRefreshTokenExpiresIn(), // Optional, for refresh token expiration
	}
}

func NewOAuthSessionFromToken(token SignedTokenWithRefresh) *OAuthSession {
	return &OAuthSession{
		AccessToken:      token.GetAccessToken(),
		TokenType:        "Bearer",
		RefreshToken:     token.GetRefreshToken(),
		ExpiresIn:        token.GetExpiresIn(),
		ExpiresAt:        token.GetExpiresAt(),
		Meta:             token.GetMeta(),
		ExpiresRefreshAt: token.GetRefreshTokenExpiresAt(),
		ExpiresRefreshIn: token.GetRefreshTokenExpiresIn(), // Optional, for refresh token expiration
	}
}

func NewOAuthSession(token SignedTokenWithRefresh, user *User, policies *UserPolicy, meta map[string]any) *OAuthSession {
	return &OAuthSession{
		AccessToken:      token.GetAccessToken(),
		TokenType:        "Bearer",
		RefreshToken:     token.GetRefreshToken(),
		ExpiresIn:        token.GetExpiresIn(),
		ExpiresAt:        token.GetExpiresAt(),
		User:             user,
		Policies:         policies,
		Meta:             meta,
		ExpiresRefreshAt: token.GetRefreshTokenExpiresAt(),
		ExpiresRefreshIn: token.GetRefreshTokenExpiresIn(), // Optional, for refresh token
	}
}

type OtpProvider struct {
	Email    []string `json:"email,omitempty"`    // Email OTP provider
	SMS      []string `json:"sms,omitempty"`      // SMS OTP provider
	WhatsApp []string `json:"whatsapp,omitempty"` // WhatsApp OTP provider
}

type AuthContext struct {
	IP        string  `json:"ip"`
	UserAgent string  `json:"user_agent"`
	Method    string  `json:"method"`
	Sid       *string `json:"sid,omitempty"` // Session ID for tracking
}

type VerificationCodes struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	Code      string    `json:"code"`
	Type      string    `json:"type"` // e.g. "email", "phone"
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"` // Expiration time for the
	Attempts  int       `json:"attempts"`   // Number of attempts made

	User *User
}

type UserVerification struct {
	ID                       snowflake.ID
	Email                    string
	EmailConfirmedAt         *time.Time
	VerificationCode         *string
	ConfirmationToken        *string
	ConfirmationSentAt       *time.Time
	RecoveryToken            *string
	RecoverySentAt           *time.Time
	EmailChangeTokenNew      *string
	EmailChange              *string
	EmailChangeSentAt        *time.Time
	EmailChangeTokenCurrent  *string
	EmailChangeConfirmStatus *int32
	Phone                    *string
	PhoneConfirmedAt         *time.Time
	PhoneChange              *string
	PhoneChangeToken         *string
	PhoneChangeSentAt        *time.Time
	OTPCode                  *string
	OTPSentAt                *time.Time
	ReauthenticationToken    *string
	ReauthenticationSentAt   *time.Time
	ConfirmedAt              *time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                *time.Time
}

type UserWithVerificationCodes struct {
	User
	Code      string    `json:"code"`
	Type      string    `json:"type"`       // e.g. "email", "phone"
	ExpiresAt time.Time `json:"expires_at"` // Expiration time for the
	Attempts  int       `json:"attempts"`   // Number of attempts made
}

func NewUser(username, email string, password *string, role UserRole) *User {
	return &User{
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type IdentifierType string

const (
	IdentifierTypeEmail    IdentifierType = "email"
	IdentifierTypePhone    IdentifierType = "phone"
	IdentifierTypeUsername IdentifierType = "username"
	IdentifierUserId       IdentifierType = "user_id"
)

func (i IdentifierType) Valid() bool {
	switch i {
	case IdentifierTypeEmail, IdentifierTypePhone, IdentifierTypeUsername, IdentifierUserId:
		return true
	}

	return false
}
