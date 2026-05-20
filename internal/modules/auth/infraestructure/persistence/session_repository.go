package persistence

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type SessionRepository struct {
	database connection.DatabaseConnection
}

func NewSessionRepository(database connection.DatabaseConnection) *SessionRepository {
	return &SessionRepository{database: database}
}

func (repo *SessionRepository) sessions() connection.QueryBuilder {
	return repo.database.From("sessions")
}

func (repo *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	data := map[string]any{
		"id":           session.ID,
		"user_id":      session.UserID,
		"project_id":   session.ProjectID,
		"ip_address":   session.IPAddress,
		"user_agent":   session.UserAgent,
		"last_seen_at": session.LastSeenAt,
		"expires_at":   session.ExpiresAt,
		"created_at":   session.CreatedAt,
	}
	return repo.sessions().InsertOne(data).Exec(ctx)
}

func (repo *SessionRepository) Update(ctx context.Context, session *domain.Session) error {
	data := map[string]any{
		"last_seen_at": session.LastSeenAt,
		"expires_at":   session.ExpiresAt,
		"revoked_at":   session.RevokedAt,
	}
	return repo.sessions().Where("id", "=", session.ID).UpdateOne(data).Exec(ctx)
}

func (repo *SessionRepository) FindByID(ctx context.Context, id string) (*domain.Session, error) {
	var model SessionModel
	if err := repo.sessions().Where("id", "=", id).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *SessionRepository) ListByUser(ctx context.Context, userID snowflake.ID) ([]*domain.Session, error) {
	var models []*SessionModel
	err := repo.sessions().Where("user_id", "=", userID).Where("revoked_at", "IS", nil).Select().Get(ctx, &models)
	if err != nil {
		return nil, err
	}

	return ToDomainSessions(models), nil
}

func (repo *SessionRepository) Revoke(ctx context.Context, id string) error {
	return repo.sessions().Where("id", "=", id).UpdateOne(map[string]any{
		"revoked_at": time.Now().UTC(),
	}).Exec(ctx)
}

func (repo *SessionRepository) RevokeAllByUser(ctx context.Context, userID snowflake.ID) error {
	return repo.sessions().Where("user_id", "=", userID).UpdateOne(map[string]any{
		"revoked_at": time.Now().UTC(),
	}).Exec(ctx)
}

type RefreshTokenRepository struct {
	database connection.DatabaseConnection
}

func NewRefreshTokenRepository(database connection.DatabaseConnection) *RefreshTokenRepository {
	return &RefreshTokenRepository{database: database}
}

func (repo *RefreshTokenRepository) tokens() connection.QueryBuilder {
	return repo.database.From("refresh_tokens")
}

func (repo *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	data := map[string]any{
		"id":         token.ID,
		"token_hash": token.TokenHash,
		"user_id":    token.UserID,
		"family_id":  token.FamilyID,
		"session_id": token.SessionID,
		"parent_id":  token.ParentID,
		"expires_at": token.ExpiresAt,
		"created_at": token.CreatedAt,
	}
	return repo.tokens().InsertOne(data).Exec(ctx)
}

func (repo *RefreshTokenRepository) Update(ctx context.Context, token *domain.RefreshToken) error {
	data := map[string]any{
		"used_at":        token.UsedAt,
		"revoked_at":     token.RevokedAt,
		"revoked_reason": token.RevokedReason,
	}
	return repo.tokens().Where("id", "=", token.ID).UpdateOne(data).Exec(ctx)
}

func (repo *RefreshTokenRepository) FindByID(ctx context.Context, id string) (*domain.RefreshToken, error) {
	var model RefreshTokenModel
	if err := repo.tokens().Where("id", "=", id).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *RefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	var model RefreshTokenModel
	if err := repo.tokens().Where("token_hash", "=", hash).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *RefreshTokenRepository) RevokeFamily(ctx context.Context, familyID string, reason string) error {
	return repo.tokens().Where("family_id", "=", familyID).UpdateOne(map[string]any{
		"revoked_at":     time.Now().UTC(),
		"revoked_reason": reason,
	}).Exec(ctx)
}

var _ domain.SessionRepository = (*SessionRepository)(nil)
var _ domain.RefreshTokenRepository = (*RefreshTokenRepository)(nil)
