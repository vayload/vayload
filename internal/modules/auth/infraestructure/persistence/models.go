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
	case domain.IdentifierTypeUsername:
		return "username"
	case domain.IdentifierTypeID:
		return "id"
	default:
		return ""
	}
}

func (i IdentifierKey) IsSupported() bool {
	return domain.IdentifierType(i).Valid()
}

// UserModel represents the 'users' table.
type UserModel struct {
	ID           snowflake.ID   `db:"id"`
	Username     string         `db:"username"`
	Email        sql.NullString `db:"email"`
	Phone        sql.NullString `db:"phone"`
	PasswordHash sql.NullString `db:"password_hash"`
	FirstName    sql.NullString `db:"first_name"`
	LastName     sql.NullString `db:"last_name"`
	AvatarURL    sql.NullString `db:"avatar_url"`

	EmailConfirmedAt   sql.NullTime   `db:"email_confirmed_at"`
	PhoneConfirmedAt   sql.NullTime   `db:"phone_confirmed_at"`
	ConfirmedAt        sql.NullTime   `db:"confirmed_at"`
	ConfirmationToken  sql.NullString `db:"confirmation_token"`
	RecoveryToken      sql.NullString `db:"recovery_token"`
	EmailChangeToken   sql.NullString `db:"email_change_token"`
	PhoneChangeToken   sql.NullString `db:"phone_change_token"`
	OTPCode            sql.NullString `db:"otp_code"`
	ConfirmationSentAt sql.NullTime   `db:"confirmation_sent_at"`
	RecoverySentAt     sql.NullTime   `db:"recovery_sent_at"`
	EmailChangeSentAt  sql.NullTime   `db:"email_change_sent_at"`
	PhoneChangeSentAt  sql.NullTime   `db:"phone_change_sent_at"`
	OTPSentAt          sql.NullTime   `db:"otp_sent_at"`
	EmailChange        sql.NullString `db:"email_change"`
	PhoneChange        sql.NullString `db:"phone_change"`

	BannedUntil  sql.NullTime `db:"banned_until"`
	LastSignInAt sql.NullTime `db:"last_sign_in_at"`

	Metadata   sql.NullString `db:"metadata"`
	Settings   sql.NullString `db:"settings"`
	Attributes sql.NullString `db:"attributes"`

	IsSuperAdmin bool `db:"is_super_admin"`
	IsSSOUser    bool `db:"is_sso_user"`

	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (m *UserModel) MapToDomain() *domain.User {
	user := &domain.User{
		ID:           m.ID,
		Username:     m.Username,
		Email:        m.Email.String,
		Phone:        database.NilIfInvalidString(m.Phone),
		PasswordHash: database.NilIfInvalidString(m.PasswordHash),
		FirstName:    database.NilIfInvalidString(m.FirstName),
		LastName:     database.NilIfInvalidString(m.LastName),
		AvatarURL:    database.NilIfInvalidString(m.AvatarURL),
		IsSuperAdmin: m.IsSuperAdmin,
		IsSSOUser:    m.IsSSOUser,

		EmailConfirmedAt: database.NilIfInvalidTime(m.EmailConfirmedAt),
		PhoneConfirmedAt: database.NilIfInvalidTime(m.PhoneConfirmedAt),
		ConfirmedAt:      database.NilIfInvalidTime(m.ConfirmedAt),

		ConfirmationToken: database.NilIfInvalidString(m.ConfirmationToken),
		RecoveryToken:     database.NilIfInvalidString(m.RecoveryToken),
		EmailChangeToken:  database.NilIfInvalidString(m.EmailChangeToken),
		PhoneChangeToken:  database.NilIfInvalidString(m.PhoneChangeToken),
		OTPCode:           database.NilIfInvalidString(m.OTPCode),
		EmailChange:       database.NilIfInvalidString(m.EmailChange),
		PhoneChange:       database.NilIfInvalidString(m.PhoneChange),
		BannedUntil:       database.NilIfInvalidTime(m.BannedUntil),

		ConfirmationSentAt: database.NilIfInvalidTime(m.ConfirmationSentAt),
		RecoverySentAt:     database.NilIfInvalidTime(m.RecoverySentAt),
		EmailChangeSentAt:  database.NilIfInvalidTime(m.EmailChangeSentAt),
		PhoneChangeSentAt:  database.NilIfInvalidTime(m.PhoneChangeSentAt),
		OTPSentAt:          database.NilIfInvalidTime(m.OTPSentAt),
		LastSignInAt:       database.NilIfInvalidTime(m.LastSignInAt),

		CreatedAt: m.CreatedAt,
		UpdatedAt: database.NilIfInvalidTime(m.UpdatedAt),
		DeletedAt: database.NilIfInvalidTime(m.DeletedAt),
	}

	return user
}

// ProjectModel represents the 'projects' table.
type ProjectModel struct {
	ID        snowflake.ID   `db:"id"`
	Name      string         `db:"name"`
	Slug      string         `db:"slug"`
	OwnerID   snowflake.ID   `db:"owner_id"`
	Settings  sql.NullString `db:"settings"`
	Locale    string         `db:"locale"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

func (m *ProjectModel) MapToDomain() *domain.Project {
	return &domain.Project{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		OwnerID:   m.OwnerID,
		Locale:    m.Locale,
		CreatedAt: m.CreatedAt,
		UpdatedAt: database.NilIfInvalidTime(m.UpdatedAt),
	}
}

// SessionModel represents the 'sessions' table.
type SessionModel struct {
	ID         string         `db:"id"`
	UserID     snowflake.ID   `db:"user_id"`
	ProjectID  sql.NullInt64  `db:"project_id"`
	IPAddress  sql.NullString `db:"ip_address"`
	UserAgent  sql.NullString `db:"user_agent"`
	LastSeenAt time.Time      `db:"last_seen_at"`
	ExpiresAt  time.Time      `db:"expires_at"`
	RevokedAt  sql.NullTime   `db:"revoked_at"`
	CreatedAt  time.Time      `db:"created_at"`
}

func (m *SessionModel) MapToDomain() *domain.Session {
	return &domain.Session{
		ID:         m.ID,
		UserID:     m.UserID,
		ProjectID:  database.NilIfInvalidFlakeId(m.ProjectID),
		IPAddress:  m.IPAddress.String,
		UserAgent:  m.UserAgent.String,
		LastSeenAt: m.LastSeenAt,
		ExpiresAt:  m.ExpiresAt,
		RevokedAt:  database.NilIfInvalidTime(m.RevokedAt),
		CreatedAt:  m.CreatedAt,
	}
}

// RefreshTokenModel represents the 'refresh_tokens' table.
type RefreshTokenModel struct {
	ID            string         `db:"id"`
	TokenHash     string         `db:"token_hash"`
	UserID        snowflake.ID   `db:"user_id"`
	FamilyID      string         `db:"family_id"`
	SessionID     sql.NullString `db:"session_id"`
	ParentID      sql.NullString `db:"parent_id"`
	UsedAt        sql.NullTime   `db:"used_at"`
	RevokedAt     sql.NullTime   `db:"revoked_at"`
	RevokedReason sql.NullString `db:"revoked_reason"`
	ExpiresAt     time.Time      `db:"expires_at"`
	CreatedAt     time.Time      `db:"created_at"`
}

func (m *RefreshTokenModel) MapToDomain() *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:            m.ID,
		TokenHash:     m.TokenHash,
		UserID:        m.UserID,
		FamilyID:      m.FamilyID,
		SessionID:     m.SessionID.String,
		ParentID:      m.ParentID.String,
		UsedAt:        database.NilIfInvalidTime(m.UsedAt),
		RevokedAt:     database.NilIfInvalidTime(m.RevokedAt),
		RevokedReason: m.RevokedReason.String,
		ExpiresAt:     m.ExpiresAt,
		CreatedAt:     m.CreatedAt,
	}
}

// RoleModel represents the 'roles' table.
type RoleModel struct {
	ID          snowflake.ID   `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	IsSystem    bool           `db:"is_system"`
	CreatedAt   time.Time      `db:"created_at"`
}

func (m *RoleModel) MapToDomain() *domain.Role {
	return &domain.Role{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description.String,
		IsSystem:    m.IsSystem,
		CreatedAt:   m.CreatedAt,
	}
}

// PermissionModel represents the 'permissions' table.
type PermissionModel struct {
	ID        snowflake.ID `db:"id"`
	Action    string       `db:"action"`
	Resource  string       `db:"resource"`
	CreatedAt time.Time    `db:"created_at"`
}

func (m *PermissionModel) MapToDomain() *domain.Permission {
	return &domain.Permission{
		ID:        m.ID,
		Action:    m.Action,
		Resource:  m.Resource,
		CreatedAt: m.CreatedAt,
	}
}

// ProjectMemberModel represents the 'project_members' table.
type ProjectMemberModel struct {
	ProjectID snowflake.ID `db:"project_id"`
	UserID    snowflake.ID `db:"user_id"`
	RoleID    snowflake.ID `db:"role_id"`
	CreatedAt time.Time    `db:"created_at"`
}

func (m *ProjectMemberModel) MapToDomain() *domain.ProjectMember {
	return &domain.ProjectMember{
		ProjectID: m.ProjectID,
		UserID:    m.UserID,
		RoleID:    m.RoleID,
		CreatedAt: m.CreatedAt,
	}
}

func ToDomainUsers(models []*UserModel) []*domain.User {
	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = m.MapToDomain()
	}
	return users
}

func ToDomainProjects(models []*ProjectModel) []*domain.Project {
	projects := make([]*domain.Project, len(models))
	for i, m := range models {
		projects[i] = m.MapToDomain()
	}
	return projects
}

func ToDomainProjectMembers(models []*ProjectMemberModel) []*domain.ProjectMember {
	members := make([]*domain.ProjectMember, len(models))
	for i, m := range models {
		members[i] = m.MapToDomain()
	}
	return members
}

func ToDomainSessions(models []*SessionModel) []*domain.Session {
	sessions := make([]*domain.Session, len(models))
	for i, m := range models {
		sessions[i] = m.MapToDomain()
	}
	return sessions
}

func ToDomainRefreshTokens(models []*RefreshTokenModel) []*domain.RefreshToken {
	tokens := make([]*domain.RefreshToken, len(models))
	for i, m := range models {
		tokens[i] = m.MapToDomain()
	}
	return tokens
}

func ToDomainRoles(models []*RoleModel) []*domain.Role {
	roles := make([]*domain.Role, len(models))
	for i, m := range models {
		roles[i] = m.MapToDomain()
	}
	return roles
}

func ToDomainPermissions(models []*PermissionModel) []*domain.Permission {
	perms := make([]*domain.Permission, len(models))
	for i, m := range models {
		perms[i] = m.MapToDomain()
	}
	return perms
}
