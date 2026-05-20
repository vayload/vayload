package persistence

import (
	"context"
	"time"

	"github.com/vayload/vayload/internal/modules/auth/domain"
	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type ProjectRepository struct {
	database connection.DatabaseConnection
}

func NewProjectRepository(database connection.DatabaseConnection) *ProjectRepository {
	return &ProjectRepository{database: database}
}

func (repo *ProjectRepository) projects() connection.QueryBuilder {
	return repo.database.From("projects")
}

func (repo *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	data := map[string]any{
		"id":         project.ID,
		"name":       project.Name,
		"slug":       project.Slug,
		"owner_id":   project.OwnerID,
		"locale":     project.Locale,
		"created_at": project.CreatedAt,
	}
	return repo.projects().InsertOne(data).Exec(ctx)
}

func (repo *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	data := map[string]any{
		"name":       project.Name,
		"slug":       project.Slug,
		"owner_id":   project.OwnerID,
		"locale":     project.Locale,
		"updated_at": time.Now().UTC(),
	}
	return repo.projects().Where("id", "=", project.ID).UpdateOne(data).Exec(ctx)
}

func (repo *ProjectRepository) Delete(ctx context.Context, id snowflake.ID) error {
	return repo.projects().Where("id", "=", id).Delete().Exec(ctx)
}

func (repo *ProjectRepository) FindByID(ctx context.Context, id snowflake.ID) (*domain.Project, error) {
	var model ProjectModel
	if err := repo.projects().Where("id", "=", id).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *ProjectRepository) FindBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	var model ProjectModel
	if err := repo.projects().Where("slug", "=", slug).Get(ctx, &model); err != nil {
		return nil, selectQueryError(err)
	}
	return model.MapToDomain(), nil
}

func (repo *ProjectRepository) ListByOwner(ctx context.Context, ownerID snowflake.ID) ([]*domain.Project, error) {
	var models []*ProjectModel
	err := repo.projects().Where("owner_id", "=", ownerID).Select().Get(ctx, &models)
	if err != nil {
		return nil, err
	}

	return ToDomainProjects(models), nil
}

func (repo *ProjectRepository) AddMember(ctx context.Context, member *domain.ProjectMember) error {
	data := map[string]any{
		"project_id": member.ProjectID,
		"user_id":    member.UserID,
		"role_id":    member.RoleID,
		"created_at": member.CreatedAt,
	}
	return repo.database.From("project_members").InsertOne(data).Exec(ctx)
}

func (repo *ProjectRepository) RemoveMember(ctx context.Context, projectID, userID snowflake.ID) error {
	return repo.database.From("project_members").
		Where("project_id", "=", projectID).
		Where("user_id", "=", userID).
		Delete().Exec(ctx)
}

func (repo *ProjectRepository) ListMembers(ctx context.Context, projectID snowflake.ID) ([]*domain.ProjectMember, error) {
	var results []struct {
		ProjectID snowflake.ID `db:"project_id"`
		UserID    snowflake.ID `db:"user_id"`
		RoleID    snowflake.ID `db:"role_id"`
		CreatedAt time.Time    `db:"created_at"`
	}
	if err := repo.database.From("project_members").Where("project_id", "=", projectID).Select().Get(ctx, &results); err != nil {
		return nil, err
	}
	members := make([]*domain.ProjectMember, len(results))
	for i, r := range results {
		members[i] = &domain.ProjectMember{
			ProjectID: r.ProjectID,
			UserID:    r.UserID,
			RoleID:    r.RoleID,
			CreatedAt: r.CreatedAt,
		}
	}
	return members, nil
}

var _ domain.ProjectRepository = (*ProjectRepository)(nil)
