/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

// =================================== password strategy ===================================
type PasswordStrategy interface {
	HashPassword(password string) string
	VerifyPassword(password, hashedPassword string) bool
	VerboseVerifyPassword(password, hashedPassword string) (valid bool, algo string)
}

const PasswordStrategyType = "password"

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
