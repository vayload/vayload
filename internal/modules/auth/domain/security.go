/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain


// The UserTokenManager interface defines methods for generating and validating user access tokens.
type UserTokenManager interface {
	GenerateJwtTokenWithRefresh(user *User) (SignedTokenWithRefresh, error)
	GenerateJwtToken(user *User) (SignedToken, error)
	ValidateToken(token string) (*User, error)

	// refresh
	CreateRefreshToken(payload string) string
	ValidateRefreshToken(tokenString string) (*User, error)
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
