package repositories

import (
	"context"
	"encoding/json"

	"github.com/vayload/vayload/internal/modules/database/connection"
	"github.com/vayload/vayload/internal/modules/storage/domain"
)

type SQLStorageRepository struct {
	db connection.DatabaseConnection
}

func NewSQLStorageRepository(db connection.DatabaseConnection) *SQLStorageRepository {
	return &SQLStorageRepository{db: db}
}

// Folder Operations

func (r *SQLStorageRepository) CreateFolder(ctx context.Context, folder *domain.FolderObject) error {
	return r.db.From("folders").Insert(map[string]any{
		"id":         folder.ID,
		"owner_id":   folder.OwnerID,
		"project_id": folder.ProjectID,
		"parent_id":  folder.ParentID,
		"name":       folder.Name,
		"path":       folder.Path,
		"depth":      folder.Depth,
		"created_at": folder.CreatedAt,
		"updated_at": folder.UpdatedAt,
	}).Exec()
}

func (r *SQLStorageRepository) GetFolderByID(ctx context.Context, id string) (*domain.FolderObject, error) {
	var folder domain.FolderObject
	err := r.db.From("folders").Where("id", "=", id).Get(&folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *SQLStorageRepository) GetFolderByPath(ctx context.Context, projectID, path string) (*domain.FolderObject, error) {
	var folder domain.FolderObject
	err := r.db.From("folders").
		Where("project_id", "=", projectID).
		Where("path", "=", path).
		Get(&folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *SQLStorageRepository) ListFolders(ctx context.Context, projectID string, parentID *string) ([]domain.FolderObject, error) {
	var folders []domain.FolderObject
	query := r.db.From("folders").Where("project_id", "=", projectID)
	
	if parentID == nil {
		query = query.WhereNull("parent_id")
	} else {
		query = query.Where("parent_id", "=", *parentID)
	}

	err := query.GetAll(&folders)
	return folders, err
}

func (r *SQLStorageRepository) UpdateFolder(ctx context.Context, folder *domain.FolderObject) error {
	return r.db.From("folders").Where("id", "=", folder.ID).Update(map[string]any{
		"name":       folder.Name,
		"parent_id":  folder.ParentID,
		"path":       folder.Path,
		"depth":      folder.Depth,
		"updated_at": folder.UpdatedAt,
	}).Exec()
}

func (r *SQLStorageRepository) DeleteFolder(ctx context.Context, id string) error {
	return r.db.From("folders").Where("id", "=", id).Delete().Exec()
}

// File Operations

func (r *SQLStorageRepository) CreateFile(ctx context.Context, file *domain.FileObject) error {
	metadata, _ := json.Marshal(file.Metadata)
	
	return r.db.From("files").Insert(map[string]any{
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
	}).Exec()
}

func (r *SQLStorageRepository) GetFileByID(ctx context.Context, id string) (*domain.FileObject, error) {
	var file domain.FileObject
	err := r.db.From("files").Where("id", "=", id).Get(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *SQLStorageRepository) ListFiles(ctx context.Context, projectID string, folderID *string) ([]domain.FileObject, error) {
	var files []domain.FileObject
	query := r.db.From("files").Where("project_id", "=", projectID)

	if folderID == nil {
		query = query.WhereNull("folder_id")
	} else {
		query = query.Where("folder_id", "=", *folderID)
	}

	err := query.GetAll(&files)
	return files, err
}

func (r *SQLStorageRepository) UpdateFile(ctx context.Context, file *domain.FileObject) error {
	metadata, _ := json.Marshal(file.Metadata)

	return r.db.From("files").Where("id", "=", file.ID).Update(map[string]any{
		"name":       file.Name,
		"folder_id":  file.FolderID,
		"updated_at": file.UpdatedAt,
		"metadata":   string(metadata),
	}).Exec()
}

func (r *SQLStorageRepository) DeleteFile(ctx context.Context, id string) error {
	return r.db.From("files").Where("id", "=", id).Delete().Exec()
}
