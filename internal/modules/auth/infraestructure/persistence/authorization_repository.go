/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Infraestructure/Persistence/Authorization
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package persistence

import (
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type AuthorizationRepositoryConfig struct {
}

type AuthorizationRepository struct {
	db     connection.DatabaseConnection
	config AuthorizationRepositoryConfig
}

func NewAuthorizationRepository(db connection.DatabaseConnection, config AuthorizationRepositoryConfig) *AuthorizationRepository {
	return &AuthorizationRepository{
		db:     db,
		config: config,
	}
}

func (r *AuthorizationRepository) Create(channelID int, identifier string, productId int) (any, error) {
	return nil, nil
}

func (r *AuthorizationRepository) GetPermissions(clientId int, productId int) (*domain.RawUserPolicy, error) {
	return nil, nil
}

func (r *AuthorizationRepository) Setup(channelID int, identifier string, productId int, profileId int) (*domain.RawUserPolicy, error) {
	return nil, nil
}

func (r *AuthorizationRepository) UpdateIdentifier(clientId int, productId int, newIdentifier string) (bool, error) {
	return false, nil
}

func (r *AuthorizationRepository) Subscribe(clientId int64, productId int, source string, subscription *domain.SubscriptionInput) (bool, error) {
	return false, nil
}

func (r *AuthorizationRepository) Unsubscribe(clientId int64, productId int, subscriptionId int64, profileId int) (bool, error) {
	return false, nil
}

func (r *AuthorizationRepository) GetSubscription(clientId int64, productId int) (*domain.Subscription, error) {
	return nil, nil
}

func (r *AuthorizationRepository) GetAuthorizationBinding(userId string) (*domain.AuthorizationBindingOptional, error) {
	return nil, nil
}

func (r *AuthorizationRepository) UpdateProfileId(userId snowflake.ID, clientId int64, profileId int, productId int) error {
	return nil
}

var _ domain.AuthorizationRepository = (*AuthorizationRepository)(nil)
