package repositories

import (
	"context"
	"encoding/json"

	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/storage/domain"
)

type SQLStorageRepository struct {
	db            connection.DatabaseConnection
	fileObjects   connection.QueryBuilder
	folderObjects connection.QueryBuilder
}

func NewSQLStorageRepository(db connection.DatabaseConnection) *SQLStorageRepository {
	return &SQLStorageRepository{
		db:            db,
		fileObjects:   db.From("files"),
		folderObjects: db.From("folders"),
	}
}

// Folder Operations

func (r *SQLStorageRepository) CreateFolder(ctx context.Context, folder *domain.FolderObject) error {
	return r.folderObjects.InsertOne(map[string]any{
		"id":         folder.ID,
		"owner_id":   folder.OwnerID,
		"project_id": folder.ProjectID,
		"parent_id":  folder.ParentID,
		"name":       folder.Name,
		"path":       folder.Path,
		"depth":      folder.Depth,
		"created_at": folder.CreatedAt,
		"updated_at": folder.UpdatedAt,
	}).Exec(ctx)
}

func (r *SQLStorageRepository) GetFolderByID(ctx context.Context, id string) (*domain.FolderObject, error) {
	var folder domain.FolderObject
	err := r.folderObjects.Where("id", "=", id).First(ctx, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *SQLStorageRepository) GetFolderByPath(ctx context.Context, projectID, path string) (*domain.FolderObject, error) {
	var folder domain.FolderObject
	err := r.folderObjects.
		Where("project_id", "=", projectID).
		Where("path", "=", path).
		First(ctx, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *SQLStorageRepository) ListFolders(ctx context.Context, projectID string, parentID *string) ([]domain.FolderObject, error) {
	var folders []domain.FolderObject
	query := r.folderObjects.Where("project_id", "=", projectID)

	if parentID == nil {
		query = query.WhereNull("parent_id")
	} else {
		query = query.Where("parent_id", "=", *parentID)
	}

	err := query.Get(ctx, &folders)
	return folders, err
}

func (r *SQLStorageRepository) UpdateFolder(ctx context.Context, folder *domain.FolderObject) error {
	return r.folderObjects.Where("id", "=", folder.ID).UpdateOne(map[string]any{
		"name":       folder.Name,
		"parent_id":  folder.ParentID,
		"path":       folder.Path,
		"depth":      folder.Depth,
		"updated_at": folder.UpdatedAt,
	}).Exec(ctx)
}

func (r *SQLStorageRepository) DeleteFolder(ctx context.Context, id string) error {
	return r.folderObjects.Where("id", "=", id).Delete().Exec(ctx)
}

// File Operations

func (r *SQLStorageRepository) CreateFile(ctx context.Context, file *domain.FileObject) error {
	metadata, _ := json.Marshal(file.Metadata)

	return r.fileObjects.InsertOne(map[string]any{
		"id":           file.ID,
		"owner_id":     file.OwnerID,
		"project_id":   file.ProjectID,
		"folder_id":    file.FolderID,
		"name":         file.Name,
		"mime_type":    file.MimeType,
		"category":     file.Category,
		"size":         file.Size,
		"provider":     file.Provider,
		"provider_key": file.ProviderKey,
		"metadata":     string(metadata),
		"created_at":   file.CreatedAt,
		"updated_at":   file.UpdatedAt,
	}).Exec(ctx)
}

func (r *SQLStorageRepository) GetFileByID(ctx context.Context, id string) (*domain.FileObject, error) {
	var file domain.FileObject
	err := r.fileObjects.Where("id", "=", id).Get(ctx, &file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *SQLStorageRepository) ListFiles(ctx context.Context, projectID string, folderID *string) ([]domain.FileObject, error) {
	var files []domain.FileObject
	query := r.fileObjects.Where("project_id", "=", projectID)

	if folderID == nil {
		query = query.WhereNull("folder_id")
	} else {
		query = query.Where("folder_id", "=", *folderID)
	}

	err := query.Get(ctx, &files)
	return files, err
}

func (r *SQLStorageRepository) UpdateFile(ctx context.Context, file *domain.FileObject) error {
	metadata, _ := json.Marshal(file.Metadata)

	return r.fileObjects.Where("id", "=", file.ID).UpdateOne(map[string]any{
		"name":       file.Name,
		"folder_id":  file.FolderID,
		"updated_at": file.UpdatedAt,
		"metadata":   string(metadata),
	}).Exec(ctx)
}

func (r *SQLStorageRepository) DeleteFile(ctx context.Context, id string) error {
	return r.fileObjects.Where("id", "=", id).Delete().Exec(ctx)
}
