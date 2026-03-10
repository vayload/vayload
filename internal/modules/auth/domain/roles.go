/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Domain
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package domain

import "github.com/vayload/vayload/internal/shared/snowflake"

// The binding between a user and a client application
// with optional profile association.
// Values can be nil to indicate absence.
type AuthorizationBindingOptional struct {
	UserId    snowflake.ID `json:"user_id"`
	ClientId  *int64       `json:"client_id"`  // Optional client ID (because a user may not be associated with a specific client)
	ProfileId *int64       `json:"profile_id"` // Optional profile ID (because a user may not be associated with a specific profile)
	Signature string       `json:"signature"`
}

type AuthorizationBinding struct {
	UserId    snowflake.ID `json:"user_id"`
	ClientId  int64        `json:"client_id"`
	ProfileId int64        `json:"profile_id"`
	Signature string       `json:"signature"`
}

type AuthorizationRepository interface {
	Create(channelID int, identifier string, productId int) (any, error)
	// Get user permissions and quotas
	GetPermissions(clientId int, productId int) (*RawUserPolicy, error)
	// ProvisionUser creates a new user with the given parameters if not exists or updates existing user
	Setup(channelID int, identifier string, productId int, profileId int) (*RawUserPolicy, error)

	UpdateIdentifier(clientId int, productId int, newIdentifier string) (bool, error)

	// Subscribe links subscription to user
	Subscribe(clientId int64, productId int, source string, subscription *SubscriptionInput) (bool, error)
	// Unsubscribe removes subscription from user (use profileId to change in on unsubscribe)
	Unsubscribe(clientId int64, productId int, subscriptionId int64, profileId int) (bool, error)

	// GetSubscription retrieves user subscription
	GetSubscription(clientId int64, productId int) (*Subscription, error)

	GetAuthorizationBinding(userId string) (*AuthorizationBindingOptional, error)

	UpdateProfileId(userId snowflake.ID, clientId int64, profileId int, productId int) error
}

type PayTransaction struct {
	ID             string  `json:"id"`
	SubscriptionID string  `json:"subscription_id"`
	Status         string  `json:"status"`
	Method         string  `json:"method"`
	Amount         float64 `json:"amount"`
	Paid           int     `json:"paid"`
}

type SubscriptionInput struct {
	PlanID          string  `json:"plan_id"`
	ProviderPlanID  string  `json:"provider_plan_id"`
	RecurrenceID    string  `json:"recurrence_id"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
	OriginID        string  `json:"origin_id"`
	TrialPeriod     bool    `json:"trial_period"`
	StartDate       string  `json:"start_date"`
	TrialEndDate    string  `json:"trial_end_date"`
	EndDate         string  `json:"end_date"`
	NextBillingDate *string `json:"next_billing_date"`
}
