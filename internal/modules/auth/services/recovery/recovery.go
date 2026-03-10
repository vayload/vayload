/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package recovery

import (
	"context"

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

// Initiates the password recovery process for a user identified by their email.
func (service *RecoveryService) RequestPasswordRecovery(ctx context.Context, email string) error {
	user, err := service.repository.FindByIdentifier(ctx, email, string(domain.IdentifierTypeEmail))
	if err != nil || user == nil {
		return domain.ErrUserNotFound(err)
	}

	token := service.randomizer.GenerateRandomString(32, true)
	if err := service.repository.SaveRecoveryToken(ctx, user.ID, token); err != nil {
		return err
	}

	go service.EventBus.Publish(ctx, domain.UserPasswordRecoveryRequestedEvent{
		User:  user,
		Token: token,
	})

	return nil
}

// Resets the password for a user using a recovery token.
func (service *RecoveryService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	user, err := service.repository.FindUserByRecoveryToken(ctx, token)
	if err != nil || user == nil {
		return domain.ErrUserNotFound(err)
	}

	hashed := service.strategies.Password.HashPassword(newPassword)
	if err := service.repository.ResetPasswordWithRecoveryToken(ctx, token, hashed); err != nil {
		return err
	}

	go service.EventBus.Publish(ctx, domain.UserPasswordResetCompletedEvent{
		User: user,
	})

	return nil
}
