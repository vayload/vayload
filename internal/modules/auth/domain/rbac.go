package domain

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/shared/snowflake"
)

// Role represents a system or custom role.
type Role struct {
	ID          snowflake.ID `json:"id,string"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	IsSystem    bool         `json:"is_system"`
	CreatedAt   time.Time    `json:"created_at"`
}

// Permission represents an action on a resource.
type Permission struct {
	ID        snowflake.ID `json:"id,string"`
	Action    string       `json:"action"`
	Resource  string       `json:"resource"`
	CreatedAt time.Time    `json:"created_at"`
}

// UserRole links a user to a role in a specific project.
type UserRole struct {
	UserID    snowflake.ID `json:"user_id,string"`
	RoleID    snowflake.ID `json:"role_id,string"`
	ProjectID *snowflake.ID `json:"project_id,string,omitempty"`
}

type RbacRepository interface {
	CreateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id snowflake.ID) error
	FindRoleByName(ctx context.Context, name string) (*Role, error)

	CreatePermission(ctx context.Context, perm *Permission) error
	FindPermission(ctx context.Context, action, resource string) (*Permission, error)

	AssignRoleToUser(ctx context.Context, userID, roleID snowflake.ID, projectID *snowflake.ID) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID snowflake.ID, projectID *snowflake.ID) error
	GetUserRoles(ctx context.Context, userID snowflake.ID, projectID *snowflake.ID) ([]*Role, error)

	AttachPermissionToRole(ctx context.Context, roleID, permID snowflake.ID) error
	DetachPermissionFromRole(ctx context.Context, roleID, permID snowflake.ID) error
	GetRolePermissions(ctx context.Context, roleID snowflake.ID) ([]*Permission, error)
}
