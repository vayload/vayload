package dtos

import (
	"io"

	"github.com/vayload/vayload/internal/modules/storage/domain"
)

type FileUploadInput struct {
	File      io.ReadSeekCloser `json:"file"`
	Name      string            `json:"name"`
	MimeType  string            `json:"mime_type"`
	OwnerID   string            `json:"owner_id"`
	ProjectID string            `json:"project_id"`
	FolderID  *string           `json:"folder_id"`
}

type FolderCreateInput struct {
	Name      string  `json:"name"`
	OwnerID   string  `json:"owner_id"`
	ProjectID string  `json:"project_id"`
	ParentID  *string `json:"parent_id"`
}

type RenameInput struct {
	ID      string `json:"id"`
	NewName string `json:"new_name"`
	Type    string `json:"type"` // "file" or "folder"
}

type MoveInput struct {
	ID          string  `json:"id"`
	NewParentID *string `json:"new_parent_id"`
	Type        string  `json:"type"` // "file" or "folder"
}

type DeleteInput struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "file" or "folder"
}

type SearchInput struct {
	Query     string  `json:"query"`
	ProjectID string  `json:"project_id"`
	Category  *string `json:"category"`
}

type FolderContentsResponse struct {
	Folders []domain.FolderObject `json:"folders"`
	Files   []domain.FileObject   `json:"files"`
}
