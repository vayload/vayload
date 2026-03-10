package domain

import (
	"strings"
	"time"
)

type Provider string

const (
	ProviderLocal Provider = "local"
	ProviderS3    Provider = "s3"
	ProviderR2    Provider = "r2"
	ProviderGCS   Provider = "gcs"
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

type FolderObject struct {
	ID             string    `json:"id"`
	OwnerID        string    `json:"owner_id"`
	ProjectID      string    `json:"project_id"`
	ParentID       *string   `json:"parent_id"`
	Name           string    `json:"name"`
	FileCount      int       `json:"file_count"`
	SubfolderCount int       `json:"subfolder_count"`
	Path           string    `json:"path"`
	Depth          int       `json:"depth"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FileObject struct {
	ID          string         `json:"id"`
	OwnerID     string         `json:"owner_id"`
	ProjectID   string         `json:"project_id"`
	FolderID    *string        `json:"folder_id"`
	Name        string         `json:"name"`
	MimeType    string         `json:"mime_type"`
	Category    FileCategory   `json:"category"`
	Size        int64          `json:"size"`
	Provider    Provider       `json:"provider"`
	ProviderKey string         `json:"provider_key"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func NewFolderObject(id, ownerID, projectID, name string, parentID *string, path string, depth int) FolderObject {
	now := time.Now().UTC()

	return FolderObject{
		ID:        id,
		OwnerID:   ownerID,
		ProjectID: projectID,
		ParentID:  parentID,
		Name:      name,
		Path:      path,
		Depth:     depth,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewFileObject(id, ownerID, projectID, name, mimeType string, size int64, provider Provider, providerKey string, metadata map[string]any) FileObject {
	now := time.Now().UTC()

	return FileObject{
		ID:          id,
		OwnerID:     ownerID,
		ProjectID:   projectID,
		Name:        name,
		MimeType:    mimeType,
		Category:    DetectCategory(mimeType),
		Size:        size,
		Provider:    provider,
		ProviderKey: providerKey,
		Metadata:    metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
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
