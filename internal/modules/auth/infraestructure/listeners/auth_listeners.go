/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Infraestructure/Listeners
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package listeners

import (
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
)

type AuthListeners struct {
	db     connection.DatabaseConnection
	config *config.Config
}

func NewAuthListeners(db connection.DatabaseConnection, cfg *config.Config) *AuthListeners {
	return &AuthListeners{
		db:     db,
		config: cfg,
	}
}

func (l *AuthListeners) ListenOf(bus domain.EventBus) {
}

var _ interface{} = (*AuthListeners)(nil)
