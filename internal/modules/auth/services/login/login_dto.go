/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package login

import (
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type IdentifyInput struct {
	UserIdentifier string `json:"user_identifier" validate:"required,email-or-phone"`           // Email or phone
	IdentifierType string `json:"identifier_type" validate:"required,oneof=email phone"`        // e.g. "email", "phone"
	ForFactor      string `json:"for_factor" validate:"required,oneof=password otp magic_link"` // e.g. "password", "otp", "magic_link"
	ClientType     string `json:"client_type" validate:"required,oneof=web mobile"`             // e.g. "web", "mobile"
}

type LoginInput struct {
	Identifier string `json:"identifier" validate:"required,email-or-phone"`
	Password   string `json:"password" validate:"required,min=8"`
}

type TokenIdInput struct {
	TokenID string `json:"token_id" validate:"required"`
}

type OtpCodeInput struct {
	Identifier string `json:"identifier" validate:"required,email-or-phone"`
	Code       string `json:"code" validate:"required"`
}

type MagicLinkInput struct {
	Token string `json:"token"`
}

type MagicLink struct {
	Code string `json:"code"`
}

type OtpCodeGenInput struct {
	Identifier string `json:"identifier"`
	Channel    string `json:"channel"`
}

type MagicLinkGenInput struct {
	Redirect       string `json:"redirect"`        // URL to redirect after login
	UserIdentifier string `json:"user_identifier"` // Email or phone
	IdentifierType string `json:"identifier_type"` // e.g. "email", "phone"
	ChannelToSend  string `json:"channel_to_send"` // e.g. "email", "phone", "whatsapp"
	Meta           any    `json:"meta,omitempty"`  // Additional metadata
}

type SetupUserInput struct {
	Username       string `json:"username" validate:"required,min=2,max=200"`
	Identifier     string `json:"identifier" validate:"required,email-or-phone"`
	IdentifierType string `json:"identifier_type" validate:"required,oneof=email phone"`

	CountryId snowflake.ID `json:"country_id" validate:"required"`
	ProfileId int          `json:"profile_id"` // The profile id as linking this current user
	Method    string       `json:"method"`     // e.g. "email", "sms", "oauth", "sso"
}

type TelcoInput struct {
	Code     string `json:"code"`
	ClientId uint32 `json:"client_id"`
	Origin   string `json:"origin"`
	Profile  string `json:"profile"`
	Product  string `json:"product"`
	Idlist   string `json:"idlist"`
}

type TelcoUserSync struct {
	Phone     string       `json:"phone"`
	UserId    snowflake.ID `json:"user_id"`
	ClientId  int64        `json:"client_id"`
	CoreId    string       `json:"core_id"`
	IdList    string       `json:"idlist"`
	ProductId int64        `json:"product_id"`
	CountryID snowflake.ID `json:"country_id"`
}

type AuthStepCredentials struct {
	UserIdentifier string              `json:"user_identifier"`         // Email or phone
	Factor         string              `json:"factor"`                  // e.g. "password", "otp", "magic_link"
	ClientType     string              `json:"client_type"`             // e.g. "web", "mobile"
	Destinations   map[string]any      `json:"destinations,omitempty"`  // e.g. {"email": "<email>", "phone": "<phone>"}
	OtpProviders   *domain.OtpProvider `json:"otp_providers,omitempty"` // Optional, for OTP factors
}
