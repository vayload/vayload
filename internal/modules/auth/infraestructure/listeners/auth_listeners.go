/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package listeners

import (
	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/auth/infraestructure/persistence"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/vayload"
)

type authEvents struct {
	repository domain.AuthLogRepository
	config     *config.Config
	// JobQueue   vayload.JobQueue
}

func NewEventListeners(database connection.DatabaseConnection, config *config.Config) *authEvents {
	repository := persistence.NewLogRepository(database)

	return &authEvents{
		repository: repository,
		config:     config,
		// JobQueue:   queueClient,
	}
}

func (event *authEvents) ListenOf(EventBus vayload.EventBus) {
	// EventBus.Listen(func(payload domain.UserLoggedInEvent) {
	// 	user := payload.User
	// 	authCtx := payload.Context

	// 	err := event.repository.SaveLogin(context.Background(), user.ID, authCtx.IP, authCtx.UserAgent, authCtx.Method, user.Phone, user.Email)
	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserLoggedInEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.UserCreatedEvent) {
	// 	user := payload.User
	// 	code := payload.Code

	// 	_, err := event.queue.Publish(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-created", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{"email"},
	// 			"kind":       "email-verification",
	// 			"email":      user.Email,
	// 			"subject":    "Verificación de correo electrónico",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username": user.Username,
	// 				"code":     code,
	// 			},
	// 			"retriable": false,
	// 			"ttl":       60 * 5,
	// 		},
	// 	})
	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserCreatedEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.UserUpdateCodeEvent) {
	// 	user := payload.User
	// 	code := payload.Code

	// 	_, err := event.JobQueue.Enqueue(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-update", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{"email"},
	// 			"kind":       "email-verification",
	// 			"email":      user.Email,
	// 			"subject":    "Actualización de código de verificación",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username": user.Username,
	// 				"code":     code,
	// 			},
	// 			"retriable": false,
	// 			"ttl":       60 * 5,
	// 		},
	// 	})

	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserUpdateCodeEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.OtpCodeGeneratedEvent) {
	// 	user := payload.User
	// 	code := payload.Code
	// 	channel := payload.Channel

	// 	_, err := event.JobQueue.Enqueue(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-otp", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{channel},
	// 			"kind":       "otp",
	// 			"email":      user.Email,
	// 			"phone":      user.Phone,
	// 			"subject":    "Código OTP generado",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username": user.Username,
	// 				"code":     code,
	// 			},
	// 			"retriable": false,
	// 			"ttl":       60 * 5,
	// 		},
	// 	})

	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "OtpCodeGeneratedEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.UserMagicLinkGeneratedEvent) {
	// 	user := payload.User
	// 	code := payload.Code
	// 	channel := payload.Channel
	// 	expiresIn := payload.ExpiresIn

	// 	_, err := event.JobQueue.Enqueue(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-otp", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{channel},
	// 			"kind":       "magic-link",
	// 			"email":      user.Email,
	// 			"phone":      user.Phone,
	// 			"subject":    "Código de acceso de Vayload",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username":   user.Username,
	// 				"code":       code,
	// 				"magic_link": fmt.Sprintf("%s/?mglk=%s", event.config.AppLink, code),
	// 				"expires_in": expiresIn.Seconds(),
	// 			},
	// 			"retriable": false,
	// 			"ttl":       expiresIn.Seconds(),
	// 		},
	// 	})

	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserMagicLinkGeneratedEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.UserEmailVerifiedEvent) {
	// 	user := payload.User

	// 	_, err := event.JobQueue.Enqueue(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-email-verified", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{"email"},
	// 			"kind":       "welcome",
	// 			"email":      user.Email,
	// 			"subject":    "Bienvenido a Vayload",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username": user.Username,
	// 			},
	// 			"retriable": false,
	// 			"ttl":       60 * 5,
	// 		},
	// 	})

	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserEmailVerifiedEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.UserPasswordRecoveryRequestedEvent) {
	// 	user := payload.User
	// 	token := payload.Token

	// 	_, err := event.JobQueue.Enqueue(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-password-reset", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{"email"},
	// 			"kind":       "password-reset",
	// 			"email":      user.Email,
	// 			"subject":    "Recuperación de contraseña",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username": user.Username,
	// 				"token":    token,
	// 				"link":     fmt.Sprintf("https://%s/reset-password?tk=%s", event.config.AppLink, token),
	// 			},
	// 			"retriable": false,
	// 			"ttl":       60 * 10,
	// 		},
	// 	})
	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserPasswordRecoveryRequestedEvent", "userId": user.ID})
	// 	}
	// })

	// EventBus.Listen(func(payload domain.UserEmailChangeRequestedEvent) {
	// 	user := payload.User
	// 	_, err := event.JobQueue.Enqueue(context.Background(), queue.NOTIFICATIONS_QUEUE, queue.TaskOptions{
	// 		Name: fmt.Sprintf("user-%s-email-change", user.ID.String()),
	// 		Payload: map[string]any{
	// 			"user_id":    user.ID.String(),
	// 			"country_id": user.CountryID,
	// 			"channels":   []string{"email"},
	// 			"kind":       "email-change",
	// 			"email":      user.Email,
	// 			"subject":    "Confirmación de cambio de correo",
	// 			"language":   "es",
	// 			"data": map[string]any{
	// 				"username":      user.Username,
	// 				"current_token": payload.CurrentToken,
	// 				"new_token":     payload.NewToken,
	// 				"new_email":     payload.NewEmail,
	// 			},
	// 			"retriable": false,
	// 			"ttl":       60 * 10,
	// 		},
	// 	})
	// 	if err != nil {
	// 		logger.E(err, logger.Fields{"context": "UserEmailChangeRequestedEvent", "userId": user.ID})
	// 	}
	// })
}
