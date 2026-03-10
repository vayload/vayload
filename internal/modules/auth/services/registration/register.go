/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package registration

import (
	"context"
	"crypto/subtle"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/auth/services/authorization"
	"github.com/vayload/vayload/internal/shared/errors"
)

type RegistrationStrategies struct {
	Password domain.PasswordStrategy
}

type RegisterService struct {
	repository    domain.AuthRepository
	tokenManager  domain.UserTokenManager
	strategies    *RegistrationStrategies
	randomizer    domain.SecureRandomizer
	EventBus      domain.EventBus
	authorization authorization.AuthorizationService
}

func NewRegisterService(
	repository domain.AuthRepository,
	tokenManager domain.UserTokenManager,
	strategies *RegistrationStrategies,
	randomizer domain.SecureRandomizer,
	EventBus domain.EventBus,
	authorization authorization.AuthorizationService,
) *RegisterService {
	return &RegisterService{
		repository:    repository,
		tokenManager:  tokenManager,
		strategies:    strategies,
		randomizer:    randomizer,
		EventBus:      EventBus,
		authorization: authorization,
	}
}

// Handles the registration of a new user.
func (service *RegisterService) RegisterUser(ctx context.Context, input RegisterInput) (*domain.User, error) {
	existingUser, err := service.repository.FindByIdentifier(ctx, input.Email, "email")

	// Reject when email already exists
	if existingUser != nil || !errors.Is(err, domain.ErrEmptyResultSet) {
		return nil, domain.ErrEmailAlreadyExists(err)
	}

	// Create new user with required fields
	user := domain.NewUser(input.Username, input.Email, &input.Password, domain.PatientRole)
	hashedPassword := service.strategies.Password.HashPassword(*user.Password)
	user.SetPassword(&hashedPassword)
	if user.Password == nil {
		return nil, errors.New("failed to hash password")
	}

	code := service.randomizer.GenerateRandomNumericCode(6)
	token := service.randomizer.GenerateRandomString(32, true)

	// Create user with pending verification
	user, err = service.repository.CreateUserWithCode(ctx, user, token, code)
	if err != nil {
		return nil, err
	}

	go service.EventBus.Publish(ctx, domain.UserCreatedEvent{
		User: user,
		Code: code,
	})

	return user, nil
}

// Validates the registration process by checking the verification code and setting up user authorization.
func (service *RegisterService) ValidateRegister(ctx context.Context, input RegisterValidationInput) (*domain.OAuthSession, error) {
	codes, err := service.repository.FindCodesByIdentifier(ctx, input.Identifier, input.Type, "email")
	if err != nil || codes == nil {
		return nil, domain.ErrUserNotFound(err)
	}

	if subtle.ConstantTimeCompare([]byte(*codes.VerificationCode), []byte(input.Code)) != 1 {
		return nil, domain.ErrInvalidOtpCode(nil)
	}

	user, err := service.repository.FindByIdentifier(ctx, input.Identifier, input.Type)
	if err != nil || user == nil {
		return nil, domain.ErrUserNotFound(err)
	}

	var policy *domain.UserPolicy
	var policyErr error

	// upsert user in authorization service
	switch domain.IdentifierType(input.Type) {
	case domain.IdentifierTypeEmail:
		policy, policyErr = service.authorization.SetupWithEmail(user.Email, int(service.authorization.GetFreeProfileId()))
	case domain.IdentifierTypePhone:
		policy, policyErr = service.authorization.SetupWithPhone(*user.Phone, int(service.authorization.GetFreeProfileId()))
	}

	// If error occurred while setting up user policy
	if policyErr != nil {
		return nil, policyErr
	}

	binding := &domain.AuthorizationBinding{
		UserId:    user.ID,
		ClientId:  policy.GetClientId(),
		Signature: policy.GetSignature(),
		ProfileId: policy.GetProfileId(),
	}

	err = service.repository.BindAuthorization(ctx, user.ID, binding, &domain.UserMeta{
		ConfirmationToken: "",
		VerificationCode:  "",
		EmailVerified:     true,
	})

	if err != nil {
		return nil, err
	}

	go service.EventBus.Publish(ctx, domain.UserEmailVerifiedEvent{
		User: user,
	})

	jwtToken, err := service.tokenManager.GenerateJwtTokenWithRefresh(&domain.AuthUser{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CountryId: user.CountryID,
		ClientId:  binding.ClientId,
	})
	if err != nil {
		return nil, err
	}

	return domain.NewOAuthSessionFromToken(jwtToken), nil
}

// Sends a verification code to the user's email.
func (service *RegisterService) SendEmailVerificationCode(ctx context.Context, dto RegisterValidationInput) error {
	user, err := service.repository.FindByIdentifier(ctx, dto.Identifier, dto.Type)
	if err != nil || user == nil {
		if errors.Is(err, domain.ErrEmptyResultSet) {
			return domain.ErrUserNotFound(err)
		}

		return err
	}

	code := service.randomizer.GenerateRandomNumericCode(6)
	err = service.repository.UpdateVerificationCode(ctx, user.ID, code, "email")
	if err != nil {
		return err
	}

	go service.EventBus.Publish(ctx, domain.UserUpdateCodeEvent{
		User: user,
		Code: code,
	})

	return nil
}

// Initiates the process to change the user's email by sending verification tokens.
func (service *RegisterService) RequestEmailVerificationChange(ctx context.Context, input RegisterEmailChangeInput) error {
	user, err := service.repository.FindByIdentifier(ctx, input.Identifier, input.Type)
	if err != nil || user == nil {
		return domain.ErrUserNotFound(err)
	}

	currentToken := service.randomizer.GenerateRandomString(32, true)
	newToken := service.randomizer.GenerateRandomString(32, true)

	if err := service.repository.SaveEmailChangeRequest(ctx, user.ID, input.NewEmail, currentToken, newToken); err != nil {
		return err
	}

	go service.EventBus.Publish(ctx, domain.UserEmailChangeRequestedEvent{
		User:         user,
		CurrentToken: currentToken,
		NewToken:     newToken,
		NewEmail:     input.NewEmail,
	})

	return nil
}

// Confirms the email change by validating the provided tokens and updating the user's email.
func (service *RegisterService) ConfirmEmailVerificationChange(ctx context.Context, input RegisterEmailChangeConfirmInput) error {
	user, err := service.repository.FindUserByEmailChangeTokens(ctx, input.CurrentToken, input.NewToken)
	if err != nil || user == nil {
		return domain.ErrUserNotFound(err)
	}

	if err := service.repository.ApplyEmailChange(ctx, user.ID); err != nil {
		return err
	}

	updatedUser, _ := service.repository.FindByIdentifier(ctx, user.ID.String(), string(domain.IdentifierUserId))
	if updatedUser != nil && updatedUser.ClientId != nil && *updatedUser.ClientId > 0 {
		if err := service.authorization.With(updatedUser.ID, *updatedUser.ClientId).UpdateIdentifier(updatedUser.Email); err != nil {
			return err
		}
	}

	code := service.randomizer.GenerateRandomNumericCode(6)
	err = service.repository.UpdateVerificationCode(ctx, user.ID, code, "email")
	if err != nil {
		return err
	}

	go service.EventBus.Publish(ctx, domain.UserUpdateCodeEvent{
		User: updatedUser,
		Code: code,
	})

	return nil
}
