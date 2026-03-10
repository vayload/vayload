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

	"github.com/vayload/vayload/internal/shared/snowflake"
)

type AuthRepository interface {
	FindByIdentifier(ctx context.Context, identifier string, identifierType string) (*User, error)
	// find user by identifier, gets user identity and metadata like phone, email, country_id, passport, etc.
	FindUserIdentity(ctx context.Context, identifier string, identifierType string) (*UserIdentity, error)
	// FindMetadataByEmail(ctx context.Context, email string) (*UserMetadata, error)
	// FindMetadataByIdentifier(ctx context.Context, identifier string, identifierType string) (*UserMetadata, error)
	FindUserAuthByIdentifier(ctx context.Context, identifier string, identifierType string) (*UserAuth, error)

	UpdateMetadata(ctx context.Context, metadata *UserMetadata) error
	UpdatePassword(ctx context.Context, userId snowflake.ID, newPassword string) error

	CreateUserWithSettings(ctx context.Context, user *User, settings *UserSettings) error
	CreateUserWithCode(ctx context.Context, user *User, token string, code string) (*User, error)
	SaveOtpCode(ctx context.Context, userId snowflake.ID, otpCode string) error
	FindCodesByIdentifier(ctx context.Context, identifier string, identifierType string, typeCode string) (*UserVerification, error)
	BindAuthorization(ctx context.Context, userId snowflake.ID, binding *AuthorizationBinding, meta *UserMeta) error
	UpdateVerificationCode(ctx context.Context, userId snowflake.ID, code string, typeCode string) error

	// Password recovery
	SaveRecoveryToken(ctx context.Context, userId snowflake.ID, token string) error
	FindUserByRecoveryToken(ctx context.Context, token string) (*User, error)
	ResetPasswordWithRecoveryToken(ctx context.Context, token string, newPassword string) error

	// Email change during registration
	SaveEmailChangeRequest(ctx context.Context, userId snowflake.ID, newEmail string, currentToken string, newToken string) error
	FindUserByEmailChangeTokens(ctx context.Context, currentToken string, newToken string) (*User, error)
	ApplyEmailChange(ctx context.Context, userId snowflake.ID) error

	FindCountryOtpProviders(ctx context.Context, countryId snowflake.ID) (*OtpProvider, error)

	// Telco
	GetTelcoProfiles(ctx context.Context, telcoId string, countryId snowflake.ID) ([]int64, error)
	GetTelcoLocation(ctx context.Context, productName string, countryCode string) (*TelcoLocation, error)
}

type AuthLogRepository interface {
	SaveLogin(ctx context.Context, userId snowflake.ID, ipAddress string, userAgent string, method string, phone *string, email string) error
	SaveExternalRequestLog(ctx context.Context, log *ExternalRequestLogging) error
}
