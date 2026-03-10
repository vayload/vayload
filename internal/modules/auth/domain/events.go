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

	"github.com/vayload/vayload/internal/vayload"
)

// This event is triggered when a user successfully logs in
type UserLoggedInEvent struct {
	User    *User
	Context *AuthContext
}

func (event *UserLoggedInEvent) Name() string {
	return "auth.UserLoggedIn"
}

// This event is triggered when a new user is created
type UserCreatedEvent struct {
	User *User
	Code string
}

func (event *UserCreatedEvent) Name() string {
	return "auth.UserCreated"
}

// This event is triggered when a user is updated
type UserUpdatedEvent struct {
	User *User
}

func (event *UserUpdatedEvent) Name() string {
	return "auth.UserUpdated"
}

// This event is triggered when a user update code is generated
type UserUpdateCodeEvent struct {
	User *User
	Code string
}

func (event *UserUpdateCodeEvent) Name() string {
	return "auth.UserUpdateCode"
}

// This event is triggered when an OTP code is generated
type OtpCodeGeneratedEvent struct {
	User    *User
	Code    string
	Channel string
}

func (event *OtpCodeGeneratedEvent) Name() string {
	return "auth.OtpCodeGenerated"
}

// This event is triggered when a magic link is generated
type UserMagicLinkGeneratedEvent struct {
	User      *User
	Code      string
	Channel   string
	ExpiresIn time.Duration
}

func (event *UserMagicLinkGeneratedEvent) Name() string {
	return "auth.UserMagicLinkGenerated"
}

// This event is triggered when a user's email is verified
type UserEmailVerifiedEvent struct {
	User *User
}

func (event *UserPasswordRecoveryRequestedEvent) Name() string {
	return "auth.UserPasswordRecoveryRequested"
}

// This event is triggered when a user requests a password recovery
type UserPasswordRecoveryRequestedEvent struct {
	User  *User
	Token string
}

func (event *UserPasswordResetCompletedEvent) Name() string {
	return "auth.UserPasswordResetCompleted"
}

// This event is triggered when a user completes a password reset
type UserPasswordResetCompletedEvent struct {
	User *User
}

func (event *UserEmailChangeRequestedEvent) Name() string {
	return "auth.UserEmailChangeRequested"
}

// This event is triggered when a user requests an email change
type UserEmailChangeRequestedEvent struct {
	User         *User
	CurrentToken string
	NewToken     string
	NewEmail     string
}

func (event *UserEmailChangeConfirmedEvent) Name() string {
	return "auth.UserEmailChangeConfirmed"
}

// This event is triggered when a user confirms an email change
type UserEmailChangeConfirmedEvent struct {
	User     *User
	NewEmail string
}

// Domain expose only publish method
// to avoid direct access to the event bus
type EventBus interface {
	Publish(ctx context.Context, event vayload.Event)
}
