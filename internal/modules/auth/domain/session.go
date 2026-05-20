package domain

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

// Session represents an active user session.
type Session struct {
	ID        string       `json:"id"`
	UserID    snowflake.ID `json:"user_id,string"`
	ProjectID *snowflake.ID `json:"project_id,string,omitempty"`

	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`

	LastSeenAt time.Time  `json:"last_seen_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// RefreshToken represents a refresh token associated with a session.
type RefreshToken struct {
	ID        string       `json:"id"`
	TokenHash string       `json:"-"`
	UserID    snowflake.ID `json:"user_id,string"`

	FamilyID string  `json:"family_id"`
	SessionID string `json:"session_id,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`

	UsedAt        *time.Time `json:"used_at,omitempty"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	RevokedReason string     `json:"revoked_reason,omitempty"`

	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id string) (*Session, error)
	ListByUser(ctx context.Context, userID snowflake.ID) ([]*Session, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllByUser(ctx context.Context, userID snowflake.ID) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	Update(ctx context.Context, token *RefreshToken) error
	FindByID(ctx context.Context, id string) (*RefreshToken, error)
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	RevokeFamily(ctx context.Context, familyID string, reason string) error
}
