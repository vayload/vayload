/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type UserRepository struct {
	database connection.DatabaseConnection
}

func NewUserRepository(database connection.DatabaseConnection) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (repo *UserRepository) users() connection.QueryBuilder {
	return repo.database.From("users")
}

func (repo *UserRepository) FindByID(ctx context.Context, id snowflake.ID) (*domain.User, error) {
	var model UserModel
	if err := repo.users().Where("id", "=", id).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	if err := repo.users().Where("email", "=", email).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model UserModel
	if err := repo.users().Where("username", "=", username).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) FindByIdentifier(ctx context.Context, identifier string, identifierType domain.IdentifierType) (*domain.User, error) {
	key := IdentifierKey(identifierType).Key()
	if key == "" {
		return nil, fmt.Errorf("unsupported identifier type: %v", identifierType)
	}

	var model UserModel
	if err := repo.users().Where(key, "=", identifier).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	data := map[string]any{
		"id":             user.ID,
		"username":       user.Username,
		"email":          user.Email,
		"password_hash":  user.PasswordHash,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"avatar_url":     user.AvatarURL,
		"is_super_admin": user.IsSuperAdmin,
		"is_sso_user":    user.IsSSOUser,
		"created_at":     user.CreatedAt,
	}
	if user.Phone != nil {
		data["phone"] = *user.Phone
	}
	return repo.users().InsertOne(data).Exec(ctx)
}

func (repo *UserRepository) Update(ctx context.Context, user *domain.User) error {
	data := map[string]any{
		"username":       user.Username,
		"email":          user.Email,
		"first_name":     user.FirstName,
		"last_name":      user.LastName,
		"avatar_url":     user.AvatarURL,
		"is_super_admin": user.IsSuperAdmin,
		"is_sso_user":    user.IsSSOUser,
		"updated_at":     time.Now().UTC(),

		"email_confirmed_at": user.EmailConfirmedAt,
		"phone_confirmed_at": user.PhoneConfirmedAt,
		"confirmed_at":       user.ConfirmedAt,
		"confirmation_token": user.ConfirmationToken,
		"recovery_token":     user.RecoveryToken,
		"otp_code":           user.OTPCode,
		"banned_until":       user.BannedUntil,
	}
	if user.Phone != nil {
		data["phone"] = *user.Phone
	}
	return repo.users().Where("id", "=", user.ID).UpdateOne(data).Exec(ctx)
}

func (repo *UserRepository) Delete(ctx context.Context, id snowflake.ID) error {
	// Soft delete
	return repo.users().Where("id", "=", id).UpdateOne(map[string]any{
		"deleted_at": time.Now().UTC(),
	}).Exec(ctx)
}

func (repo *UserRepository) UpdatePassword(ctx context.Context, userID snowflake.ID, passwordHash string) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"password_hash": passwordHash,
		"updated_at":    time.Now().UTC(),
	}).Exec(ctx)
}

func (repo *UserRepository) UpdateLastSignIn(ctx context.Context, userID snowflake.ID, ip, userAgent string) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"last_sign_in_at": time.Now().UTC(),
		// We could log IP/UA in a separate table or metadata if needed
	}).Exec(ctx)
}

func (repo *UserRepository) SaveConfirmationToken(ctx context.Context, userID snowflake.ID, token string) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"confirmation_token":   token,
		"confirmation_sent_at": time.Now().UTC(),
	}).Exec(ctx)
}

func (repo *UserRepository) SaveRecoveryToken(ctx context.Context, userID snowflake.ID, token string) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"recovery_token":   token,
		"recovery_sent_at": time.Now().UTC(),
	}).Exec(ctx)
}
func (repo *UserRepository) SaveOtpCode(ctx context.Context, userID snowflake.ID, code string) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"otp_code":    code,
		"otp_sent_at": time.Now().UTC(),
	}).Exec(ctx)
}

func (repo *UserRepository) FindUserByConfirmationToken(ctx context.Context, token string) (*domain.User, error) {
	var model UserModel
	if err := repo.users().Where("confirmation_token", "=", token).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) FindUserByRecoveryToken(ctx context.Context, token string) (*domain.User, error) {
	var model UserModel
	if err := repo.users().Where("recovery_token", "=", token).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) ConfirmEmail(ctx context.Context, userID snowflake.ID) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"email_confirmed_at": time.Now().UTC(),
		"confirmed_at":       time.Now().UTC(),
		"confirmation_token": nil,
	}).Exec(ctx)
}

func (repo *UserRepository) ResetPasswordWithRecoveryToken(ctx context.Context, token string, hashedPassword string) error {
	return repo.users().Where("recovery_token", "=", token).UpdateOne(map[string]any{
		"password_hash":    hashedPassword,
		"recovery_token":   nil,
		"recovery_sent_at": nil,
		"updated_at":       time.Now().UTC(),
	}).Exec(ctx)
}

// LogRepository implements domain.AuthLogRepository
type LogRepository struct {
	database connection.DatabaseConnection
}

func NewLogRepository(database connection.DatabaseConnection) *LogRepository {
	return &LogRepository{
		database: database,
	}
}

func (repo *LogRepository) SaveLogin(ctx context.Context, userID snowflake.ID, ipAddress string, userAgent string, method string, email string) error {
	data := map[string]any{
		"id":         snowflake.Node.Generate(),
		"actor_id":   userID,
		"action":     "login",
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"payload":    fmt.Sprintf("method: %s, email: %s", method, email),
		"created_at": time.Now().UTC(),
	}
	return repo.database.From("audit_log_entries").InsertOne(data).Exec(ctx)
}

func selectQueryError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrEmptyResultSet
	}
	return err
}

var _ domain.AuthRepository = (*UserRepository)(nil)

func (repo *UserRepository) SaveEmailChangeRequest(ctx context.Context, userID snowflake.ID, newEmail string, currentToken string, newToken string) error {
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"email_change":         newEmail,
		"email_change_token":   currentToken,
		"phone_change_token":   newToken,
		"email_change_sent_at": time.Now().UTC(),
	}).Exec(ctx)
}

func (repo *UserRepository) FindUserByEmailChangeTokens(ctx context.Context, currentToken string, newToken string) (*domain.User, error) {
	var model UserModel
	if err := repo.users().Where("email_change_token", "=", currentToken).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) ApplyEmailChange(ctx context.Context, userID snowflake.ID) error {
	var model UserModel
	if err := repo.users().Where("id", "=", userID).Get(ctx, &model); err != nil {
		return err
	}
	if !model.EmailChange.Valid {
		return fmt.Errorf("no pending email change")
	}
	return repo.users().Where("id", "=", userID).UpdateOne(map[string]any{
		"email":              model.EmailChange.String,
		"email_change":       nil,
		"email_change_token": nil,
		"phone_change_token": nil,
		"email_confirmed_at": time.Now().UTC(),
		"updated_at":         time.Now().UTC(),
	}).Exec(ctx)
}

var _ domain.AuthLogRepository = (*LogRepository)(nil)
