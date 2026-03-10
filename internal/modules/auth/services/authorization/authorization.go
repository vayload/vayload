/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package authorization

import (
	"errors"
	"strconv"
	"sync"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

const (
	AUTHORIZATION_SERVICE_KEY = "authorization_service"
)

type AuthorizationService interface {
	With(userId snowflake.ID, clientId int64) AuthorizationService
	Get() (*domain.UserPolicy, error)
	GetAuthorized(userId snowflake.ID, clientId int64) (*domain.UserPolicy, error)
	GetSubscription(clientId int64) (*domain.Subscription, error)
	SetupWithEmail(email string, profileId int) (*domain.UserPolicy, error)
	SetupWithPhone(phone string, profileId int) (*domain.UserPolicy, error)
	Subscribe(source string, subscription *domain.SubscriptionInput) error
	Unsubscribe(subscriptionId int64, profileId int) error
	UpdateProfileId(newProfileId int64) error
	UpdateIdentifier(identifier string) error
	GetProfiles() map[string]int64
	GetFreeProfileId() int64
	GetPremiumProfileId() int64
}

type authorizationService struct {
	repository domain.AuthorizationRepository
	productId  int

	profiles map[string]int64
	mu       sync.RWMutex

	// Seeds
	UserId   snowflake.ID
	ClientId int64
}

func NewAuthorizationService(repository domain.AuthorizationRepository) *authorizationService {
	return &authorizationService{
		repository: repository,
		productId:  10052,
		profiles: map[string]int64{
			"free":    546,
			"premium": 547,
		},
	}
}

func (service *authorizationService) With(userId snowflake.ID, clientId int64) AuthorizationService {
	service.UserId = userId
	service.ClientId = clientId

	return service
}

func (service *authorizationService) Get() (*domain.UserPolicy, error) {
	if service.UserId == 0 || service.ClientId == 0 {
		return nil, errors.New("user or client not set")
	}

	return service.GetAuthorized(service.UserId, service.ClientId)
}

func (service *authorizationService) GetAuthorized(userId snowflake.ID, clientId int64) (*domain.UserPolicy, error) {
	// Binding is relationship between user and client
	binding := &domain.AuthorizationBindingOptional{}

	// When client id is zero get authorization binding
	if clientId == 0 {
		var err error
		binding, err = service.repository.GetAuthorizationBinding(strconv.Itoa(int(userId)))
		if err != nil {
			return nil, err
		}
	} else {
		binding.ClientId = &clientId
	}

	if binding.ClientId == nil || *binding.ClientId <= 0 {
		return nil, domain.ErrUserNotFound(errors.New("client id not found"))
	}

	rawUserPolicy, err := service.repository.GetPermissions(int(*binding.ClientId), service.productId)
	if err != nil {
		return nil, err
	}

	return domain.NewUserPolicy(rawUserPolicy), nil
}

func (service *authorizationService) GetSubscription(clientId int64) (*domain.Subscription, error) {
	return service.repository.GetSubscription(clientId, service.productId)
}

func (service *authorizationService) SetupWithEmail(email string, profileId int) (*domain.UserPolicy, error) {
	rawUserPolicy, err := service.repository.Setup(int(domain.PolicyEmailChannel), email, service.productId, profileId)
	if err != nil {
		return nil, err
	}

	return domain.NewUserPolicy(rawUserPolicy), nil
}

func (service *authorizationService) SetupWithPhone(phone string, profileId int) (*domain.UserPolicy, error) {
	rawUserPolicy, err := service.repository.Setup(int(domain.PolicyPhoneChannel), phone, service.productId, profileId)
	if err != nil {
		return nil, err
	}

	return domain.NewUserPolicy(rawUserPolicy), nil
}

func (service *authorizationService) Subscribe(source string, subscription *domain.SubscriptionInput) error {
	ok, err := service.repository.Subscribe(service.ClientId, service.productId, source, subscription)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrUserNotFound(errors.New("failed to subscribe user"))
	}

	return nil
}

func (service *authorizationService) Unsubscribe(subscriptionId int64, profileId int) error {
	ok, err := service.repository.Unsubscribe(service.ClientId, service.productId, subscriptionId, profileId)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrUserNotFound(errors.New("failed to unsubscribe user"))
	}

	return nil
}

func (service *authorizationService) UpdateProfileId(newProfileId int64) error {
	if err := service.repository.UpdateProfileId(service.UserId, service.ClientId, int(newProfileId), service.productId); err != nil {
		return err
	}

	return nil
}

func (service *authorizationService) UpdateIdentifier(identifier string) error {
	if _, err := service.repository.UpdateIdentifier(int(service.ClientId), service.productId, identifier); err != nil {
		return err
	}

	return nil
}

func (service *authorizationService) GetProfiles() map[string]int64 {
	service.mu.RLock()
	defer service.mu.RUnlock()

	return service.profiles
}

func (service *authorizationService) GetFreeProfileId() int64 {
	service.mu.RLock()
	defer service.mu.RUnlock()

	return service.profiles["free"]
}

func (service *authorizationService) GetPremiumProfileId() int64 {
	service.mu.RLock()
	defer service.mu.RUnlock()

	return service.profiles["premium"]
}

func (service *authorizationService) SetProfiles(profiles map[string]int64) {
	service.mu.Lock()
	defer service.mu.Unlock()
	service.profiles = profiles
}
