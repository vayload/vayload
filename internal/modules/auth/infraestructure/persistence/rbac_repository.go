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

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type RbacRepository struct {
	database connection.DatabaseConnection
}

func NewRbacRepository(db connection.DatabaseConnection) *RbacRepository {
	return &RbacRepository{database: db}
}

func (repo *RbacRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	return repo.database.From("roles").InsertOne(map[string]any{
		"id":          role.ID,
		"name":        role.Name,
		"description": role.Description,
		"is_system":   role.IsSystem,
		"created_at":  role.CreatedAt,
	}).Exec(ctx)
}

func (repo *RbacRepository) DeleteRole(ctx context.Context, id snowflake.ID) error {
	return repo.database.From("roles").Where("id", "=", id).Delete().Exec(ctx)
}

func (repo *RbacRepository) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	var model RoleModel
	err := repo.database.From("roles").Where("name", "=", name).Select().Get(ctx, &model)
	if err != nil {
		return nil, err
	}
	return model.MapToDomain(), nil
}

func (repo *RbacRepository) CreatePermission(ctx context.Context, perm *domain.Permission) error {
	return repo.database.From("permissions").InsertOne(map[string]any{
		"id":         perm.ID,
		"action":     perm.Action,
		"resource":   perm.Resource,
		"created_at": perm.CreatedAt,
	}).Exec(ctx)
}

func (repo *RbacRepository) FindPermission(ctx context.Context, action, resource string) (*domain.Permission, error) {
	var model PermissionModel
	err := repo.database.From("permissions").
		Where("action", "=", action).
		Where("resource", "=", resource).
		Select().Get(ctx, &model)
	if err != nil {
		return nil, err
	}

	return model.MapToDomain(), nil
}

func (repo *RbacRepository) AssignRoleToUser(ctx context.Context, userID, roleID snowflake.ID, projectID *snowflake.ID) error {
	values := map[string]any{
		"user_id":    userID,
		"role_id":    roleID,
		"project_id": projectID,
	}
	return repo.database.From("user_roles").InsertOne(values).Exec(ctx)
}

func (repo *RbacRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID snowflake.ID, projectID *snowflake.ID) error {
	qb := repo.database.From("user_roles").Where("user_id", "=", userID).Where("role_id", "=", roleID)
	if projectID != nil {
		qb = qb.Where("project_id", "=", *projectID)
	} else {
		qb = qb.Where("project_id", "IS", nil)
	}
	return qb.Delete().Exec(ctx)
}

func (repo *RbacRepository) GetUserRoles(ctx context.Context, userID snowflake.ID, projectID *snowflake.ID) ([]*domain.Role, error) {
	var models []*RoleModel
	qb := repo.database.From("roles").Join("user_roles", "roles.id", "=", "user_roles.role_id")
	qb = qb.Where("user_roles.user_id", "=", userID)
	if projectID != nil {
		qb = qb.Where("user_roles.project_id", "=", *projectID)
	} else {
		qb = qb.Where("user_roles.project_id", "IS", nil)
	}

	if err := qb.Select("roles.*").Get(ctx, &models); err != nil {
		return nil, err
	}
	return ToDomainRoles(models), nil
}

func (repo *RbacRepository) AttachPermissionToRole(ctx context.Context, roleID, permID snowflake.ID) error {
	return repo.database.From("role_permissions").InsertOne(map[string]any{
		"role_id":       roleID,
		"permission_id": permID,
	}).Exec(ctx)
}

func (repo *RbacRepository) DetachPermissionFromRole(ctx context.Context, roleID, permID snowflake.ID) error {
	return repo.database.From("role_permissions").
		Where("role_id", "=", roleID).
		Where("permission_id", "=", permID).
		Delete().Exec(ctx)
}

func (repo *RbacRepository) GetRolePermissions(ctx context.Context, roleID snowflake.ID) ([]*domain.Permission, error) {
	var models []*PermissionModel
	err := repo.database.From("permissions").
		Join("role_permissions", "permissions.id", "=", "role_permissions.permission_id").
		Where("role_permissions.role_id", "=", roleID).
		Select("permissions.*").
		Get(ctx, &models)
	if err != nil {
		return nil, err
	}
	return ToDomainPermissions(models), nil
}

var _ domain.RbacRepository = (*RbacRepository)(nil)
