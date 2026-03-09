/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import "time"

// The UserTokenManager interface defines methods for generating and validating user access tokens.
type UserTokenManager interface {
	GenerateJwtTokenWithRefresh(user *AuthUser) (SignedTokenWithRefresh, error)
	GenerateJwtToken(user *AuthUser) (SignedToken, error)
	ValidateToken(token string) (*AuthUser, error)

	// refresh
	CreateRefreshToken(payload string) string
	ValidateRefreshToken(tokenString string) (*AuthUser, error)
}

// The SignedToken interface defines methods for accessing token properties.
type SignedToken interface {
	GetPayload() any
	GetMeta() map[string]any
	GetAccessToken() string
	GetExpiresAt() time.Time
	GetExpiresIn() int64
}

// SignedTokenWithRefresh extends SignedToken with refresh token capabilities.
type SignedTokenWithRefresh interface {
	SignedToken
	GetRefreshToken() string
	GetRefreshTokenExpiresAt() time.Time
	GetRefreshTokenExpiresIn() int64
}

// The SecureRandomizer interface defines methods for generating random strings and codes.
type SecureRandomizer interface {
	// GenerateRandomString returns a secure random alphanumeric string of specified length in base64 encoding.
	GenerateRandomString(length int, urlSafe bool) string

	// GenerateRandomNumericCode returns a secure random numeric code (as string) of specified length (with leading zeros).
	GenerateRandomNumericCode(length int) string

	// GenerateRandomBytes returns secure random bytes.
	GenerateRandomBytes(length int) ([]byte, error)

	// GenerateUUID returns a secure random UUIDv4 string.
	GenerateUUID() string

	// SecureCompare compares two strings in constant time to avoid timing attacks.
	SecureCompare(a, b string) bool
}
