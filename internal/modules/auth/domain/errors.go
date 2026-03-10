/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import (
	"github.com/vayload/vayload/internal/shared/errors"
)

const ErrContext = "auth"

func ErrInvalidCredentials(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeInvalidCredentials, ErrContext, "Invalid email or password", nil, cause)
}

func ErrInvalidOtpCode(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeInvalidCredentials, ErrContext, "Invalid OTP code", nil, cause)
}

func ErrUserNotFound(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindNotFound, errors.CodeResourceNotFound, ErrContext, "User not found", map[string]any{"resource": "user"}, cause)
}

func ErrEmailAlreadyExists(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindConflict, errors.CodeAlreadyExists, ErrContext, "Email already exists", map[string]any{"field": "email"}, cause)
}

func ErrInvalidToken(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeTokenInvalid, ErrContext, "Token is invalid or expired", nil, cause)
}

func ErrTokenExpired(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeTokenExpired, ErrContext, "Token has expired", nil, cause)
}

func ErrRefreshTokenInvalid(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeTokenInvalid, ErrContext, "Refresh token is invalid", map[string]any{"token_type": "refresh"}, cause)
}

func ErrUnauthorized(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeUnauthorized, ErrContext, "Unauthorized access", nil, cause)
}

func ErrForbidden(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindForbidden, errors.CodeForbidden, ErrContext, "Access forbidden", nil, cause)
}

func ErrAccountLocked(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindForbidden, errors.CodeForbidden, ErrContext, "Account is locked", map[string]any{"reason": "locked"}, cause)
}

func ErrAccountDisabled(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindForbidden, errors.CodeForbidden, ErrContext, "Account is disabled", map[string]any{"reason": "disabled"}, cause)
}

func ErrPasswordResetRequired(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindPreconditionFailed, errors.CodePreconditionFailed, ErrContext, "Password reset is required", map[string]any{"action_required": "password_reset"}, cause)
}

func ErrSessionExpired(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindUnauthorized, errors.CodeSessionExpired, ErrContext, "Session has expired", nil, cause)
}

func ErrInvalidPasswordFormat(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindValidation, errors.CodeValidationFailed, ErrContext, "Password does not meet requirements", map[string]any{"field": "password"}, cause)
}

func ErrTooManyAttempts(cause error) *errors.Err {
	return errors.NewDomainErr(errors.KindRateLimited, errors.CodeTooManyAttempts, ErrContext, "Too many login attempts, please try again later", nil, cause)
}

func ErrJwtTokenGenerationFailed(cause error) *errors.ErrBuilder {
	return errors.Internal("Failed to generate JWT token").Context("auth.jwt").Cause(cause)
}

func ErrMagicLinkValidation(cause error) *errors.ErrBuilder {
	return errors.Internal("Magic link validation failed").Context("auth.magic-link").Cause(cause)
}

// ============================================================================
// Telco-specific Errors
// ============================================================================

// ErrTelcoCustodyNotFound - invalid or expired code (404)
func ErrTelcoCustodyNotFound(reason string, cause error) *errors.ErrBuilder {
	return errors.NotFound("Could not verify your access. Please try again.").Context("auth.telco-login").Reason(reason).Cause(cause)
}

// ErrTelcoClientNotFound - client does not exist in Telco (404)
func ErrTelcoClientNotFound(reason string) *errors.ErrBuilder {
	return errors.NotFound("Account not found. Please verify your information.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoProductNotFound - internal configuration missing (500)
func ErrTelcoProductNotFound(reason string) *errors.ErrBuilder {
	return errors.Internal("Service unavailable. Please try again later.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoServicesFailed - external service down (503)
func ErrTelcoServicesFailed(reason string) *errors.ErrBuilder {
	return errors.ExternalError("Service temporarily unavailable.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoServiceTimeout - service timeout (504)
func ErrTelcoServiceTimeout(reason string) *errors.ErrBuilder {
	return errors.Timeout("Request took too long. Please try again.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoParsingFailed - parsing error, possible contract change (500)
func ErrTelcoParsingFailed(reason string) *errors.ErrBuilder {
	return errors.Internal("Error processing request.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoValidationFailed - invalid service response (500)
func ErrTelcoValidationFailed(reason string) *errors.ErrBuilder {
	return errors.Internal("Error processing request.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoNetworkError - network issue (503)
func ErrTelcoNetworkError(reason string) *errors.ErrBuilder {
	return errors.Unavailable("Could not connect to service. Please try again later.").Context("auth.telco-login").Reason(reason)
}

// ErrTelcoRateLimited - rate limited (429)
func ErrTelcoRateLimited(reason string) *errors.ErrBuilder {
	return errors.RateLimited("Too many attempts. Please wait a moment.").Context("auth.telco-login").Reason(reason)
}

var (
	// General errors
	ErrNotResults = errors.New("no results found")

	// Database/Repository errors
	ErrEmptyResultSet   = errors.New("repository: empty result set")
	ErrNoRowsAffected   = errors.New("repository: no rows affected")
	ErrDuplicateEntry   = errors.New("repository: duplicate entry")
	ErrForeignKeyFailed = errors.New("repository: foreign key constraint failed")
	ErrRecordNotFound   = errors.New("repository: record not found")

	// Delete operation errors
	ErrDeleteFailed    = errors.New("repository: delete operation failed")
	ErrNothingToDelete = errors.New("repository: nothing to delete")

	// Update operation errors
	ErrUpdateFailed    = errors.New("repository: update operation failed")
	ErrNothingToUpdate = errors.New("repository: nothing to update")

	// Insert operation errors
	ErrInsertFailed = errors.New("repository: insert operation failed")

	// Authentication specific errors
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrUserInactive    = errors.New("user is inactive")
)
