/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth Service
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package auth

import (
	"context"
	"fmt"

	"github.com/vayload/vayload/config"
	"github.com/vayload/vayload/internal/kernel"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/auth/infraestructure/listeners"
	"github.com/vayload/vayload/internal/modules/auth/infraestructure/persistence"
	"github.com/vayload/vayload/internal/modules/auth/infraestructure/providers"
	"github.com/vayload/vayload/internal/modules/auth/services/analytics"
	"github.com/vayload/vayload/internal/modules/auth/services/login"
	"github.com/vayload/vayload/internal/modules/auth/services/recovery"
	"github.com/vayload/vayload/internal/modules/auth/services/registration"
	auth_http "github.com/vayload/vayload/internal/modules/auth/transport/http"
	"github.com/vayload/vayload/internal/modules/database"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/security"
	"github.com/vayload/vayload/internal/vayload"
)

const (
	ServiceName = "auth"

	JWTManagerKey           = "auth.jwt_manager"
	LoginServiceKey         = "auth.login_service"
	RegisterServiceKey      = "auth.register_service"
	RecoveryServiceKey      = "auth.recovery_service"
	AnalyticsServiceKey     = "auth.analytics_service"
	AuthorizationServiceKey = "auth.authorization_service"
)

type AuthService struct {
	kernel.BaseService

	config      *config.Config
	db          connection.DatabaseConnection
	httpHandler *auth_http.AuthHttpHandler
}

func NewAuthService(cfg *config.Config) *AuthService {
	deps := []vayload.ServiceName{
		vayload.ServiceDatabaseName,
	}

	return &AuthService{
		BaseService: kernel.NewBaseService(vayload.ServiceAuthName, true, deps...),
		config:      cfg,
	}
}

func (s *AuthService) Bootstrap(ctx context.Context, args map[string]any, reply *map[string]any) error {
	container := s.Container()
	if container == nil {
		return fmt.Errorf("container not provided for auth service")
	}

	db, err := kernel.MapTo[connection.DatabaseConnection](container, database.DATABASE_CONNECTION)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	s.db = db

	userRepository := persistence.NewUserRepository(db)
	analyticsRepository := persistence.NewAnalyticsRepository(db)

	tokenManager := security.NewJwtManager(security.JwtConfig{
		PublicKeyBytes:     s.config.Security.JwtPublicKey,
		PrivateKeyBytes:    s.config.Security.JwtPrivateKey,
		ExpireAccessToken:  int64(s.config.Security.JwtExpirationTime),
		ExpireRefreshToken: int64(s.config.Security.JwtRefreshExpirationDays * 24),
	})
	container.SetInstance(JWTManagerKey, tokenManager)

	hashing := security.NewHasher()
	randomizer := security.NewRandomizer()
	oAuth2Facade := providers.NewOAuth2Facade(s.config)

	authStrategies := &login.AuthStrategies{
		OAuth2:   oAuth2Facade,
		Password: hashing,
		OtpCode:  security.NewOtpCodeStrategy(),
	}

	registrationStrategies := &registration.RegistrationStrategies{
		Password: hashing,
	}

	eventBus := s.EventBus()

	loginService := login.NewLoginService(userRepository, tokenManager, randomizer, eventBus, authStrategies)
	registerService := registration.NewRegisterService(userRepository, tokenManager, registrationStrategies, randomizer, eventBus)
	recoveryService := recovery.NewRecoveryService(userRepository, &recovery.RecoveryStrategies{Password: hashing}, randomizer, eventBus)
	analyticsService := analytics.NewAnalyticsService(analyticsRepository)

	container.SetInstance(LoginServiceKey, loginService)
	container.SetInstance(RegisterServiceKey, registerService)
	container.SetInstance(RecoveryServiceKey, recoveryService)
	container.SetInstance(AnalyticsServiceKey, analyticsService)

	s.httpHandler = auth_http.NewHttpHandler(s.config, container, auth_http.HttpServices{
		LoginService:     loginService,
		RegisterService:  registerService,
		RecoveryService:  recoveryService,
		AnalyticsService: analyticsService,
	})

	authListeners := listeners.NewEventListeners(db, s.config)
	authListeners.ListenOf(eventBus)

	_ = domain.ErrContext

	return nil
}

func (s *AuthService) Shutdown(ctx context.Context) error {
	return nil
}

func (s *AuthService) HttpRoutes() []vayload.HttpRoutesGroup {
	if s.httpHandler == nil {
		return []vayload.HttpRoutesGroup{}
	}

	return s.httpHandler.HttpRoutes()
}

var _ vayload.Service = (*AuthService)(nil)
