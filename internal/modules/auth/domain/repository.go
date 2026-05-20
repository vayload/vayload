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

type AuthRepository interface {
	FindByID(ctx context.Context, id snowflake.ID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByIdentifier(ctx context.Context, identifier string, identifierType IdentifierType) (*User, error)

	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id snowflake.ID) error

	UpdatePassword(ctx context.Context, userID snowflake.ID, passwordHash string) error
	UpdateLastSignIn(ctx context.Context, userID snowflake.ID, ip, userAgent string) error
	SaveOtpCode(ctx context.Context, userID snowflake.ID, code string) error

	// Verification tokens
	SaveConfirmationToken(ctx context.Context, userID snowflake.ID, token string) error
	SaveRecoveryToken(ctx context.Context, userID snowflake.ID, token string) error
	
	FindUserByConfirmationToken(ctx context.Context, token string) (*User, error)
	FindUserByRecoveryToken(ctx context.Context, token string) (*User, error)
	
	ConfirmEmail(ctx context.Context, userID snowflake.ID) error
	ResetPasswordWithRecoveryToken(ctx context.Context, token string, hashedPassword string) error

	// Email change
	SaveEmailChangeRequest(ctx context.Context, userID snowflake.ID, newEmail string, currentToken string, newToken string) error
	FindUserByEmailChangeTokens(ctx context.Context, currentToken string, newToken string) (*User, error)
	ApplyEmailChange(ctx context.Context, userID snowflake.ID) error
}

type AuthLogRepository interface {
	SaveLogin(ctx context.Context, userID snowflake.ID, ipAddress string, userAgent string, method string, email string) error
}

type ExternalRequestLogging struct {
	RequestID      string
	URL            string
	Method         string
	RequestBody    string
	RequestAt      time.Time
	ResponseBody   string
	ResponseStatus int
	RequestElapsed int64
}
