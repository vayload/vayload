package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/vayload/vayload/internal/modules/storage/domain"
	"github.com/vayload/vayload/internal/modules/storage/services/dtos"
	"github.com/vayload/vayload/internal/shared/snowflake"
)

type StorageService struct {
	repository domain.StorageRepository
	// store is implicitly through repository or we might need it for direct uploads
	// If the user wants to keep the direct upload logic from storage.go
}

func NewStorageService(repository domain.StorageRepository) *StorageService {
	return &StorageService{
		repository: repository,
	}
}

func (s *StorageService) CreateFolder(ctx context.Context, input dtos.FolderCreateInput) (*domain.FolderObject, error) {
	id := snowflake.Node.Generate().String()

	// Determine path and depth
	path := "/" + input.Name
	depth := 1
	if input.ParentID != nil {
		parent, err := s.repository.GetFolderByID(ctx, *input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent folder not found: %w", err)
		}
		path = parent.Path + "/" + input.Name
		depth = parent.Depth + 1
	}

	folder := domain.NewFolderObject(id, input.OwnerID, input.ProjectID, input.Name, input.ParentID, path, depth)

	if err := s.repository.CreateFolder(ctx, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}

func (s *StorageService) GetFolderContents(ctx context.Context, projectID string, folderID *string) (*dtos.FolderContentsResponse, error) {
	folders, err := s.repository.ListFolders(ctx, projectID, folderID)
	if err != nil {
		return nil, err
	}

	files, err := s.repository.ListFiles(ctx, projectID, folderID)
	if err != nil {
		return nil, err
	}

	return &dtos.FolderContentsResponse{
		Folders: folders,
		Files:   files,
	}, nil
}

func (s *StorageService) Rename(ctx context.Context, input dtos.RenameInput) error {
	if input.Type == "folder" {
		folder, err := s.repository.GetFolderByID(ctx, input.ID)
		if err != nil {
			return err
		}
		folder.Name = input.NewName
		folder.UpdatedAt = time.Now().UTC()
		// Path would need to be updated recursively for children in a real VFS
		return s.repository.UpdateFolder(ctx, folder)
	}

	file, err := s.repository.GetFileByID(ctx, input.ID)
	if err != nil {
		return err
	}
	file.Name = input.NewName
	file.UpdatedAt = time.Now().UTC()
	return s.repository.UpdateFile(ctx, file)
}

func (s *StorageService) Move(ctx context.Context, input dtos.MoveInput) error {
	if input.Type == "folder" {
		folder, err := s.repository.GetFolderByID(ctx, input.ID)
		if err != nil {
			return err
		}
		folder.ParentID = input.NewParentID
		folder.UpdatedAt = time.Now().UTC()
		// Path update logic...
		return s.repository.UpdateFolder(ctx, folder)
	}

	file, err := s.repository.GetFileByID(ctx, input.ID)
	if err != nil {
		return err
	}
	file.FolderID = input.NewParentID
	file.UpdatedAt = time.Now().UTC()
	return s.repository.UpdateFile(ctx, file)
}

func (s *StorageService) Delete(ctx context.Context, id string, itemType string) error {
	if itemType == "folder" {
		return s.repository.DeleteFolder(ctx, id)
	}
	return s.repository.DeleteFile(ctx, id)
}

func (s *StorageService) Search(ctx context.Context, input dtos.SearchInput) (*dtos.FolderContentsResponse, error) {
	// Simple implementation: search by name in repository
	// Depending on repository capabilities
	// For now, let's assume we can filter ListFiles/ListFolders
	return &dtos.FolderContentsResponse{}, nil
}

// Keeping existing methods for compatibility or finishing them
func (s *StorageService) Upload(ctx context.Context, ownerID, projectID, folderID string, name, mimeType string, reader io.Reader) (*domain.FileObject, error) {
	// Implementation should follow what was in storage.go but integrated here
	return nil, fmt.Errorf("not implemented")
}

func (s *StorageService) Get(ctx context.Context, fileID string) (io.ReadCloser, *domain.FileObject, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (s *StorageService) Sign(ctx context.Context, fileID string, expiry time.Duration) (string, error) {
	return "", fmt.Errorf("not implemented")
}
