package storage

import (
	"strings"
	"time"
)

type FileCategory string

const (
	CategoryImage        FileCategory = "image"
	CategoryVideo        FileCategory = "video"
	CategoryAudio        FileCategory = "audio"
	CategoryText         FileCategory = "text"
	CategoryDocument     FileCategory = "document"
	CategorySpreadsheet  FileCategory = "spreadsheet"
	CategoryPresentation FileCategory = "presentation"
	CategoryArchive      FileCategory = "archive"
	CategoryExecutable   FileCategory = "executable"
	CategoryFont         FileCategory = "font"
	CategoryOther        FileCategory = "other"
)

type FileObject struct {
	ID        string       `json:"id" db:"id"`
	OwnerID   string       `json:"owner_id" db:"owner_id"`
	ProjectID string       `json:"project_id" db:"project_id"`
	Name      string       `json:"name" db:"name"`
	MimeType  string       `json:"mime_type" db:"mime_type"`
	Category  FileCategory `json:"category" db:"category"`
	Size      int64        `json:"size" db:"size"`
	SHA256    [32]byte     `json:"sha256" db:"sha256"`
	Key       string       `json:"key" db:"key"`       // storage key (local o S3)
	Folder    string       `json:"folder" db:"folder"` // storage folder for UI
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

func (f FileObject) IsOwner(ownerId string) bool {
	return f.OwnerID == ownerId
}

func (f FileObject) IsCategory(c FileCategory) bool {
	return f.Category == c
}

func (f FileObject) IsMedia() bool {
	return f.Category == CategoryImage ||
		f.Category == CategoryVideo ||
		f.Category == CategoryAudio
}

func (f FileObject) Extension() string {
	if i := strings.LastIndexByte(f.Name, '.'); i >= 0 {
		return strings.ToLower(f.Name[i+1:])
	}
	return ""
}

func NewFileObject(id string, ownerID string, name string, mimeType string, key string) FileObject {
	now := time.Now().UTC()

	return FileObject{
		ID:        id,
		OwnerID:   ownerID,
		Name:      name,
		MimeType:  mimeType,
		Category:  DetectCategory(mimeType),
		Key:       key,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (f FileObject) WithContent(size int64, sha256 [32]byte) FileObject {
	f.Size = size
	f.SHA256 = sha256
	f.UpdatedAt = time.Now().UTC()

	return f
}

func (f FileObject) WithKey(key string) FileObject {
	f.Key = key
	f.UpdatedAt = time.Now().UTC()

	return f
}

func (f FileObject) WithFolder(folder string) FileObject {
	f.Folder = folder
	f.UpdatedAt = time.Now().UTC()

	return f
}

func DetectCategory(mime string) FileCategory {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return CategoryImage
	case strings.HasPrefix(mime, "video/"):
		return CategoryVideo
	case strings.HasPrefix(mime, "audio/"):
		return CategoryAudio
	case mime == "application/pdf":
		return CategoryDocument
	case strings.Contains(mime, "spreadsheet"):
		return CategorySpreadsheet
	case strings.Contains(mime, "presentation"):
		return CategoryPresentation
	case strings.Contains(mime, "zip"):
		return CategoryArchive
	default:
		return CategoryOther
	}
}
