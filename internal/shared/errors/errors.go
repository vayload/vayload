/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Errors
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package errors

import (
	"errors"
	"fmt"

	httpi "github.com/vayload/vayload/pkg/http"
)

type ErrorKind int

const (
	KindUnknown            ErrorKind = iota // Unknown error
	KindNotFound                            // Resource not found
	KindValidation                          // Validation error
	KindConflict                            // Conflict (duplicate, already exists)
	KindUnauthorized                        // Authentication required
	KindForbidden                           // Permission denied
	KindInternal                            // Internal server error
	KindUnavailable                         // Service unavailable
	KindBadRequest                          // Bad request / invalid input
	KindTimeout                             // Operation timeout
	KindRateLimited                         // Too many requests
	KindPreconditionFailed                  // Precondition failed
)

const (
	// Resource codes
	CodeResourceNotFound       = "RESOURCE_NOT_FOUND"
	CodeResourceAlreadyExists  = "RESOURCE_ALREADY_EXISTS"
	CodeResourceCreationFailed = "RESOURCE_CREATION_FAILED"
	CodeResourceUpdateFailed   = "RESOURCE_UPDATE_FAILED"
	CodeResourceDeletionFailed = "RESOURCE_DELETION_FAILED"

	// Validation codes
	CodeValidationFailed = "VALIDATION_FAILED"
	CodeInvalidInput     = "INVALID_INPUT"
	CodeMissingField     = "MISSING_REQUIRED_FIELD"
	CodeInvalidFormat    = "INVALID_FORMAT"

	// Auth codes
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeTokenExpired       = "TOKEN_EXPIRED"
	CodeTokenInvalid       = "TOKEN_INVALID"
	CodeSessionExpired     = "SESSION_EXPIRED"

	// Operation codes
	CodeOperationFailed    = "OPERATION_FAILED"
	CodeOperationTimeout   = "OPERATION_TIMEOUT"
	CodeOperationCancelled = "OPERATION_CANCELLED"

	// State codes
	CodeConflict           = "CONFLICT"
	CodeAlreadyExists      = "ALREADY_EXISTS"
	CodeAlreadyProcessed   = "ALREADY_PROCESSED"
	CodeExpired            = "EXPIRED"
	CodePreconditionFailed = "PRECONDITION_FAILED"

	// External service codes
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	CodeExternalError      = "EXTERNAL_SERVICE_ERROR"

	// Rate limiting
	CodeRateLimited     = "RATE_LIMITED"
	CodeTooManyAttempts = "TOO_MANY_ATTEMPTS"

	// Database codes (internal use)
	CodeNoRowsAffected = "NO_ROWS_AFFECTED"
	CodeDuplicateEntry = "DUPLICATE_ENTRY"
)

// Err represents a domain error with kind, code, context and optional details
type Err struct {
	Kind       ErrorKind      `json:"kind"`
	StatusCode int            `json:"status_code"` // Kept for backward compatibility
	Code       string         `json:"code"`        // Centralized code (e.g., RESOURCE_NOT_FOUND)
	Reason     string         `json:"reason"`      // Reason for the error (e.g., "not_found", "invalid_input")
	Context    string         `json:"context"`     // Module context (e.g., "users", "payments", "medications")
	Message    string         `json:"message"`     // Human-readable message (informative, in English)
	Details    map[string]any `json:"details,omitempty"`
	Cause      error          `json:"cause,omitempty"`
}

type Error interface {
	Status() int
	Error() string
	Cause() error
	Details() map[string]any
	Message() string
	Code() string
}

type SimpleError struct {
	message string
}

func New(message string) error {
	return &SimpleError{
		message: message,
	}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func (e *SimpleError) Error() string {
	return e.message
}

// NewErr creates a new error (deprecated: use NewDomainErr instead)
func NewErr(status int, code, message string, details map[string]any, cause error) *Err {
	return &Err{
		Kind:       statusToKind(status),
		StatusCode: status,
		Code:       code,
		Message:    message,
		Details:    details,
		Cause:      cause,
	}
}

// NewDomainErr creates a new domain error with ErrorKind and context
func NewDomainErr(kind ErrorKind, code, context, message string, details map[string]any, cause error) *Err {
	return &Err{
		Kind:       kind,
		StatusCode: kindToStatus(kind),
		Code:       code,
		Context:    context,
		Message:    message,
		Details:    details,
		Cause:      cause,
	}
}

func (e *Err) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Context, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Context, e.Code, e.Message)
}

func (e *Err) Status() int {
	return e.StatusCode
}

func (e *Err) GetKind() ErrorKind {
	return e.Kind
}

func (e *Err) GetContext() string {
	return e.Context
}

func (e *Err) GetCode() string {
	return e.Code
}

func statusToKind(status int) ErrorKind {
	switch {
	case status == 400:
		return KindBadRequest
	case status == 401:
		return KindUnauthorized
	case status == 403:
		return KindForbidden
	case status == 404:
		return KindNotFound
	case status == 409:
		return KindConflict
	case status == 412:
		return KindPreconditionFailed
	case status == 422:
		return KindValidation
	case status == 429:
		return KindRateLimited
	case status == 503:
		return KindUnavailable
	case status == 504:
		return KindTimeout
	case status >= 500:
		return KindInternal
	default:
		return KindUnknown
	}
}

func kindToStatus(kind ErrorKind) int {
	switch kind {
	case KindNotFound:
		return 404
	case KindValidation:
		return 422
	case KindBadRequest:
		return 400
	case KindConflict:
		return 409
	case KindUnauthorized:
		return 401
	case KindForbidden:
		return 403
	case KindInternal:
		return 500
	case KindUnavailable:
		return 503
	case KindTimeout:
		return 504
	case KindRateLimited:
		return 429
	case KindPreconditionFailed:
		return 412
	default:
		return 500
	}
}

func MappingErrToHttp(err error) error {
	if err == nil {
		return nil
	}

	// Use type switch for direct type checking (not recursive like errors.As)
	// This prevents matching errors nested in Cause
	switch e := err.(type) {
	case *ErrBuilder:
		return mapErrToHttp(e.Err)
	case *Err:
		return mapErrToHttp(e)
	case *httpi.HttpClientErr:
		return &httpi.Err{
			Status: e.Status,
			Err: httpi.HttpError{
				Code:    e.Code,
				Message: e.Message,
			},
			Cause: e.Cause,
		}
	case *httpi.Err:
		return e
	}

	// Fallback: try errors.As for wrapped errors (but only if direct check failed)
	var builderErr *ErrBuilder
	if errors.As(err, &builderErr) {
		return mapErrToHttp(builderErr.Err)
	}

	var domainErr *Err
	if errors.As(err, &domainErr) {
		return mapErrToHttp(domainErr)
	}

	return httpi.ErrInternal(fmt.Errorf("unexpected error: %w", err))
}

func mapErrToHttp(domainErr *Err) *httpi.Err {
	httpErr := httpi.HttpError{
		Code:    domainErr.Code,
		Message: domainErr.Message,
	}
	if domainErr.Details != nil {
		httpErr.Details = domainErr.Details
	}
	if domainErr.Reason != "" {
		httpErr.Reason = domainErr.Reason
	}

	return &httpi.Err{
		Status: domainErr.StatusCode,
		Err:    httpErr,
		Cause:  domainErr.Cause,
	}
}
