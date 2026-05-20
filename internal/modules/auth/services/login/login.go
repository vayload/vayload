/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package login

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/shared/errors"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type AuthStrategies struct {
	OAuth2   domain.OAuth2StrategyFacade
	Password domain.PasswordStrategy
	OtpCode  domain.OtpStrategy
}

type LoginService struct {
	userRepo          domain.AuthRepository
	rbacRepo          domain.RbacRepository
	sessionRepo       domain.SessionRepository
	refreshRepository domain.RefreshTokenRepository
	tokenManager      domain.UserTokenManager
	randomizer        domain.SecureRandomizer
	eventBus          domain.EventBus
	strategies        *AuthStrategies
}

func NewLoginService(
	userRepo domain.AuthRepository,
	rbacRepo domain.RbacRepository,
	sessionRepo domain.SessionRepository,
	refreshRepository domain.RefreshTokenRepository,
	tokenManager domain.UserTokenManager,
	randomizer domain.SecureRandomizer,
	eventBus domain.EventBus,
	authStrategies *AuthStrategies,
) *LoginService {
	return &LoginService{
		userRepo:          userRepo,
		rbacRepo:          rbacRepo,
		sessionRepo:       sessionRepo,
		refreshRepository: refreshRepository,
		tokenManager:      tokenManager,
		randomizer:        randomizer,
		eventBus:          eventBus,
		strategies:        authStrategies,
	}
}

func (s *LoginService) IdentifyUserForFactor(ctx context.Context, input IdentifyInput) (*AuthStepCredentials, error) {
	user, err := s.userRepo.FindByIdentifier(ctx, input.UserIdentifier, domain.IdentifierType(input.IdentifierType))
	if user == nil || err != nil {
		return nil, domain.ErrUserNotFound(err)
	}

	credentials := &AuthStepCredentials{
		UserIdentifier: input.UserIdentifier,
		Factor:         input.ForFactor,
		ClientType:     input.ClientType,
	}

	if input.ForFactor == domain.OtpStrategyType {
		otpProviders := &domain.OtpProvider{
			Email:    []string{"sendia"},
			SMS:      []string{},
			WhatsApp: []string{"vayload-ws"},
		}
		credentials.OtpProviders = otpProviders
		credentials.Destinations = map[string]any{
			"email": MaskIdentity(user.Email, "email"),
		}
		if user.Phone != nil {
			credentials.Destinations["phone"] = MaskIdentity(*user.Phone, "phone")
		}
	}

	return credentials, nil
}

func (s *LoginService) GenerateOtpCode(ctx context.Context, input OtpCodeGenInput) error {
	user, err := s.userRepo.FindByIdentifier(ctx, input.Identifier, domain.IdentifierType(detectIdentifierType(input.Identifier)))
	if err != nil || user == nil {
		return domain.ErrUserNotFound(err)
	}

	otpCode := s.strategies.OtpCode.GenerateOtpCode()
	if err := s.userRepo.SaveOtpCode(ctx, user.ID, otpCode); err != nil {
		return fmt.Errorf("saving OTP code: %w", err)
	}

	go s.eventBus.Publish(ctx, domain.OtpCodeGeneratedEvent{
		User:    user,
		Code:    otpCode,
		Channel: input.Channel,
	})

	return nil
}

func (s *LoginService) LoginWithPassword(ctx context.Context, input LoginInput, authCtx domain.AuthContext) (*domain.OAuthSession, error) {
	user, err := s.userRepo.FindByIdentifier(ctx, input.Identifier, domain.IdentifierType(detectIdentifierType(input.Identifier)))
	if err != nil {
		return nil, domain.ErrInvalidCredentials(err)
	}

	if user.PasswordHash == nil {
		return nil, domain.ErrInvalidCredentials(fmt.Errorf("user has no password set"))
	}

	passing, algo := s.strategies.Password.VerboseVerifyPassword(input.Password, *user.PasswordHash)
	if !passing {
		return nil, domain.ErrInvalidCredentials(fmt.Errorf("invalid password"))
	}

	if algo != "scrypt" {
		newHash := s.strategies.Password.HashPassword(input.Password)
		if newHash != "" {
			_ = s.userRepo.UpdatePassword(ctx, user.ID, newHash)
		}
	}

	return s.createSession(ctx, user, authCtx)
}

func (s *LoginService) LoginWithOtpCode(ctx context.Context, input OtpCodeInput, authCtx domain.AuthContext) (*domain.OAuthSession, error) {
	user, err := s.userRepo.FindByIdentifier(ctx, input.Identifier, domain.IdentifierType(detectIdentifierType(input.Identifier)))
	if err != nil || user == nil {
		return nil, domain.ErrInvalidCredentials(err)
	}

	if user.OTPCode == nil || !s.strategies.OtpCode.CompareOtpCode(input.Code, *user.OTPCode) {
		return nil, domain.ErrInvalidCredentials(fmt.Errorf("invalid OTP code"))
	}

	return s.createSession(ctx, user, authCtx)
}

func (s *LoginService) LoginWithOAuth2(ctx context.Context, provider domain.OAuth2Provider, code string, authCtx domain.AuthContext) (*domain.OAuthSession, error) {
	oauth, err := s.strategies.OAuth2.ExchangeCode(provider, code)
	if err != nil {
		return nil, domain.ErrInvalidCredentials(err)
	}

	user, err := s.userRepo.FindByIdentifier(ctx, oauth.Email, domain.IdentifierTypeEmail)
	if err != nil {
		// Auto-register user if not found
		user = &domain.User{
			ID:        snowflake.Node.Generate(),
			Email:     oauth.Email,
			FirstName: &oauth.FirstName,
			LastName:  &oauth.LastName,
			AvatarURL: &oauth.AvatarURL,
			CreatedAt: time.Now().UTC(),
		}
		// TODO: Implement user creation logic here if needed for OAuth
	}

	return s.createSession(ctx, user, authCtx)
}

func (s *LoginService) GetOAuth2URL(ctx context.Context, provider domain.OAuth2Provider, payload domain.OAuth2State) (string, error) {
	return s.strategies.OAuth2.GetAuthRedirectURL(provider, &payload)
}

func (s *LoginService) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.OAuthSession, error) {
	storedToken, err := s.refreshRepository.FindByHash(ctx, refreshTokenStr)
	if err != nil {
		return nil, errors.NotFound("refresh token not found").Cause(err)
	}

	if storedToken.RevokedAt != nil {
		return nil, errors.Unauthorized("refresh token revoked")
	}
	if storedToken.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.Unauthorized("refresh token expired")
	}

	user, err := s.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, errors.NotFound("user not found").Cause(err)
	}

	token, err := s.tokenManager.GenerateJwtToken(user)
	if err != nil {
		return nil, errors.Internal("failed to generate token").Cause(err)
	}

	_ = s.sessionRepo.Update(ctx, &domain.Session{
		ID:         storedToken.SessionID,
		LastSeenAt: time.Now().UTC(),
	})

	return &domain.OAuthSession{
		AccessToken: token.GetAccessToken(),
		TokenType:   "Bearer",
		ExpiresIn:   token.GetExpiresIn(),
		ExpiresAt:   token.GetExpiresAt(),
		User:        user,
	}, nil
}

func (s *LoginService) createSession(ctx context.Context, user *domain.User, authCtx domain.AuthContext) (*domain.OAuthSession, error) {
	token, err := s.tokenManager.GenerateJwtTokenWithRefresh(user)
	if err != nil {
		return nil, domain.ErrJwtTokenGenerationFailed(err)
	}

	sessionID := strings.ReplaceAll(s.randomizer.GenerateUUID(), "-", "")
	session := &domain.Session{
		ID:         sessionID,
		UserID:     user.ID,
		IPAddress:  authCtx.IP,
		UserAgent:  authCtx.UserAgent,
		LastSeenAt: time.Now().UTC(),
		ExpiresAt:  token.GetRefreshTokenExpiresAt(),
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	refreshToken := &domain.RefreshToken{
		ID:        strings.ReplaceAll(s.randomizer.GenerateUUID(), "-", ""),
		TokenHash: token.GetRefreshToken(),
		UserID:    user.ID,
		SessionID: sessionID,
		ExpiresAt: token.GetRefreshTokenExpiresAt(),
		CreatedAt: time.Now().UTC(),
	}
	if err := s.refreshRepository.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	_ = s.userRepo.UpdateLastSignIn(ctx, user.ID, authCtx.IP, authCtx.UserAgent)

	go s.eventBus.Publish(ctx, domain.UserLoggedInEvent{
		User:    user,
		Context: &authCtx,
	})

	return domain.NewOAuthSession(token, user, user.Metadata), nil
}

func (s *LoginService) GetCurrentSession(ctx context.Context, userId string, accessToken string) (*domain.OAuthSession, error) {
	id, _ := snowflake.ParseString(userId)
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrUserNotFound(err)
	}

	return &domain.OAuthSession{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User:        user,
	}, nil
}

func (s *LoginService) GetUserPermissions(ctx context.Context, userId snowflake.ID) ([]string, error) {
	return []string{"read", "write"}, nil
}
