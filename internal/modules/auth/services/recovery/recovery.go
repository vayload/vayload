/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Services/Recovery
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package recovery

import (
	"github.com/vayload/vayload/internal/modules/auth/domain"
)

const SERVICE_KEY = "auth.recovery"

type RecoveryService struct {
	repository domain.AuthRepository
	strategies *RecoveryStrategies
	randomizer domain.SecureRandomizer
	EventBus   domain.EventBus
}

type RecoveryStrategies struct {
	Password domain.PasswordStrategy
}

func NewRecoveryService(repository domain.AuthRepository, strategies *RecoveryStrategies, randomizer domain.SecureRandomizer, EventBus domain.EventBus) *RecoveryService {
	return &RecoveryService{
		repository: repository,
		strategies: strategies,
		randomizer: randomizer,
		EventBus:   EventBus,
	}
}

func (service *RecoveryService) RequestPasswordRecovery(ctx interface{}, email string) error {
	return nil
}

func (service *RecoveryService) ResetPassword(ctx interface{}, token string, newPassword string) error {
	return nil
}

var _ interface{} = (*RecoveryService)(nil)
