/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Services/Registration
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package registration

import (
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/auth/services/authorization"
)

type RegistrationStrategies struct {
	Password domain.PasswordStrategy
}

type RegisterService struct {
	repository    domain.AuthRepository
	tokenManager  domain.UserTokenManager
	strategies    *RegistrationStrategies
	randomizer    domain.SecureRandomizer
	EventBus      domain.EventBus
	authorization authorization.AuthorizationService
}

func NewRegisterService(
	repository domain.AuthRepository,
	tokenManager domain.UserTokenManager,
	strategies *RegistrationStrategies,
	randomizer domain.SecureRandomizer,
	EventBus domain.EventBus,
	authorization authorization.AuthorizationService,
) *RegisterService {
	return &RegisterService{
		repository:    repository,
		tokenManager:  tokenManager,
		strategies:    strategies,
		randomizer:    randomizer,
		EventBus:      EventBus,
		authorization: authorization,
	}
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

func (service *RegisterService) RegisterUser(ctx interface{}, input RegisterInput) (*domain.User, error) {
	return nil, nil
}

type RegisterValidationInput struct {
	Identifier string
	Type       string
	Code       string
}

func (service *RegisterService) ValidateRegister(ctx interface{}, input RegisterValidationInput) (*domain.OAuthSession, error) {
	return nil, nil
}

func (service *RegisterService) SendEmailVerificationCode(ctx interface{}, dto RegisterValidationInput) error {
	return nil
}

type RegisterEmailChangeInput struct {
	Identifier string
	Type       string
	NewEmail   string
}

func (service *RegisterService) RequestEmailVerificationChange(ctx interface{}, input RegisterEmailChangeInput) error {
	return nil
}

type RegisterEmailChangeConfirmInput struct {
	CurrentToken string
	NewToken     string
}

func (service *RegisterService) ConfirmEmailVerificationChange(ctx interface{}, input RegisterEmailChangeConfirmInput) error {
	return nil
}

var _ interface{} = (*RegisterService)(nil)
