/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package auth

import (
	"github.com/vayload/vayload/internal/vayload"
)

type AuthProvider interface {
	// With(userId snowflake.ID, clientId int64) AuthProvider
	// Get() (*domain.UserPolicy, error)
	// GetAuthorized(userId snowflake.ID, clientId int64) (*domain.UserPolicy, error)
	// GetSubscription(clientId int64) (*domain.Subscription, error)
	// SetupWithEmail(email string, profileId int) (*domain.UserPolicy, error)
	// SetupWithPhone(phone string, profileId int) (*domain.UserPolicy, error)
	// Subscribe(source string, subscription *domain.SubscriptionInput) error
	// Unsubscribe(subscriptionId int64, profileId int) error
	// UpdateProfileId(newProfileId int64) error
	// UpdateIdentifier(identifier string) error
	// GetProfiles() map[string]int64
	// GetFreeProfileId() int64
	// GetPremiumProfileId() int64
}

type authProvider struct {
	registry vayload.Container
}

// func NewAuthProvider(registry vayload.Container) *authProvider {
// 	var svc authorization.AuthorizationService
// 	if err := registry.GetInto(authorization.AUTHORIZATION_SERVICE_KEY, &svc); err != nil {
// 		logger.F(err, logger.Fields{"context": "auth_service_resolved"})
// 	}

// 	return &authProvider{
// 		AuthorizationService: svc,
// 		registry:             registry,
// 	}
// }

// func (p *authProvider) With(userId snowflake.ID, clientId int64) AuthProvider {
// 	newService := p.AuthorizationService.With(userId, clientId)

// 	return &authProvider{
// 		AuthorizationService: newService,
// 		registry:             p.registry,
// 	}
// }
