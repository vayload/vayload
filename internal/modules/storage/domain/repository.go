package domain

import "context"

type StorageRepository interface {
	CreateFolder(ctx context.Context, folder *FolderObject) error
	GetFolderByID(ctx context.Context, id string) (*FolderObject, error)
	GetFolderByPath(ctx context.Context, projectID, path string) (*FolderObject, error)
	ListFolders(ctx context.Context, projectID string, parentID *string) ([]FolderObject, error)
	UpdateFolder(ctx context.Context, folder *FolderObject) error
	DeleteFolder(ctx context.Context, id string) error

	CreateFile(ctx context.Context, file *FileObject) error
	GetFileByID(ctx context.Context, id string) (*FileObject, error)
	ListFiles(ctx context.Context, projectID string, folderID *string) ([]FileObject, error)
	UpdateFile(ctx context.Context, file *FileObject) error
	DeleteFile(ctx context.Context, id string) error
}
