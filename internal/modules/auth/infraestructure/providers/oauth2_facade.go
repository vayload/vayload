/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Infraestructure/Providers
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package providers

import (
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/domain"
)

type OAuth2Facade struct {
	config *config.Config
}

func NewOAuth2Facade(cfg *config.Config) *OAuth2Facade {
	return &OAuth2Facade{
		config: cfg,
	}
}

func (f *OAuth2Facade) Select(provider domain.OAuth2Provider) (domain.OAuth2Strategy, error) {
	return nil, nil
}

func (f *OAuth2Facade) GetAuthRedirectURL(provider domain.OAuth2Provider, state *domain.OAuth2State) (string, error) {
	return "", nil
}

func (f *OAuth2Facade) ExchangeCode(provider domain.OAuth2Provider, code string) (*domain.OAuthUser, error) {
	return nil, nil
}

func (f *OAuth2Facade) ValidateToken(provider domain.OAuth2Provider, token string) (*domain.OAuthUser, error) {
	return nil, nil
}

var _ domain.OAuth2StrategyFacade = (*OAuth2Facade)(nil)
