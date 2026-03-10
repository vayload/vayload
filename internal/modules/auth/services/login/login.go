/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Services/Login
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package login

import (
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/auth/services/authorization"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type AuthStrategies struct {
	OAuth2   domain.OAuth2StrategyFacade
	Password domain.PasswordStrategy
	OtpCode  domain.OtpStrategy
}

type LoginService struct {
	repository    domain.AuthRepository
	tokenManager  domain.UserTokenManager
	randomizer    domain.SecureRandomizer
	EventBus      domain.EventBus
	strategies    *AuthStrategies
	authorization authorization.AuthorizationService
}

func NewLoginService(
	repository domain.AuthRepository,
	tokenManager domain.UserTokenManager,
	randomizer domain.SecureRandomizer,
	EventBus domain.EventBus,
	authStrategies *AuthStrategies,
	authorization authorization.AuthorizationService,
) *LoginService {
	return &LoginService{
		repository:    repository,
		tokenManager:  tokenManager,
		randomizer:    randomizer,
		EventBus:      EventBus,
		strategies:    authStrategies,
		authorization: authorization,
	}
}

type IdentifyInput struct {
	UserIdentifier string
	IdentifierType string
	ForFactor      string
	ClientType     string
}

type AuthStepCredentials struct {
	UserIdentifier string
	Factor         string
	ClientType     string
	OtpProviders   *domain.OtpProvider
	Destinations   map[string]any
}

func (service *LoginService) IdentifyUserForFactor(ctx interface{}, input IdentifyInput) (*AuthStepCredentials, error) {
	return nil, nil
}

type OtpCodeGenInput struct {
	Identifier string
	Channel    string
}

func (service *LoginService) GenerateOtpCode(ctx interface{}, input OtpCodeGenInput) error {
	return nil
}

type LoginInput struct {
	Identifier string
	Password   string
}

func (service *LoginService) LoginWithPassword(ctx interface{}, input LoginInput, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	return nil, nil
}

func (service *LoginService) GetOAuth2URL(ctx interface{}, provider domain.OAuth2Provider, payload domain.OAuth2State) (string, error) {
	return "", nil
}

func (service *LoginService) LoginWithOAuth2(ctx interface{}, provider domain.OAuth2Provider, code string, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	return nil, nil
}

type OtpCodeInput struct {
	Identifier string
	Code       string
}

func (service *LoginService) LoginWithOtpCode(ctx interface{}, input OtpCodeInput, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	return nil, nil
}

func (service *LoginService) RefreshToken(ctx interface{}, token string) (*domain.OAuthSession, error) {
	return nil, nil
}

type SetupUserInput struct {
	Identifier     string
	IdentifierType string
	Username       string
	CountryId      snowflake.ID
	ProfileId      int
	Method         string
}

func (service *LoginService) SetupUser(ctx interface{}, input SetupUserInput, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	return nil, nil
}

func (service *LoginService) GetCurrentSession(ctx interface{}, userId string, accessToken string) (*domain.OAuthSession, error) {
	return nil, nil
}

func (service *LoginService) GetUserPermissions(ctx interface{}, userId snowflake.ID) (*domain.UserPolicy, error) {
	return nil, nil
}

var _ interface{} = (*LoginService)(nil)
