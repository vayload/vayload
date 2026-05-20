package domain

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

// Project represents a project owned by a user.
type Project struct {
	ID        snowflake.ID   `json:"id,string"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	OwnerID   snowflake.ID   `json:"owner_id,string"`
	Settings  map[string]any `json:"settings,omitempty"`
	Locale    string         `json:"locale"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty"`
}

// ProjectMember represents a user participating in a project with a specific role.
type ProjectMember struct {
	ProjectID snowflake.ID `json:"project_id,string"`
	UserID    snowflake.ID `json:"user_id,string"`
	RoleID    snowflake.ID `json:"role_id,string"`
	CreatedAt time.Time    `json:"created_at"`
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id snowflake.ID) error
	FindByID(ctx context.Context, id snowflake.ID) (*Project, error)
	FindBySlug(ctx context.Context, slug string) (*Project, error)
	ListByOwner(ctx context.Context, ownerID snowflake.ID) ([]*Project, error)

	AddMember(ctx context.Context, member *ProjectMember) error
	RemoveMember(ctx context.Context, projectID, userID snowflake.ID) error
	ListMembers(ctx context.Context, projectID snowflake.ID) ([]*ProjectMember, error)
}
