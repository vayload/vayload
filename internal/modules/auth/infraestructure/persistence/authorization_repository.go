/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
	httpi "github.com/vayload/vayload/pkg/http"
)

type authorizationRepository struct {
	database   connection.DatabaseConnection
	httpClient *httpi.HttpClient
}

type AuthorizationRepositoryConfig struct {
	ProductURL   string
	ProductToken string
}

func NewAuthorizationRepository(database connection.DatabaseConnection, config AuthorizationRepositoryConfig) *authorizationRepository {
	httpClient := httpi.NewHttpClient(httpi.HttpClientConfig{
		BaseURL: config.ProductURL,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", config.ProductToken),
		},
	})

	return &authorizationRepository{
		database:   database,
		httpClient: httpClient,
	}
}

func (repository *authorizationRepository) Create(channelId int, identifier string, productId int) (any, error) {
	ctx := context.Background()
	body := map[string]any{
		"product_id": productId,
		"channel_id": channelId,
		"identifier": identifier,
	}
	type data struct {
		ID    string `json:"client_id"`
		IsNew bool   `json:"is_new_client"`
	}

	response, err := httpi.UnwrapBody[data](repository.httpClient.Post(ctx, "/clients", body))
	if err != nil {
		return nil, err
	}

	return response.Data, nil

}

func (repository *authorizationRepository) GetPermissions(clientId int, productId int) (*domain.RawUserPolicy, error) {
	ctx := context.Background()

	path := fmt.Sprintf("/clients/%d?product_id=%d&return=%s", clientId, productId, "minimal")
	response, err := httpi.UnwrapBody[domain.RawUserPolicy](repository.httpClient.Get(ctx, path))
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (repository *authorizationRepository) Setup(channelId int, identifier string, productId int, profileId int) (*domain.RawUserPolicy, error) {
	ctx := context.Background()
	body := map[string]any{
		"product_id": productId,
		"channel_id": channelId,
		"identifier": identifier,
		"profile_id": profileId,
	}

	response, err := httpi.UnwrapBody[domain.RawUserPolicy](repository.httpClient.Post(ctx, "/clients/setup?return=minimal", body))
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (repository *authorizationRepository) UpdateIdentifier(clientId int, productId int, newIdentifier string) (bool, error) {
	ctx := context.Background()
	body := map[string]any{
		"product_id": productId,
		"channel_id": clientId,
		"identifier": newIdentifier,
	}

	response, err := repository.httpClient.Patch(ctx, "/clients/"+strconv.Itoa(clientId), body)
	if err != nil || response.StatusCode != 200 {
		return false, fmt.Errorf("failed to update identifier: %w", err)
	}

	return true, nil
}

func (repository *authorizationRepository) Subscribe(clientId int64, productId int, source string, subscription *domain.SubscriptionInput) (bool, error) {
	ctx := context.Background()
	data := map[string]any{
		"client_id":    clientId,
		"product_id":   productId,
		"source":       source,
		"subscription": subscription,
	}

	response, err := repository.httpClient.Post(ctx, "/subscriptions", data)
	if err != nil {
		return false, err
	}

	return (response.StatusCode == 200 || response.StatusCode == 201), nil
}

func (repository *authorizationRepository) Unsubscribe(clientId int64, productId int, subscriptionId int64, profileId int) (bool, error) {
	ctx := context.Background()
	data := map[string]any{
		"id":         subscriptionId,
		"client_id":  clientId,
		"product_id": productId,
		"profile_id": profileId, // The profileId to change when subscription is cancelled
	}

	_, err := repository.httpClient.Post(ctx, "/subscriptions/cancel", data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repository *authorizationRepository) GetSubscription(clientId int64, productId int) (*domain.Subscription, error) {
	ctx := context.Background()
	response, err := repository.httpClient.Get(ctx, fmt.Sprintf("/subscriptions/query?product_id=%d&client_id=%d", productId, clientId))
	if err != nil {
		return nil, err
	}

	var subscriptionData struct {
		Data domain.Subscription `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&subscriptionData); err != nil {
		return nil, err
	}

	return &subscriptionData.Data, nil
}

func (repository *authorizationRepository) GetAuthorizationBinding(userId string) (*domain.AuthorizationBindingOptional, error) {
	ctx := context.Background()
	query := `SELECT id, client_id, profile_id, profile_signature FROM users WHERE id = ?`

	rawResult := AuthorizationBindingModel{}
	filterId, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}

	err = repository.database.SelectOne(ctx, &rawResult, query, int64(filterId))
	if err == nil {
		return nil, err
	}

	return rawResult.MapToDomain(), nil
}

func (repository *authorizationRepository) UpdateProfileId(userId snowflake.ID, clientId int64, profileId int, productId int) error {
	ctx := context.Background()
	body := map[string]any{
		"product_id": productId,
		"profile_id": profileId, // profile id to change
	}

	response, err := repository.httpClient.Put(ctx, fmt.Sprintf("/clients/%d", clientId), body)
	if err != nil {
		return fmt.Errorf("failed to update profile id: %w", err)
	}

	if response.StatusCode != 200 {
		var errorResponse struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if errDec := json.NewDecoder(response.Body).Decode(&errorResponse); errDec != nil {
			return fmt.Errorf("failed to update profile id: %w", errDec)
		}

		return fmt.Errorf("%s", errorResponse.Message)
	}

	query := `UPDATE users SET client_id = ?, profile_id = ? WHERE id = ?`
	err = repository.database.Prepared(ctx, query, clientId, profileId, userId)
	if err != nil {
		return fmt.Errorf("failed to update profile id: %w", err)
	}

	return nil
}
