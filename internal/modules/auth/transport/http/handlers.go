/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Transport/Http
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package http

import (
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/services/analytics"
	"github.com/vayload/vayload/internal/modules/auth/services/login"
	"github.com/vayload/vayload/internal/modules/auth/services/recovery"
	"github.com/vayload/vayload/internal/modules/auth/services/registration"
	"github.com/vayload/vayload/internal/vayload"
)

type HttpServices struct {
	LoginService     *login.LoginService
	RegisterService  *registration.RegisterService
	RecoveryService  *recovery.RecoveryService
	AnalyticsService *analytics.AnalyticsService
}

type AuthHttpHandler struct {
	config    *config.Config
	container vayload.Container
	services  HttpServices
}

func NewHttpHandler(cfg *config.Config, container vayload.Container, services HttpServices) *AuthHttpHandler {
	return &AuthHttpHandler{
		config:    cfg,
		container: container,
		services:  services,
	}
}

func (h *AuthHttpHandler) RegisterRoutes() {
}

var _ interface{} = (*AuthHttpHandler)(nil)
