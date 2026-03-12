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

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/shared/errors"
	"github.com/vayload/vayload/internal/shared/snowflake"
	"github.com/vayload/vayload/pkg/collect"
	"github.com/vayload/vayload/pkg/logger"
)

type AuthStrategies struct {
	OAuth2   domain.OAuth2StrategyFacade
	Password domain.PasswordStrategy
	OtpCode  domain.OtpStrategy
}

type LoginService struct {
	repository   domain.AuthRepository
	tokenManager domain.UserTokenManager
	randomizer   domain.SecureRandomizer
	EventBus     domain.EventBus
	strategies   *AuthStrategies
}

func NewLoginService(
	repository domain.AuthRepository,
	tokenManager domain.UserTokenManager,
	randomizer domain.SecureRandomizer,
	EventBus domain.EventBus,
	authStrategies *AuthStrategies,
) *LoginService {

	return &LoginService{
		repository:   repository,
		tokenManager: tokenManager,
		randomizer:   randomizer,
		EventBus:     EventBus,
		strategies:   authStrategies,
	}
}

func (service *LoginService) IdentifyUserForFactor(ctx context.Context, input IdentifyInput) (*AuthStepCredentials, error) {
	user, err := service.repository.FindUserIdentity(ctx, input.UserIdentifier, input.IdentifierType)
	if user == nil || err != nil {
		return nil, domain.ErrUserNotFound(err)
	}

	credentials := &AuthStepCredentials{
		UserIdentifier: input.UserIdentifier,
		Factor:         input.ForFactor,
		ClientType:     input.ClientType,
	}

	// if the factor is otp required notifications provider for send otp codes
	// send client this informacion for better UX and UI
	if input.ForFactor == domain.OtpStrategyType {
		otpProviders := &domain.OtpProvider{
			Email:    []string{},
			SMS:      []string{},
			WhatsApp: []string{},
		}

		if user.CountryID != nil {
			countryId := snowflake.ID(*user.CountryID)
			otpProvider, err := service.repository.FindCountryOtpProviders(context.Background(), countryId)
			if err != nil {
				// Only log the error, because it's not critical and needs to be handled gracefully
				logger.E(err, logger.Fields{"context": "IdentifyUserForFactor", "countryId": countryId})
			} else if otpProvider != nil {
				otpProviders = &domain.OtpProvider{
					Email:    collect.Filter(otpProvider.Email, func(item string) bool { return item != "" && item != "null" }),
					SMS:      collect.Filter(otpProvider.SMS, func(item string) bool { return item != "" && item != "null" }),
					WhatsApp: collect.Filter(otpProvider.WhatsApp, func(item string) bool { return item != "" && item != "null" }),
				}
			}
		}

		credentials.OtpProviders = otpProviders
		credentials.Destinations = map[string]any{
			"email": MaskIdentity(user.Email, "email"),
		}
		if user.Phone != nil {
			credentials.Destinations["phone"] = MaskIdentity(*user.Phone, "phone")
		}

		// Include a default if not foudn any provider for this country
		switch domain.IdentifierType(input.IdentifierType) {
		case domain.IdentifierTypeEmail:
			if len(otpProviders.Email) == 0 {
				otpProviders.Email = []string{"sendia"}
			}
		case domain.IdentifierTypePhone:
			if len(otpProviders.SMS) == 0 {
				otpProviders.SMS = []string{}
			}
			if len(otpProviders.WhatsApp) == 0 {
				otpProviders.WhatsApp = []string{"vayload-ws"}
			}
		}
	}

	return credentials, nil
}

func (service *LoginService) GenerateOtpCode(ctx context.Context, input OtpCodeGenInput) error {
	identifierType := detectIdentifierType(input.Identifier)

	user, err := service.repository.FindByIdentifier(ctx, input.Identifier, identifierType)
	if err != nil || user == nil {
		return domain.ErrUserNotFound(err)
	}

	otpCode := service.strategies.OtpCode.GenerateOtpCode()

	// Send OTP code to user
	if err := service.repository.SaveOtpCode(context.Background(), user.ID, otpCode); err != nil {
		return fmt.Errorf("saving OTP code: %w", err)
	}

	go service.EventBus.Publish(ctx, domain.OtpCodeGeneratedEvent{
		User:    user,
		Code:    otpCode,
		Channel: input.Channel,
	})

	return nil
}

// =========================== LOGIN METHODS =====================
func (service *LoginService) LoginWithPassword(ctx context.Context, input LoginInput, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	identifierType := detectIdentifierType(input.Identifier)

	user, err := service.repository.FindByIdentifier(ctx, input.Identifier, identifierType)
	if err != nil {
		return nil, domain.ErrInvalidCredentials(err)
	}

	if user.Password == nil {
		return nil, domain.ErrInvalidCredentials(fmt.Errorf("user password is not set"))
	}

	passing, algo := service.strategies.Password.VerboseVerifyPassword(input.Password, *user.Password)
	// if not passing, reject login attempt
	if !passing {
		return nil, domain.ErrInvalidCredentials(fmt.Errorf("invalid password"))
	}

	// migrate to new hashing algorithm
	// previous passwords build with wordpress hashing algorithm
	if algo != "scrypt" {
		newPassword := service.strategies.Password.HashPassword(input.Password)
		if len(newPassword) > 0 {
			passUpdateErr := service.repository.UpdatePassword(context.Background(), user.ID, newPassword)
			if passUpdateErr != nil {
				// Only log this error, because migrate to new hashing algorithm is not critical
				// In next version, we can remove this logic (when user passwords are already migrated)
				logger.E(passUpdateErr, logger.Fields{"context": "LoginWithPassword", "action": "migrate password hashing"})
			}
		}
	}

	token, err := service.tokenManager.GenerateJwtTokenWithRefresh(&domain.AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		ClientId:  *user.ClientId,
		CountryId: user.CountryID,
		Meta: map[string]any{
			"avatar_url": user.AvatarURL,
		},
	})

	if err != nil {
		return nil, domain.ErrJwtTokenGenerationFailed(err)
	}

	// For extenal event logging
	go service.EventBus.Publish(ctx, domain.UserLoggedInEvent{
		User: user,
		Context: &domain.AuthContext{
			IP:        authContext.IP,
			UserAgent: authContext.UserAgent,
			Method:    "password",
		},
	})

	session := domain.NewOAuthSessionFromToken(token)

	return session, nil
}

func (service *LoginService) GetOAuth2URL(ctx context.Context, provider domain.OAuth2Provider, payload domain.OAuth2State) (string, error) {
	return service.strategies.OAuth2.GetAuthRedirectURL(provider, &payload)
}

func (service *LoginService) LoginWithOAuth2(ctx context.Context, provider domain.OAuth2Provider, code string, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	oauth, err := service.strategies.OAuth2.ExchangeCode(provider, code)
	if err != nil {
		return nil, domain.ErrInvalidCredentials(err)
	}

	user := domain.NewUser(oauth.FirstName, oauth.Email, nil, domain.PatientRole)

	rawUser, _ := service.repository.FindByIdentifier(ctx, oauth.Email, string(domain.IdentifierTypeEmail))

	// If user exists, validate account status (banned/deleted)
	if rawUser != nil {
		user.ID = rawUser.ID
		user.Role = rawUser.Role
		user.ClientId = rawUser.ClientId
		user.CountryID = rawUser.CountryID
	} else {
		// If user not found, create a new user
		user.Username = oauth.FirstName
		user.LastName = &oauth.LastName
		user.AvatarURL = &oauth.AvatarURL

		settings := &domain.UserSettings{
			Language: "es",
			Notifications: domain.UserNotificationSettings{
				Email: true,
			},
		}

		if createErr := service.repository.CreateUserWithSettings(ctx, user, settings); createErr != nil {
			return nil, createErr
		}
	}

	// If user has no clientId, create a new client and bind authorization
	if user.ClientId == nil {
		// client, policyErr := service.authorization.SetupWithEmail(user.Email, int(service.authorization.GetFreeProfileId()))
		// if policyErr != nil {
		// 	return nil, policyErr
		// }

		// user.ClientId = &client.Id
		// user.ProfileId = &client.ProfileId

		// binding := &domain.AuthorizationBinding{
		// 	UserId:    user.ID,
		// 	ProfileId: client.ProfileId,
		// 	ClientId:  client.Id,
		// 	Signature: client.Signature,
		// }
		// meta := &domain.UserMeta{
		// 	EmailVerified:     true,
		// 	ConfirmationToken: "",
		// 	VerificationCode:  "",
		// }

		// if authErr := service.repository.BindAuthorization(ctx, user.ID, binding, meta); authErr != nil {
		// 	return nil, authErr
		// }
	}

	token, err := service.tokenManager.GenerateJwtTokenWithRefresh(&domain.AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		ClientId:  *user.ClientId,
		CountryId: user.CountryID,
		Meta: map[string]any{
			"avatar_url": user.AvatarURL,
		},
	})

	if err != nil {
		return nil, domain.ErrJwtTokenGenerationFailed(err)
	}

	go service.EventBus.Publish(ctx, domain.UserLoggedInEvent{
		User: user,
		Context: &domain.AuthContext{
			IP:        authContext.IP,
			UserAgent: authContext.UserAgent,
			Method:    fmt.Sprintf("oauth2.%s", provider),
		},
	})

	session := domain.NewOAuthSessionFromToken(token)

	return session, nil
}

func (service *LoginService) LoginWithOtpCode(ctx context.Context, input OtpCodeInput, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	identifierType := detectIdentifierType(input.Identifier)
	user, err := service.repository.FindByIdentifier(ctx, input.Identifier, identifierType)
	if err != nil || user == nil {
		return nil, domain.ErrInvalidCredentials(err)
	}

	if !service.strategies.OtpCode.CompareOtpCode(input.Code, user.OTPCode) {
		return nil, domain.ErrInvalidCredentials(fmt.Errorf("invalid OTP code"))
	}

	jwtToken, err := service.tokenManager.GenerateJwtTokenWithRefresh(&domain.AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		ClientId:  *user.ClientId,
		CountryId: user.CountryID,
		Meta: map[string]any{
			"avatar_url": user.AvatarURL,
		},
	})

	if err != nil {
		return nil, domain.ErrJwtTokenGenerationFailed(err)
	}

	go service.EventBus.Publish(ctx, domain.UserLoggedInEvent{
		User: user,
		Context: &domain.AuthContext{
			IP:        authContext.IP,
			UserAgent: authContext.UserAgent,
			Method:    "otp",
		},
	})

	session := domain.NewOAuthSessionFromToken(jwtToken)

	return session, nil
}

func (service *LoginService) RefreshToken(ctx context.Context, token string) (*domain.OAuthSession, error) {
	rawUser, err := service.tokenManager.ValidateRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}

	user, err := service.repository.FindByIdentifier(ctx, rawUser.Email, "email")
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	jwtToken, err := service.tokenManager.GenerateJwtToken(&domain.AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		ClientId:  *user.ClientId,
		CountryId: user.CountryID,
		Meta: map[string]any{
			"avatar_url": user.AvatarURL,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	return &domain.OAuthSession{
		AccessToken: jwtToken.GetAccessToken(),
		ExpiresAt:   jwtToken.GetExpiresAt(),
		ExpiresIn:   jwtToken.GetExpiresIn(),
		TokenType:   "Bearer",
	}, nil
}

func (service *LoginService) SetupUser(ctx context.Context, input SetupUserInput, authContext domain.AuthContext) (*domain.OAuthSession, error) {
	user, err := service.repository.FindByIdentifier(ctx, input.Identifier, input.IdentifierType)

	// Reject if have error and this error is diferent to not found
	if err != nil && !errors.Is(err, domain.ErrNotResults) {
		return nil, fmt.Errorf("finding user: %w", err)
	}

	var policy *domain.UserPolicy
	var policyErr error

	// Perform user creation when not exists in the system
	// Email is required for user creation
	if user == nil {
		// user = domain.NewUser(input.Username, "", nil, domain.PatientRole)
		// user.CountryID = &input.CountryId
		// user.AuthType = input.Method

		// switch domain.IdentifierType(input.IdentifierType) {
		// case domain.IdentifierTypeEmail:
		// 	user.Email = input.Identifier

		// 	// Perform authorization setup with email
		// 	policy, policyErr = service.authorization.SetupWithEmail(input.Identifier, input.ProfileId)
		// case domain.IdentifierTypePhone:
		// 	user.Email = fmt.Sprintf("%s@users.vayload.com", strings.ReplaceAll(input.Identifier, "+", "")) // Dummy email for phone-based users
		// 	user.Phone = &input.Identifier

		// 	// Perform authorization setup with phone
		// 	policy, policyErr = service.authorization.SetupWithPhone(input.Identifier, input.ProfileId)
		// }

		// If error occurred while setting up user policy
		// Prevent user creation if policy setup fails
		if policyErr != nil {
			return nil, fmt.Errorf("setting up user policy: %w", policyErr)
		}

		clientId := policy.GetClientId()
		profileId := policy.GetProfileId()
		user.ClientId = &clientId
		user.ProfileId = &profileId

		// Create one user because not founded with given identifier
		settings := &domain.UserSettings{
			Language:      "es",
			Notifications: domain.UserNotificationSettings{},
		}

		switch domain.IdentifierType(input.IdentifierType) {
		case domain.IdentifierTypeEmail:
			settings.Notifications.Email = true
		case domain.IdentifierTypePhone:
			settings.Notifications.SMS = true
			settings.Notifications.WhatsApp = true
		}

		if createErr := service.repository.CreateUserWithSettings(ctx, user, settings); createErr != nil {
			return nil, fmt.Errorf("creating user: %w", createErr)
		}
	}

	// When policy is nil, it means user already exists
	if policy == nil {
		// policy, policyErr = service.authorization.GetAuthorized(user.ID, *user.ClientId)
		// if policyErr != nil {
		// 	return nil, fmt.Errorf("getting user policy: %w", policyErr)
		// }
	}

	token, err := service.tokenManager.GenerateJwtTokenWithRefresh(&domain.AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		ClientId:  *user.ClientId,
		CountryId: user.CountryID,
	})

	if err != nil {
		return nil, domain.ErrJwtTokenGenerationFailed(err)
	}

	go service.EventBus.Publish(ctx, domain.UserLoggedInEvent{
		User: user,
		Context: &domain.AuthContext{
			IP:        authContext.IP,
			UserAgent: authContext.UserAgent,
			Method:    input.Method,
		},
	})

	session := domain.NewOAuthSession(token, user, policy, nil)

	return session, nil
}

func (service *LoginService) GetCurrentSession(ctx context.Context, userId string, accessToken string) (*domain.OAuthSession, error) {
	user, err := service.repository.FindByIdentifier(ctx, userId, string(domain.IdentifierUserId))
	if err != nil && !errors.Is(err, domain.ErrNotResults) {
		return nil, fmt.Errorf("finding user: %w", err)
	}

	// Perform user creation when not exists in the system
	if user == nil {
		return nil, domain.ErrUserNotFound(fmt.Errorf("user not found with ID: %s", userId))
	}

	// When policy is nil, it means user already exists
	// policy, policyErr := service.authorization.GetAuthorized(user.ID, *user.ClientId)
	// if policyErr != nil {
	// 	return nil, fmt.Errorf("getting user policy: %w", policyErr)
	// }

	session := &domain.OAuthSession{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		User:        user,
		// Policies:    policy,
		Meta: nil,
	}

	return session, nil
}

func (service *LoginService) GetUserPermissions(ctx context.Context, userId snowflake.ID) (*domain.UserPolicy, error) {
	user, err := service.repository.FindByIdentifier(ctx, userId.String(), string(domain.IdentifierUserId))
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user.ClientId == nil {
		return nil, domain.ErrInvalidCredentials(nil)
	}

	// policy, err := service.authorization.GetAuthorized(userId, *user.ClientId)
	// if err != nil {
	// 	return nil, fmt.Errorf("getting user policy: %w", err)
	// }

	// return policy, nil
	return nil, nil
}
