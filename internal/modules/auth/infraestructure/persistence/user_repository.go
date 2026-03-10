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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
	"github.com/vayload/vayload/pkg/collect"
)

const (
	USER_SCHEMA_NAME    = "users"
	PATIENT_SCHEMA_NAME = "patients"
)

type UserRepository struct {
	database connection.DatabaseConnection
	users    connection.QueryBuilder
}

func NewUserRepository(database connection.DatabaseConnection) *UserRepository {
	return &UserRepository{
		database: database,
		users:    database.From(USER_SCHEMA_NAME),
	}
}

func (mysql *UserRepository) FindUserIdentity(ctx context.Context, identifier string, identifierType string) (*domain.UserIdentity, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	model := UserIdentityModel{}
	fields := []string{"id", "email", "phone", "country_id", "document_id", "document_type"}

	if err := mysql.database.From(USER_SCHEMA_NAME).Select(fields...).Where(identifierKey.Key(), "=", identifier).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

func (mysql *UserRepository) FindByIdentifier(ctx context.Context, identifier string, identifierType string) (*domain.User, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	model := UserModel{}
	if err := mysql.database.From(USER_SCHEMA_NAME).Where(identifierKey.Key(), "=", identifier).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

func (mysql *UserRepository) FindUserAuthByIdentifier(ctx context.Context, identifier string, identifierType string) (*domain.UserAuth, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	model := UserModel{}
	if err := mysql.users.Where(identifierKey.Key(), "=", identifier).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToAuth(), nil
}

func (repo *UserRepository) UpdateMetadata(ctx context.Context, metadata *domain.UserMetadata) error {
	data := map[string]any{
		"otp_code":           metadata.OtpCode,
		"verification_code":  metadata.VerificationCode,
		"confirmation_token": metadata.ConfirmationToken,
	}

	err := repo.users.Where("id", "=", metadata.User.ID).Update(data).Exec()

	return err
}

func (mysql *UserRepository) UpdatePassword(ctx context.Context, userId snowflake.ID, newPassword string) error {
	data := map[string]any{
		"password_hash": newPassword,
	}

	err := mysql.users.Where("id", "=", userId).Update(data).Exec()

	return err
}

func (mysql *UserRepository) CreateUserWithSettings(ctx context.Context, user *domain.User, settings *domain.UserSettings) error {
	fields := map[string]any{
		"first_name":    user.Username,
		"last_name":     user.LastName,
		"email":         user.Email,
		"role":          user.Role,
		"avatar_url":    user.AvatarURL,
		"auth_type":     user.AuthType, // e.g. "email", "phone", "oauth", "sso"
		"password_hash": user.Password,
		"country_id":    user.CountryID,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}
	if user.Phone != nil {
		fields["phone"] = user.Phone
	}
	// If user has a client ID, set it
	if user.ClientId != nil {
		fields["client_id"] = user.ClientId
	}
	// If user has a profile ID, set it
	if user.ProfileId != nil {
		fields["profile_id"] = user.ProfileId
	}

	curdate := time.Now().UTC()
	notifyFields := map[string]any{
		"sms":      false,
		"email":    false,
		"whatsapp": false,
		"purchase": false,
	}

	if settings.Notifications.SMS && user.Phone != nil {
		notifyFields["sms"] = true
		fields["phone_confirmed_at"] = curdate
	}
	if settings.Notifications.Email && user.Email != "" {
		notifyFields["email"] = true
		fields["email_confirmed_at"] = curdate
	}
	if settings.Notifications.WhatsApp && user.Phone != nil {
		notifyFields["whatsapp"] = true
		fields["phone_confirmed_at"] = curdate
	}

	// return mysql.database.Transaction(ctx, func(db storage.DatabaseContext) error {
	userId := snowflake.Node.Generate()
	fields["id"] = userId

	err := mysql.database.From(USER_SCHEMA_NAME).Insert(fields).Exec()
	if err != nil {
		return err
	}
	user.ID = userId

	err = mysql.database.From("notification_preferences").Upsert(map[string]any{
		"user_id":  userId,
		"sms":      notifyFields["sms"],
		"email":    notifyFields["email"],
		"whatsapp": notifyFields["whatsapp"],
		"purchase": notifyFields["purchase"],
	}, []string{"sms", "email", "whatsapp", "purchase"}).Exec()
	if err != nil {
		return err
	}

	if user.Role == domain.PatientRole {
		err = mysql.database.From(PATIENT_SCHEMA_NAME).Insert(map[string]any{
			"user_id":    userId,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		}).Exec()

		return err
	}

	return nil
	// })
}

func (mysql *UserRepository) CreateUserWithCode(ctx context.Context, user *domain.User, token string, code string) (*domain.User, error) {
	// err := mysql.database.Transaction(ctx, func(tx storage.DatabaseContext) error {
	userId := snowflake.Node.Generate()

	err := mysql.database.From("users").Insert(map[string]any{
		"id":                   userId,
		"password_hash":        user.Password,
		"first_name":           user.Username,
		"last_name":            "",
		"email":                user.Email,
		"verification_code":    code,
		"confirmation_token":   token,
		"confirmation_sent_at": time.Now().UTC(),
		"created_at":           user.CreatedAt,
		"updated_at":           user.UpdatedAt,
	}).Exec()

	if err != nil {
		return nil, err
	}

	user.ID = snowflake.ID(userId)

	return user, nil
	// })
}

func (mysql *UserRepository) FindCodesByIdentifier(ctx context.Context, identifier string, identifierType string, typeCode string) (*domain.UserVerification, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	fields := []string{
		"id",
		"email",
		"email_confirmed_at",
		"verification_code",
		"confirmation_token",
		"confirmation_sent_at",
		"recovery_token",
		"recovery_sent_at",
		"email_change_token_new",
		"email_change",
		"email_change_sent_at",
		"email_change_token_current",
		"email_change_confirm_status",
		"phone",
		"phone_confirmed_at",
		"phone_change",
		"phone_change_token",
		"phone_change_sent_at",
		"otp_code",
		"otp_sent_at",
		"reauthentication_token",
		"reauthentication_sent_at",
		"confirmed_at",
		"created_at",
		"updated_at",
		"deleted_at",
	}

	model := UserVerificationModel{}
	err := mysql.database.From(USER_SCHEMA_NAME).Where(identifierKey.Key(), "=", identifier).Select(fields...).Get(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotResults
		}
		return nil, err
	}

	return model.MapToDomain(), nil
}

func (mysql *UserRepository) BindAuthorization(ctx context.Context, userId snowflake.ID, binding *domain.AuthorizationBinding, meta *domain.UserMeta) error {
	data := map[string]any{
		"profile_id":        binding.ProfileId,
		"client_id":         binding.ClientId,
		"profile_signature": binding.Signature,
		"updated_at":        time.Now().UTC(),
	}

	if meta != nil {
		if meta.ConfirmationToken != "" {
			data["confirmation_token"] = meta.ConfirmationToken
		}
		if meta.VerificationCode != "" {
			data["verification_code"] = meta.VerificationCode
		}
		if meta.EmailVerified {
			data["email_confirmed_at"] = time.Now().UTC()
		}
		if meta.PhoneVerified {
			data["phone_confirmed_at"] = time.Now().UTC()
		}
		if meta.VasProvider != "" {
			data["auth_type"] = "vas"
		}
		if meta.DigevoCoreId > 0 {
			data["digevo_core_id"] = meta.DigevoCoreId
		}
	}

	// return mysql.database.Transaction(ctx, func(tx storage.DatabaseContext) error {
	err := mysql.database.From("users").Update(data).Where("id", "=", userId).Exec()
	if err != nil {
		return err
	}

	err = mysql.database.From("notification_preferences").Upsert(map[string]any{
		"user_id": userId,
		"email":   true,
	}, []string{"email"}).Exec()

	return err
	// })
}

func (mysql *UserRepository) UpdateVerificationCode(ctx context.Context, userId snowflake.ID, code string, typeCode string) error {
	// ctx := context.Background()

	// data := map[string]any{
	// 	"user_id":    userId,
	// 	"code":       code,
	// 	"type":       typeCode,
	// 	"expires_at": time.Now().UTC().Add(15 * time.Minute),
	// }

	// _, err := mysql.verification_codes.UpsertOne(ctx, data, []string{"code", "type", "expires_at"})

	return nil
}

func (mysql *UserRepository) SaveRecoveryToken(ctx context.Context, userId snowflake.ID, token string) error {
	query := `UPDATE users SET recovery_token = ?, recovery_sent_at = ? WHERE id = ?`
	return mysql.database.Prepared(ctx, query, token, time.Now().UTC(), int64(userId))
}

func (mysql *UserRepository) FindUserByRecoveryToken(ctx context.Context, token string) (*domain.User, error) {
	fields := []string{"id", "email", "password_hash", "client_id", "profile_id", "country_id", "digevo_core_id", "created_at", "updated_at"}
	filters := map[string]any{"recovery_token": token}
	model := UserModel{}

	if err := mysql.database.From("users").Wheres(filters).Select(fields...).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

func (mysql *UserRepository) ResetPasswordWithRecoveryToken(ctx context.Context, token string, newPassword string) error {
	query := `UPDATE users SET password_hash = ?, recovery_token = NULL, recovery_sent_at = NULL WHERE recovery_token = ?`
	err := mysql.database.Prepared(ctx, query, newPassword, token)

	return err
}

func (mysql *UserRepository) SaveEmailChangeRequest(ctx context.Context, userId snowflake.ID, newEmail string, currentToken string, newToken string) error {
	query := `UPDATE users SET email_change = ?, email_change_token_current = ?, email_change_token_new = ?, email_change_sent_at = ? WHERE id = ?`
	return mysql.database.Prepared(ctx, query, newEmail, currentToken, newToken, time.Now().UTC(), int64(userId))
}

func (mysql *UserRepository) FindUserByEmailChangeTokens(ctx context.Context, currentToken string, newToken string) (*domain.User, error) {
	var model UserModel
	fields := []string{"id", "email", "email_change", "client_id", "profile_id", "country_id", "digevo_core_id", "created_at", "updated_at"}
	filters := map[string]any{"email_change_token_current": currentToken, "email_change_token_new": newToken}

	if err := mysql.database.From("users").Wheres(filters).Select(fields...).Get(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotResults
		}
		return nil, err
	}

	return model.MapToDomain(), nil
}

func (mysql *UserRepository) ApplyEmailChange(ctx context.Context, userId snowflake.ID) error {
	query := `UPDATE users SET email = email_change, email_change = NULL, email_change_token_current = NULL, email_change_token_new = NULL, email_change_confirm_status = 1 WHERE id = ?`
	return mysql.database.Prepared(ctx, query, int64(userId))
}

func (mysql *UserRepository) FindCountryOtpProviders(ctx context.Context, countryId snowflake.ID) (*domain.OtpProvider, error) {
	query := "SELECT JSON_OBJECTAGG(REPLACE(`meta_key`, '_senders', ''), `meta_value`) `providers` FROM countries_meta WHERE `meta_key` IN ('sms_senders', 'email_senders', 'whatsapp_senders') AND country_id = ?"
	err := mysql.database.SelectOne(ctx, &rawProviders, query, countryId)
	if err != nil {
		return nil, err
	}

	if !rawProviders.Providers.Valid {
		return &domain.OtpProvider{
			Email: []string{"sendgrid"},
		}, nil
	}

	var rawProvidersList struct {
		Email    string `json:"email,omitempty"`
		SMS      string `json:"sms,omitempty"`
		WhatsApp string `json:"whatsapp,omitempty"`
	}

	if err := json.Unmarshal([]byte(rawProviders.Providers.String), &rawProvidersList); err != nil {
		return nil, err
	}

	providers := domain.OtpProvider{
		Email: collect.Filter(strings.Split(rawProvidersList.Email, ","), func(email string) bool {
			return email != ""
		}),
		SMS: collect.Filter(strings.Split(rawProvidersList.SMS, ","), func(sms string) bool {
			return sms != ""
		}),
		WhatsApp: collect.Filter(strings.Split(rawProvidersList.WhatsApp, ","), func(whatsapp string) bool {
			return whatsapp != ""
		}),
	}

	return &providers, nil
}

func (mysql *UserRepository) SaveOtpCode(ctx context.Context, userId snowflake.ID, otpCode string) error {
	data := map[string]any{
		"otp_code": otpCode,
	}

	err := mysql.users.Where("id", "=", userId).Update(data).Exec()

	return err
}

// Get telco profiles (registered in database as relationship with autorization service)
func (mysql *UserRepository) GetTelcoProfiles(ctx context.Context, telcoId string, countryId snowflake.ID) ([]int64, error) {
	var model struct {
		MetaValue string `db:"meta_value"`
	}

	filters := map[string]any{"meta_key": telcoId, "country_id": countryId}
	if err := mysql.database.From("countries_meta").Wheres(filters).Select("meta_value").Get(&model); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []int64{}, domain.ErrEmptyResultSet
		}

		return []int64{}, err
	}

	items := strings.Split(model.MetaValue, ",")
	profiles := make([]int64, len(items))
	for i, item := range items {
		id, err := strconv.ParseInt(item, 10, 64)
		if err == nil {
			profiles[i] = id
		}
	}

	return profiles, nil
}

const (
	USER_LOGIN_SCHEMA           = "audit_log_entries"
	EXTERNAL_REQUEST_LOG_SCHEMA = "external_services_logs"
)

type MyLogRepository struct {
	database connection.DatabaseConnection
}

func NewLogRepository(database connection.DatabaseConnection) *MyLogRepository {
	return &MyLogRepository{
		database: database,
	}
}

func (mysql *MyLogRepository) SaveLogin(ctx context.Context, userId snowflake.ID, ipAddress string, userAgent string, method string, phone *string, email string) error {
	id := snowflake.Node.Generate()

	payload := map[string]any{
		"email":  email,
		"msisdn": phone,
		"method": method,
	}
	payloadJson, _ := json.Marshal(payload)

	data := map[string]any{
		"id":         id,
		"actor_id":   userId,
		"payload":    string(payloadJson),
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"action":     "sign-in",
		"created_at": time.Now().UTC(),
	}

	err := mysql.database.From(USER_LOGIN_SCHEMA).Insert(data).Exec()
	return err
}

func (mysql *MyLogRepository) SaveExternalRequestLog(ctx context.Context, log *domain.ExternalRequestLogging) error {
	id := snowflake.Node.Generate()

	data := map[string]any{
		"id":              id,
		"request_id":      log.RequestID,
		"url":             log.URL,
		"http_method":     log.Method,
		"request_body":    log.RequestBody,
		"request_at":      log.RequestAt,
		"response_body":   log.ResponseBody,
		"response_status": log.ResponseStatus,
		"request_elapsed": log.RequestElapsed,
		"created_at":      time.Now().UTC(),
	}

	err := mysql.database.From(EXTERNAL_REQUEST_LOG_SCHEMA).Insert(data).Exec()
	return err
}

func selectQueryError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrEmptyResultSet
	}
	if errors.Is(err, domain.ErrEmptyResultSet) {
		return domain.ErrEmptyResultSet
	}

	return err
}

var _ domain.AuthLogRepository = (*MyLogRepository)(nil)
var _ domain.AuthRepository = (*UserRepository)(nil)
