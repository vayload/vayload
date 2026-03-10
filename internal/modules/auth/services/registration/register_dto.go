/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package registration

type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

// Validate registration with token and code
type RegisterValidationInput struct {
	Identifier string `json:"identifier"` // Email or phone
	Type       string `json:"type"`       // e.g. "email", "phone"
	Token      string `json:"token"`
	Code       string `json:"code"`
}

type VerifyEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}

type TokenIdInput struct {
	TokenID string `json:"token_id" validate:"required"`
}

type RegisterEmailChangeInput struct {
	Identifier string
	Type       string
	NewEmail   string
}

type RegisterEmailChangeConfirmInput struct {
	CurrentToken string
	NewToken     string
}
