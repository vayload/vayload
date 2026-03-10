/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

// =================================== password strategy ===================================
type PasswordStrategy interface {
	HashPassword(password string) string
	VerifyPassword(password, hashedPassword string) bool
	VerboseVerifyPassword(password, hashedPassword string) (valid bool, algo string)
}

const PasswordStrategyType = "password"

// =================================== magic link strategy ===================================

// For this strategy use a encrypted payload in the token
type MagicLinkStrategy interface {
	// Sign generates a magic link token.
	Sign(payload *MagicLinkPayload) (string, error)

	// Verify checks the validity of a magic link token.
	Verify(token string) (*MagicLinkPayload, error)

	// GetExpirationTime returns the expiration time of the magic link token.
	GetExpirationTime() time.Duration
}

type MagicLinkPayload struct {
	Redirect       string `json:"r"`            // URL to redirect after login
	UserIdentifier string `json:"ui"`           // Email or phone
	IdentifierType string `json:"it"`           // e.g. "email", "phone"
	Meta           any    `json:"mt,omitempty"` // Additional metadata
}

// =================================== otp strategy ===================================
type OtpStrategy interface {
	GenerateOtpCode() string
	CompareOtpCode(inputCode string, storedCode string) bool
}

const OtpStrategyType = "otp"

// =================================== oauth strategy ===================================

type OAuth2Provider string

const (
	OAuth2Google   OAuth2Provider = "google"
	OAuth2Facebook OAuth2Provider = "facebook"
	OAuth2Twitter  OAuth2Provider = "twitter"
	OAuth2LinkedIn OAuth2Provider = "linkedin"
	OAuth2Apple    OAuth2Provider = "apple"
)

func (p OAuth2Provider) IsSupported() bool {
	switch p {
	case OAuth2Google, OAuth2Facebook, OAuth2Twitter, OAuth2LinkedIn, OAuth2Apple:
		return true
	}

	return false
}

type OAuthStrategy interface {
	ValidateToken(token string) (*OAuthUser, error)
}

type OAuth2Strategy interface {
	OAuthStrategy
	GetAuthURL(state *OAuth2State) string
	ExchangeCode(code string) (*OAuthUser, error)
}

type OAuth2StrategyFacade interface {
	Select(provider OAuth2Provider) (OAuth2Strategy, error)
	GetAuthRedirectURL(provider OAuth2Provider, state *OAuth2State) (string, error)
	ExchangeCode(provider OAuth2Provider, code string) (*OAuthUser, error)
	ValidateToken(provider OAuth2Provider, token string) (*OAuthUser, error)
}

type OAuthUser struct {
	SID           string `json:"sid"`            // from "sub"
	Email         string `json:"email"`          // from "email"
	EmailVerified bool   `json:"email_verified"` // from "email_verified"
	Name          string `json:"name"`           // from "name"
	FirstName     string `json:"first_name"`     // from "given_name"
	LastName      string `json:"last_name"`      // from "family_name"
	AvatarURL     string `json:"avatar_url"`     // from "picture"
}

type OAuth2State struct {
	Origin string `json:"origin"`
	Nonce  string `json:"nonce"`
}

// =================================== telco strategy ===================================
type TelcoAuthStrategy interface {
	Name() string
	GetCustody(ctx context.Context, code string, productId int) (*TelcoCustody, error)

	GetClientById(ctx context.Context, clientId int) (*TelcoClientUser, error)
	GetClientByPhone(ctx context.Context, phone string) (*TelcoClientUser, error)
	GetStatus(ctx context.Context, id int, productId int) (*TelcoStatus, error)
	GetUserById(ctx context.Context, userId int, productId int) (*TelcoFremiumUser, error)
	SetProfileTo(ctx context.Context, clientId int, productId int, profileId int) (*TelcoFremiumUser, error)
}

type TelcoClientProfile struct {
	Email          string       `json:"email"`
	Phone          string       `json:"phone"`
	Name           string       `json:"name"`
	ProfileId      int64        `json:"profile_id"`
	CountryId      snowflake.ID `json:"country_id"`
	VasProvider    string       `json:"vas_provider"`
	SubscriptionId string       `json:"subscription_id"`
	DigevoCoreId   int64        `json:"digevo_core_id"`
	IsFree         bool         `json:"is_free"`
}

type TelcoCustody struct {
	ProfileId string `json:"profile_id"`
	ClientId  int64  `json:"client_id"`
	Product   string `json:"product"`
	IdList    int    `json:"idlist"`
	Origin    string `json:"origin"`
}

type TelcoStatus struct {
	Status bool `json:"status"`
}

type TelcoClientUser struct {
	ClientId         int64  `json:"client_id"`
	Phone            string `json:"phone"`
	OperatorId       int    `json:"operator_id"`
	User             string `json:"user"`
	Email            string `json:"email"`
	Pass             string `json:"pass"`
	UserIdentifyData string `json:"user_identify_data"`
}

type TelcoLocation struct {
	CountryID *snowflake.ID `json:"country_id"`
	IDList    int64         `json:"idlists"`
}

type TelcoFremiumProfile struct {
	ProfileId int64  `json:"idProfile"`
	Paid      bool   `json:"paid"`
	Country   string `json:"country"`
	TagName   string `json:"tagName"`
}

type TelcoFremiumUser struct {
	SubscriptionID string              `json:"idSubscription"`
	Profile        TelcoFremiumProfile `json:"profile"`
}
