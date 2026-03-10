/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Auth/Infraestructure/Persistence
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
	"github.com/vayload/vayload/pkg/collect"
)

type IdentifierKey string

func (i IdentifierKey) Key() string {
	switch domain.IdentifierType(i) {
	case domain.IdentifierTypeEmail:
		return "email"
	case domain.IdentifierTypePhone:
		return "phone"
	case domain.IdentifierUserId:
		return "id"
	default:
		return ""
	}
}

func (i IdentifierKey) IsSupported() bool {
	switch domain.IdentifierType(i) {
	case domain.IdentifierTypeEmail, domain.IdentifierTypePhone, domain.IdentifierUserId, domain.IdentifierTypeUsername:
		return true
	default:
		return false
	}
}

type UserIdentityModel struct {
	UserID       snowflake.ID   `db:"id"`
	Email        string         `db:"email"`
	Phone        sql.NullString `db:"phone"`
	CountryID    sql.NullInt64  `db:"country_id"`
	DocumentID   sql.NullString `db:"document_id"`
	DocumentType sql.NullString `db:"document_type"`
}

func (userIdentityModel *UserIdentityModel) MapToDomain() *domain.UserIdentity {
	return &domain.UserIdentity{
		UserID:       userIdentityModel.UserID,
		Email:        userIdentityModel.Email,
		Phone:        nilIfInvalidString(userIdentityModel.Phone),
		CountryID:    nilIfInvalidInt64(userIdentityModel.CountryID),
		DocumentID:   nilIfInvalidString(userIdentityModel.DocumentID),
		DocumentType: nilIfInvalidString(userIdentityModel.DocumentType),
	}
}

type UserModel struct {
	ID                       snowflake.ID   `db:"id"`
	Role                     string         `db:"role"`
	DocumentID               sql.NullString `db:"document_id"`
	DocumentType             sql.NullString `db:"document_type"`
	FirstName                string         `db:"first_name"`
	LastName                 sql.NullString `db:"last_name"`
	Username                 string         `db:"username"`
	Email                    string         `db:"email"`
	PasswordHash             sql.NullString `db:"password_hash"`
	EmailConfirmedAt         sql.NullTime   `db:"email_confirmed_at"`
	InvitedAt                sql.NullTime   `db:"invited_at"`
	ConfirmationToken        sql.NullString `db:"confirmation_token"`
	ConfirmationSentAt       sql.NullTime   `db:"confirmation_sent_at"`
	RecoveryToken            sql.NullString `db:"recovery_token"`
	RecoverySentAt           sql.NullTime   `db:"recovery_sent_at"`
	EmailChangeTokenNew      sql.NullString `db:"email_change_token_new"`
	EmailChange              sql.NullString `db:"email_change"`
	EmailChangeSentAt        sql.NullTime   `db:"email_change_sent_at"`
	Phone                    sql.NullString `db:"phone"`
	PhoneConfirmedAt         sql.NullTime   `db:"phone_confirmed_at"`
	PhoneChange              sql.NullString `db:"phone_change"`
	PhoneChangeToken         sql.NullString `db:"phone_change_token"`
	PhoneChangeSentAt        sql.NullTime   `db:"phone_change_sent_at"`
	ConfirmedAt              sql.NullTime   `db:"confirmed_at"`
	EmailChangeTokenCurrent  sql.NullString `db:"email_change_token_current"`
	EmailChangeConfirmStatus sql.NullInt16  `db:"email_change_confirm_status"`
	OTPCode                  sql.NullString `db:"otp_code"`
	OTPSentAt                sql.NullTime   `db:"otp_sent_at"`
	BannedUntil              sql.NullTime   `db:"banned_until"`
	ReauthenticationToken    sql.NullString `db:"reauthentication_token"`
	ReauthenticationSentAt   sql.NullTime   `db:"reauthentication_sent_at"`
	AvatarURL                sql.NullString `db:"avatar_url"`
	Metadata                 sql.NullString `db:"metadata"`
	Settings                 sql.NullString `db:"settings"`
	Attributes               sql.NullString `db:"attributes"`
	ClientID                 sql.NullInt64  `db:"client_id"`
	ProfileID                sql.NullInt64  `db:"profile_id"`
	DigevoCoreID             sql.NullInt64  `db:"digevo_core_id"`
	CountryID                sql.NullInt64  `db:"country_id"`
	AuthType                 sql.NullString `db:"auth_type"`
	AuthSID                  sql.NullString `db:"auth_sid"`
	LastSignInAt             sql.NullTime   `db:"last_sign_in_at"`
	IsSuperAdmin             sql.NullBool   `db:"is_super_admin"`
	IsSSOUser                sql.NullBool   `db:"is_sso_user"`
	CreatedAt                time.Time      `db:"created_at"`
	UpdatedAt                time.Time      `db:"updated_at"`
	DeletedAt                sql.NullTime   `db:"deleted_at"`
}

func (userModel *UserModel) MapToDomain() *domain.User {
	return &domain.User{
		ID:            userModel.ID,
		Username:      userModel.Username,
		Email:         userModel.Email,
		Role:          domain.UserRole(userModel.Role),
		Phone:         nilIfInvalidString(userModel.Phone),
		Password:      nilIfInvalidString(userModel.PasswordHash),
		CreatedAt:     userModel.CreatedAt,
		UpdatedAt:     userModel.UpdatedAt,
		CountryID:     nilIfInvalidFlakeId(userModel.CountryID),
		ClientId:      nilIfInvalidInt64(userModel.ClientID),
		ProfileId:     nilIfInvalidInt64(userModel.ProfileID),
		DigevoCoreId:  nilIfInvalidInt64(userModel.DigevoCoreID),
		AuthType:      userModel.AuthType.String,
		AvatarURL:     nilIfInvalidString(userModel.AvatarURL),
		EmailVerified: userModel.EmailConfirmedAt.Valid,
		PhoneVerified: userModel.PhoneConfirmedAt.Valid,
		BannedUntil:   nilIfInvalidTime(userModel.BannedUntil),
		DeletedAt:     nilIfInvalidTime(userModel.DeletedAt),
		OTPCode:       userModel.OTPCode.String,
	}
}

type ClientID sql.NullInt64
type ProfileID sql.NullInt64
type DigevoCoreID sql.NullInt64
type CountryID sql.NullInt64

func (m *UserModel) GetClientID() *int64         { return nilIfInvalidInt64(m.ClientID) }
func (m *UserModel) GetProfileID() *int64        { return nilIfInvalidInt64(m.ProfileID) }
func (m *UserModel) GetDigevoCoreID() *int64     { return nilIfInvalidInt64(m.DigevoCoreID) }
func (m *UserModel) GetCountryID() *snowflake.ID { return nilIfInvalidFlakeId(m.CountryID) }

func (userModel *UserModel) MapToAuth() *domain.UserAuth {
	return &domain.UserAuth{
		ID:                    userModel.ID,
		Role:                  domain.UserRole(userModel.Role),
		Email:                 userModel.Email,
		Phone:                 nilIfInvalidString(userModel.Phone),
		PasswordHash:          userModel.PasswordHash.String,
		EmailConfirmed:        userModel.EmailConfirmedAt.Valid,
		PhoneConfirmed:        userModel.PhoneConfirmedAt.Valid,
		ConfirmationToken:     userModel.ConfirmationToken.String,
		RecoveryToken:         userModel.RecoveryToken.String,
		OTPCode:               userModel.OTPCode.String,
		BannedUntil:           nilIfInvalidTime(userModel.BannedUntil),
		DeletedAt:             nilIfInvalidTime(userModel.DeletedAt),
		ReauthenticationToken: userModel.ReauthenticationToken.String,
		AuthType:              userModel.AuthType.String,
		AuthSID:               userModel.AuthSID.String,
		IsSuperAdmin:          userModel.IsSuperAdmin.Valid && userModel.IsSuperAdmin.Bool,
		IsSSOUser:             userModel.IsSSOUser.Valid && userModel.IsSSOUser.Bool,
		LastSignInAt:          nilIfInvalidTime(userModel.LastSignInAt),
	}
}

type UserVerificationModel struct {
	ID                       snowflake.ID   `db:"id"`
	Email                    string         `db:"email"`
	EmailConfirmedAt         sql.NullTime   `db:"email_confirmed_at"`
	VerificationCode         sql.NullString `db:"verification_code"`
	ConfirmationToken        sql.NullString `db:"confirmation_token"`
	ConfirmationSentAt       sql.NullTime   `db:"confirmation_sent_at"`
	RecoveryToken            sql.NullString `db:"recovery_token"`
	RecoverySentAt           sql.NullTime   `db:"recovery_sent_at"`
	EmailChangeTokenNew      sql.NullString `db:"email_change_token_new"`
	EmailChange              sql.NullString `db:"email_change"`
	EmailChangeSentAt        sql.NullTime   `db:"email_change_sent_at"`
	EmailChangeTokenCurrent  sql.NullString `db:"email_change_token_current"`
	EmailChangeConfirmStatus sql.NullInt32  `db:"email_change_confirm_status"`
	Phone                    sql.NullString `db:"phone"`
	PhoneConfirmedAt         sql.NullTime   `db:"phone_confirmed_at"`
	PhoneChange              sql.NullString `db:"phone_change"`
	PhoneChangeToken         sql.NullString `db:"phone_change_token"`
	PhoneChangeSentAt        sql.NullTime   `db:"phone_change_sent_at"`
	OTPCode                  sql.NullString `db:"otp_code"`
	OTPSentAt                sql.NullTime   `db:"otp_sent_at"`
	ReauthenticationToken    sql.NullString `db:"reauthentication_token"`
	ReauthenticationSentAt   sql.NullTime   `db:"reauthentication_sent_at"`
	ConfirmedAt              sql.NullTime   `db:"confirmed_at"`
	CreatedAt                sql.NullTime   `db:"created_at"`
	UpdatedAt                sql.NullTime   `db:"updated_at"`
	DeletedAt                sql.NullTime   `db:"deleted_at"`
}

func (m UserVerificationModel) MapToDomain() *domain.UserVerification {
	return &domain.UserVerification{
		ID:                       m.ID,
		Email:                    m.Email,
		EmailConfirmedAt:         nilIfInvalidTime(m.EmailConfirmedAt),
		VerificationCode:         nilIfInvalidString(m.VerificationCode),
		ConfirmationToken:        nilIfInvalidString(m.ConfirmationToken),
		ConfirmationSentAt:       nilIfInvalidTime(m.ConfirmationSentAt),
		RecoveryToken:            nilIfInvalidString(m.RecoveryToken),
		RecoverySentAt:           nilIfInvalidTime(m.RecoverySentAt),
		EmailChangeTokenNew:      nilIfInvalidString(m.EmailChangeTokenNew),
		EmailChange:              nilIfInvalidString(m.EmailChange),
		EmailChangeSentAt:        nilIfInvalidTime(m.EmailChangeSentAt),
		EmailChangeTokenCurrent:  nilIfInvalidString(m.EmailChangeTokenCurrent),
		EmailChangeConfirmStatus: nilIfInvalidInt32(m.EmailChangeConfirmStatus),
		Phone:                    nilIfInvalidString(m.Phone),
		PhoneConfirmedAt:         nilIfInvalidTime(m.PhoneConfirmedAt),
		PhoneChange:              nilIfInvalidString(m.PhoneChange),
		PhoneChangeToken:         nilIfInvalidString(m.PhoneChangeToken),
		PhoneChangeSentAt:        nilIfInvalidTime(m.PhoneChangeSentAt),
		OTPCode:                  nilIfInvalidString(m.OTPCode),
		OTPSentAt:                nilIfInvalidTime(m.OTPSentAt),
		ReauthenticationToken:    nilIfInvalidString(m.ReauthenticationToken),
		ReauthenticationSentAt:   nilIfInvalidTime(m.ReauthenticationSentAt),
		ConfirmedAt:              nilIfInvalidTime(m.ConfirmedAt),
		CreatedAt:                nilIfInvalidTime(m.CreatedAt).UTC(),
		UpdatedAt:                nilIfInvalidTime(m.UpdatedAt).UTC(),
		DeletedAt:                nilIfInvalidTime(m.DeletedAt),
	}
}

type AuthorizationBindingModel struct {
	UserId    snowflake.ID   `db:"id"`
	ClientId  sql.NullInt64  `db:"client_id"`
	ProfileId sql.NullInt64  `db:"profile_id"`
	Signature sql.NullString `db:"profile_signature"`
}

func (model *AuthorizationBindingModel) MapToDomain() *domain.AuthorizationBindingOptional {
	return &domain.AuthorizationBindingOptional{
		UserId:    model.UserId,
		ClientId:  nilIfInvalidInt64(model.ClientId),
		Signature: model.Signature.String,
		ProfileId: nilIfInvalidInt64(model.ProfileId),
	}
}

type TelcoLocationModel struct {
	CountryID sql.NullInt64  `db:"country_id"`
	IDList    sql.NullString `db:"idlists"`
}

func (tl *TelcoLocationModel) MapToDomain() *domain.TelcoLocation {
	idlist := 0
	if tl.IDList.Valid {
		idlist, _ = strconv.Atoi(tl.IDList.String)
	}

	tempId := snowflake.ID(tl.CountryID.Int64)
	return &domain.TelcoLocation{
		CountryID: &tempId,
		IDList:    int64(idlist),
	}
}

func nilIfInvalidString(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nilIfInvalidInt64(ns sql.NullInt64) *int64 {
	if ns.Valid {
		return &ns.Int64
	}
	return nil
}

func nilIfInvalidInt32(ns sql.NullInt32) *int32 {
	if ns.Valid {
		return &ns.Int32
	}
	return nil
}

func nilIfInvalidTime(ns sql.NullTime) *time.Time {
	if ns.Valid {
		return &ns.Time
	}
	return nil
}

func nilIfInvalidFlakeId(ns sql.NullInt64) *snowflake.ID {
	if ns.Valid {
		id := snowflake.ID(ns.Int64)
		return &id
	}
	return nil
}

type UserRepository struct {
	db connection.DatabaseConnection
}

func NewUserRepository(db connection.DatabaseConnection) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) FindUserIdentity(ctx context.Context, identifier string, identifierType string) (*domain.UserIdentity, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	model := UserIdentityModel{}
	fields := []string{"id", "email", "phone", "country_id", "document_id", "document_type"}

	if err := repo.db.From("users").Select(fields...).Where("email", "=", identifier).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

func (repo *UserRepository) FindByIdentifier(ctx context.Context, identifier string, identifierType string) (*domain.User, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	model := UserModel{}
	if err := repo.db.From("users").Where(identifierKey.Key(), "=", identifier).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

func (repo *UserRepository) FindUserAuthByIdentifier(ctx context.Context, identifier string, identifierType string) (*domain.UserAuth, error) {
	identifierKey := IdentifierKey(identifierType)
	if !identifierKey.IsSupported() {
		return nil, fmt.Errorf("unsupported identifier type: %s", identifierType)
	}

	model := UserModel{}
	if err := repo.db.From("users").Where(identifierKey.Key(), "=", identifier).Get(&model); err != nil {
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

	err := repo.db.From("users").Where("id", "=", metadata.User.ID).Update(data).Exec()

	return err
}

func (repo *UserRepository) UpdatePassword(ctx context.Context, userId snowflake.ID, newPassword string) error {
	data := map[string]any{
		"password_hash": newPassword,
	}

	err := repo.db.From("users").Where("id", "=", userId).Update(data).Exec()

	return err
}

func (repo *UserRepository) CreateUserWithSettings(ctx context.Context, user *domain.User, settings *domain.UserSettings) error {
	fields := map[string]any{
		"first_name":    user.Username,
		"last_name":     user.LastName,
		"username":      user.Username,
		"email":         user.Email,
		"role":          user.Role,
		"avatar_url":    user.AvatarURL,
		"auth_type":     user.AuthType,
		"password_hash": user.Password,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}
	if user.Phone != nil {
		fields["phone"] = user.Phone
	}
	if user.ClientId != nil {
		fields["client_id"] = *user.ClientId
	}
	if user.ProfileId != nil {
		fields["profile_id"] = *user.ProfileId
	}
	if user.CountryID != nil {
		fields["country_id"] = *user.CountryID
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

	userId := snowflake.Node.Generate()
	fields["id"] = userId

	err := repo.db.From("users").Insert(fields).Exec()
	if err != nil {
		return err
	}
	user.ID = userId

	err = repo.db.From("notification_preferences").Insert(map[string]any{
		"user_id":  userId,
		"sms":      notifyFields["sms"],
		"email":    notifyFields["email"],
		"whatsapp": notifyFields["whatsapp"],
		"purchase": notifyFields["purchase"],
	}).Exec()

	return err
}

func (repo *UserRepository) CreateUserWithCode(ctx context.Context, user *domain.User, token string, code string) (*domain.User, error) {
	userId := snowflake.Node.Generate()

	err := repo.db.From("users").Insert(map[string]any{
		"id":                   userId,
		"password_hash":        user.Password,
		"first_name":           user.Username,
		"last_name":            "",
		"username":             user.Username,
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
}

func (repo *UserRepository) FindCodesByIdentifier(ctx context.Context, identifier string, identifierType string, typeCode string) (*domain.UserVerification, error) {
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
	err := repo.db.From("users").Select(fields...).Where(identifierKey.Key(), "=", identifier).Get(&model)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotResults
		}
		return nil, err
	}

	return model.MapToDomain(), nil
}

func (repo *UserRepository) BindAuthorization(ctx context.Context, userId snowflake.ID, binding *domain.AuthorizationBinding, meta *domain.UserMeta) error {
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

	err := repo.db.From("users").Where("id", "=", userId).Update(data).Exec()
	if err != nil {
		return err
	}

	err = repo.db.From("notification_preferences").Insert(map[string]any{
		"user_id": userId,
		"email":   true,
	}).Exec()

	return err
}

func (repo *UserRepository) UpdateVerificationCode(ctx context.Context, userId snowflake.ID, code string, typeCode string) error {
	return nil
}

func (repo *UserRepository) SaveRecoveryToken(ctx context.Context, userId snowflake.ID, token string) error {
	query := `UPDATE users SET recovery_token = ?, recovery_sent_at = ? WHERE id = ?`
	return repo.db.Prepared(ctx, query, token, time.Now().UTC(), int64(userId))
}

func (repo *UserRepository) FindUserByRecoveryToken(ctx context.Context, token string) (*domain.User, error) {
	model := UserModel{}
	if err := repo.db.From("users").Select("id", "email", "password_hash", "client_id", "profile_id", "country_id", "digevo_core_id", "created_at", "updated_at").Where("recovery_token", "=", token).Get(&model); err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

func (repo *UserRepository) ResetPasswordWithRecoveryToken(ctx context.Context, token string, newPassword string) error {
	query := `UPDATE users SET password_hash = ?, recovery_token = NULL, recovery_sent_at = NULL WHERE recovery_token = ?`
	return repo.db.Prepared(ctx, query, newPassword, token)
}

func (repo *UserRepository) SaveEmailChangeRequest(ctx context.Context, userId snowflake.ID, newEmail string, currentToken string, newToken string) error {
	query := `UPDATE users SET email_change = ?, email_change_token_current = ?, email_change_token_new = ?, email_change_sent_at = ? WHERE id = ?`
	return repo.db.Prepared(ctx, query, newEmail, currentToken, newToken, time.Now().UTC(), int64(userId))
}

func (repo *UserRepository) FindUserByEmailChangeTokens(ctx context.Context, currentToken string, newToken string) (*domain.User, error) {
	var model UserModel
	fields := []string{"id", "email", "email_change", "client_id", "profile_id", "country_id", "digevo_core_id", "created_at", "updated_at"}
	if err := repo.db.From("users").Select(fields...).Where("email_change_token_current", "=", currentToken).Where("email_change_token_new", "=", newToken).Get(&model); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotResults
		}
		return nil, err
	}
	return model.MapToDomain(), nil
}

func (repo *UserRepository) ApplyEmailChange(ctx context.Context, userId snowflake.ID) error {
	query := `UPDATE users SET email = email_change, email_change = NULL, email_change_token_current = NULL, email_change_token_new = NULL, email_change_confirm_status = 1 WHERE id = ?`
	return repo.db.Prepared(ctx, query, int64(userId))
}

func (repo *UserRepository) FindCountryOtpProviders(ctx context.Context, countryId snowflake.ID) (*domain.OtpProvider, error) {
	query := "SELECT JSON_OBJECTAGG(REPLACE(`meta_key`, '_senders', ''), `meta_value`) `providers` FROM countries_meta WHERE `meta_key` IN ('sms_senders', 'email_senders', 'whatsapp_senders') AND country_id = ?"

	var rawProviders struct {
		Providers sql.NullString `db:"providers"`
	}

	err := repo.db.SelectOne(ctx, &rawProviders, query, countryId)
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

func (repo *UserRepository) SaveOtpCode(ctx context.Context, userId snowflake.ID, otpCode string) error {
	data := map[string]any{
		"otp_code": otpCode,
	}

	err := repo.db.From("users").Where("id", "=", userId).Update(data).Exec()

	return err
}

func (repo *UserRepository) GetTelcoProfiles(ctx context.Context, telcoId string, countryId snowflake.ID) ([]int64, error) {
	var model struct {
		MetaValue string `db:"meta_value"`
	}

	if err := repo.db.From("countries_meta").Select("meta_value").Where("meta_key", "=", telcoId).Where("country_id", "=", int64(countryId)).Get(&model); err != nil {
		if err == sql.ErrNoRows {
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

func (repo *UserRepository) GetTelcoLocation(ctx context.Context, productName string, countryCode string) (*domain.TelcoLocation, error) {
	query := `SELECT
        m.country_id,
        m.meta_value idlists
    FROM countries_meta m
        JOIN countries c ON c.id = m.country_id
    WHERE FIND_IN_SET(?, m.meta_value)
        AND m.meta_key = 'vas_products'
        AND c.code = ?`

	model := TelcoLocationModel{}
	err := repo.db.SelectOne(ctx, &model, query, productName, countryCode)
	if err != nil {
		return nil, selectQueryError(err)
	}

	return model.MapToDomain(), nil
}

type AuthLogRepository struct {
	db connection.DatabaseConnection
}

func NewAuthLogRepository(db connection.DatabaseConnection) *AuthLogRepository {
	return &AuthLogRepository{db: db}
}

func (repo *AuthLogRepository) SaveLogin(ctx context.Context, userId snowflake.ID, ipAddress string, userAgent string, method string, phone *string, email string) error {
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

	err := repo.db.From("audit_log_entries").Insert(data).Exec()
	return err
}

func (repo *AuthLogRepository) SaveExternalRequestLog(ctx context.Context, log *domain.ExternalRequestLogging) error {
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

	err := repo.db.From("external_services_logs").Insert(data).Exec()
	return err
}

func selectQueryError(err error) error {
	if err == nil {
		return nil
	}
	if err == sql.ErrNoRows {
		return domain.ErrEmptyResultSet
	}
	if err == domain.ErrEmptyResultSet {
		return domain.ErrEmptyResultSet
	}

	return err
}

var _ domain.AuthLogRepository = (*AuthLogRepository)(nil)
var _ domain.AuthRepository = (*UserRepository)(nil)
