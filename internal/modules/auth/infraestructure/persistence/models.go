/*
 * SPDX-License-Identifier: MIT
 *
 * Vayload - Container
 *
 * Copyright (c) 2026 Alex Zweiter
 */

package persistence

import (
	"database/sql"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database"
	"github.com/vayload/vayload/internal/shared/snowflake"
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
	Email        string         `db:"email"`         // Email is optional, used for login
	Phone        sql.NullString `db:"phone"`         // Phone is optional, used for login
	CountryID    sql.NullInt64  `db:"country_id"`    // Optional country ID for OTP providers
	DocumentID   sql.NullString `db:"document_id"`   // Unique identifier for the user (DNI, passport, etc.)
	DocumentType sql.NullString `db:"document_type"` // Type of document (DNI, passport, etc.)
}

func (userIdentityModel *UserIdentityModel) MapToDomain() *domain.UserIdentity {
	return &domain.UserIdentity{
		UserID:       userIdentityModel.UserID,
		Email:        userIdentityModel.Email,
		Phone:        database.NilIfInvalidString(userIdentityModel.Phone),
		CountryID:    database.NilIfInvalidInt64(userIdentityModel.CountryID),
		DocumentID:   database.NilIfInvalidString(userIdentityModel.DocumentID),
		DocumentType: database.NilIfInvalidString(userIdentityModel.DocumentType),
	}
}

type VayloadUserRaw struct {
	ID                snowflake.ID   `db:"id"`
	PermsChecksum     sql.NullString `db:"perms_checksum"`
	UserCoreID        sql.NullInt64  `db:"user_core_id"`
	ProfileID         sql.NullInt64  `db:"profile_id"`
	FirstName         string         `db:"first_name"`
	LastName          string         `db:"last_name"`
	FullName          sql.NullString `db:"full_name"`
	RUT               sql.NullString `db:"rut"`
	Email             string         `db:"email"`
	Password          string         `db:"password_hash"`
	Phone             sql.NullString `db:"phone"`
	EmailVerified     sql.NullTime   `db:"email_verified"`
	PhoneVerified     sql.NullTime   `db:"phone_verified"`
	WhatsappNotify    sql.NullBool   `db:"whatsapp_notify"`
	SmsNotify         sql.NullBool   `db:"sms_notify"`
	EmailNotify       sql.NullBool   `db:"email_notify"`
	VerificationCode  sql.NullString `db:"verification_code"`
	OtpCode           sql.NullString `db:"otp_code"`
	ConfirmationToken sql.NullString `db:"confirmation_token"`
	TourVisited       bool           `db:"tour_visited"`
	UserFrom          sql.NullString `db:"user_from"`
	AuthSID           sql.NullString `db:"auth_sid"`
	UserVAS           sql.NullString `db:"user_vas"`
	DgvCoreID         sql.NullInt64  `db:"dgv_core_id"`
	CountryID         sql.NullInt64  `db:"country_id"`
	ReminderEnabled   sql.NullBool   `db:"reminder_enabled"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
	UpdatedBy         string         `db:"updated_by"`
}

type UserModel struct {
	ID                       snowflake.ID   `db:"id"`
	Role                     string         `db:"role"`
	DocumentID               sql.NullString `db:"document_id"`
	DocumentType             sql.NullString `db:"document_type"`
	FirstName                string         `db:"first_name"`
	LastName                 sql.NullString `db:"last_name"`
	Email                    string         `db:"email"`
	PasswordHash             sql.NullString `db:"password_hash"`
	EmailConfirmedAt         sql.NullTime   `db:"email_confirmed_at"`
	VerificationSentAt       sql.NullTime   `db:"verification_sent_at"`
	InvitedAt                sql.NullTime   `db:"invited_at"`
	VerificationCode         sql.NullString `db:"verification_code"`
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
	AuthType                 sql.NullString `db:"auth_type"`
	AuthSID                  sql.NullString `db:"auth_sid"`
	ClientID                 sql.NullInt64  `db:"client_id"`
	ProfileID                sql.NullInt64  `db:"profile_id"`
	ProfileSignature         sql.NullString `db:"profile_signature"`
	CountryID                sql.NullInt64  `db:"country_id"`
	Metadata                 sql.NullString `db:"metadata"`
	Integrations             sql.NullString `db:"integrations"`
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
		Username:      userModel.FirstName,
		Email:         userModel.Email,
		Role:          domain.UserRole(userModel.Role),
		Phone:         database.NilIfInvalidString(userModel.Phone),
		Password:      database.NilIfInvalidString(userModel.PasswordHash),
		CreatedAt:     userModel.CreatedAt,
		UpdatedAt:     userModel.UpdatedAt,
		CountryID:     database.NilIfInvalidFlakeId(userModel.CountryID),
		ClientId:      database.NilIfInvalidInt64(userModel.ClientID),
		ProfileId:     database.NilIfInvalidInt64(userModel.ProfileID),
		AuthType:      userModel.AuthType.String,
		AvatarURL:     database.NilIfInvalidString(userModel.AvatarURL),
		EmailVerified: userModel.EmailConfirmedAt.Valid,
		PhoneVerified: userModel.PhoneConfirmedAt.Valid,
		BannedUntil:   database.NilIfInvalidTime(userModel.BannedUntil),
		DeletedAt:     database.NilIfInvalidTime(userModel.DeletedAt),
		OTPCode:       userModel.OTPCode.String,
	}
}

func (userModel *UserModel) MapToAuth() *domain.UserAuth {
	return &domain.UserAuth{
		ID:                    userModel.ID,
		Role:                  domain.UserRole(userModel.Role),
		Email:                 userModel.Email,
		Phone:                 database.NilIfInvalidString(userModel.Phone),
		PasswordHash:          userModel.PasswordHash.String,
		EmailConfirmed:        userModel.EmailConfirmedAt.Valid,
		PhoneConfirmed:        userModel.PhoneConfirmedAt.Valid,
		ConfirmationToken:     userModel.ConfirmationToken.String,
		RecoveryToken:         userModel.RecoveryToken.String,
		OTPCode:               userModel.OTPCode.String,
		BannedUntil:           database.NilIfInvalidTime(userModel.BannedUntil),
		DeletedAt:             database.NilIfInvalidTime(userModel.DeletedAt),
		ReauthenticationToken: userModel.ReauthenticationToken.String,
		AuthType:              userModel.AuthType.String,
		AuthSID:               userModel.AuthSID.String,
		IsSuperAdmin:          userModel.IsSuperAdmin.Valid && userModel.IsSuperAdmin.Bool,
		IsSSOUser:             userModel.IsSSOUser.Valid && userModel.IsSSOUser.Bool,
		LastSignInAt:          database.NilIfInvalidTime(userModel.LastSignInAt),
	}
}

type UserVerificationCodeModel struct {
	VayloadUserRaw

	Code      string    `db:"code"`       // For verification code
	Type      string    `db:"type"`       // Type of verification (email, phone, etc.)
	ExpiresAt time.Time `db:"expires_at"` // Expiration time for the verification code
	Attempts  int       `db:"attempts"`   // Number of attempts for the verification code
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
		EmailConfirmedAt:         database.NilIfInvalidTime(m.EmailConfirmedAt),
		VerificationCode:         database.NilIfInvalidString(m.VerificationCode),
		ConfirmationToken:        database.NilIfInvalidString(m.ConfirmationToken),
		ConfirmationSentAt:       database.NilIfInvalidTime(m.ConfirmationSentAt),
		RecoveryToken:            database.NilIfInvalidString(m.RecoveryToken),
		RecoverySentAt:           database.NilIfInvalidTime(m.RecoverySentAt),
		EmailChangeTokenNew:      database.NilIfInvalidString(m.EmailChangeTokenNew),
		EmailChange:              database.NilIfInvalidString(m.EmailChange),
		EmailChangeSentAt:        database.NilIfInvalidTime(m.EmailChangeSentAt),
		EmailChangeTokenCurrent:  database.NilIfInvalidString(m.EmailChangeTokenCurrent),
		EmailChangeConfirmStatus: database.NilIfInvalidInt32(m.EmailChangeConfirmStatus),
		Phone:                    database.NilIfInvalidString(m.Phone),
		PhoneConfirmedAt:         database.NilIfInvalidTime(m.PhoneConfirmedAt),
		PhoneChange:              database.NilIfInvalidString(m.PhoneChange),
		PhoneChangeToken:         database.NilIfInvalidString(m.PhoneChangeToken),
		PhoneChangeSentAt:        database.NilIfInvalidTime(m.PhoneChangeSentAt),
		OTPCode:                  database.NilIfInvalidString(m.OTPCode),
		OTPSentAt:                database.NilIfInvalidTime(m.OTPSentAt),
		ReauthenticationToken:    database.NilIfInvalidString(m.ReauthenticationToken),
		ReauthenticationSentAt:   database.NilIfInvalidTime(m.ReauthenticationSentAt),
		ConfirmedAt:              database.NilIfInvalidTime(m.ConfirmedAt),
		CreatedAt:                database.NilIfInvalidTime(m.CreatedAt).UTC(),
		UpdatedAt:                database.NilIfInvalidTime(m.UpdatedAt).UTC(),
		DeletedAt:                database.NilIfInvalidTime(m.DeletedAt),
	}
}

func (raw *UserVerificationCodeModel) MapToDomain() *domain.UserWithVerificationCodes {
	return &domain.UserWithVerificationCodes{
		User: domain.User{
			ID:       raw.ID,
			Username: raw.FirstName,
			Email:    raw.Email,
			Phone:    database.NilIfInvalidString(raw.Phone),
		},
		Code:      raw.Code,
		Type:      raw.Type,
		ExpiresAt: raw.ExpiresAt,
	}
}

type VayloadOtpProviderRaw struct {
	Email    []string `db:"email,omitempty"`    // Email OTP provider
	SMS      []string `db:"sms,omitempty"`      // SMS OTP provider
	WhatsApp []string `db:"whatsapp,omitempty"` // WhatsApp OTP provider
}

var rawProviders struct {
	Providers sql.NullString `db:"providers"`
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
		ClientId:  database.NilIfInvalidInt64(model.ClientId),
		Signature: model.Signature.String,
		ProfileId: database.NilIfInvalidInt64(model.ProfileId),
	}
}
