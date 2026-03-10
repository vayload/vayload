package handlers

import (
	"time"

	"github.com/vayload/vayload/internal/modules/storage/services"
	"github.com/vayload/vayload/internal/modules/storage/services/dtos"
	"github.com/vayload/vayload/internal/vayload"
)

type StorageHttpHandler struct {
	service *services.StorageService
}

func NewStorageHttpHandler(service *services.StorageService) *StorageHttpHandler {
	return &StorageHttpHandler{service: service}
}

func (h *StorageHttpHandler) Upload(req vayload.HttpRequest, res vayload.HttpResponse) error {
	file, err := req.File("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	ownerID := req.FormData("owner_id")[0]
	projectID := req.FormData("project_id")[0]
	folderID := req.FormData("folder_id")[0] // Changed to folderID

	fileObj, err := h.service.Upload(req.Context(), ownerID, projectID, folderID, file.Filename, file.Header.Get("Content-Type"), src)
	if err != nil {
		return err
	}

	return res.JSON(fileObj)
}

func (h *StorageHttpHandler) Download(req vayload.HttpRequest, res vayload.HttpResponse) error {
	id := req.GetParam("id")
	reader, _, err := h.service.Get(req.Context(), id)
	if err != nil {
		return err
	}
	defer reader.Close()

	return res.Stream(reader)
}

func (h *StorageHttpHandler) Sign(req vayload.HttpRequest, res vayload.HttpResponse) error {
	id := req.GetParam("id")
	expiry := 1 * time.Hour // Default expiry

	token, err := h.service.Sign(req.Context(), id, expiry)
	if err != nil {
		return err
	}

	return res.JSON(map[string]any{
		"token": token,
		"url":   "/storage/files/" + id + "?token=" + token,
	})
}

func (h *StorageHttpHandler) CreateFolder(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input dtos.FolderCreateInput
	if err := req.ParseBody(&input); err != nil {
		return err
	}

	folder, err := h.service.CreateFolder(req.Context(), input)
	if err != nil {
		return err
	}

	return res.JSON(folder)
}

func (h *StorageHttpHandler) GetFolderContents(req vayload.HttpRequest, res vayload.HttpResponse) error {
	projectID := req.GetQuery("project_id")
	folderIDStr := req.GetParam("id")
	var folderID *string
	if folderIDStr != "" && folderIDStr != "null" && folderIDStr != "root" {
		folderID = &folderIDStr
	}

	contents, err := h.service.GetFolderContents(req.Context(), projectID, folderID)
	if err != nil {
		return err
	}

	return res.JSON(contents)
}

func (h *StorageHttpHandler) Rename(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input dtos.RenameInput
	if err := req.ParseBody(&input); err != nil {
		return err
	}

	if err := h.service.Rename(req.Context(), input); err != nil {
		return err
	}

	return res.Status(200).JSON(map[string]string{"message": "renamed"})
}

func (h *StorageHttpHandler) Move(req vayload.HttpRequest, res vayload.HttpResponse) error {
	var input dtos.MoveInput
	if err := req.ParseBody(&input); err != nil {
		return err
	}

	if err := h.service.Move(req.Context(), input); err != nil {
		return err
	}

	return res.Status(200).JSON(map[string]string{"message": "moved"})
}

func (h *StorageHttpHandler) Delete(req vayload.HttpRequest, res vayload.HttpResponse) error {
	id := req.GetParam("id")
	itemType := req.GetQuery("type", "file")

	if err := h.service.Delete(req.Context(), id, itemType); err != nil {
		return err
	}

	return res.Status(200).JSON(map[string]string{"message": "deleted"})
}

func (h *StorageHttpHandler) Search(req vayload.HttpRequest, res vayload.HttpResponse) error {
	query := req.GetQuery("q")
	projectID := req.GetQuery("project_id")
	category := req.GetQuery("category")

	var catPtr *string
	if category != "" {
		catPtr = &category
	}

	results, err := h.service.Search(req.Context(), dtos.SearchInput{
		Query:     query,
		ProjectID: projectID,
		Category:  catPtr,
	})
	if err != nil {
		return err
	}

	return res.JSON(results)
}

func (s *StorageHttpHandler) HttpRoutes() []vayload.HttpRoutesGroup {
	return []vayload.HttpRoutesGroup{
		{
			Prefix: "/storage",
			Routes: []vayload.HttpRoute{
				{
					Path:    "/files/upload",
					Method:  vayload.HttpPost,
					Handler: s.Upload,
				},
				{
					Path:    "/files/:id",
					Method:  vayload.HttpGet,
					Handler: s.Download,
				},
				{
					Path:    "/files/:id/sign",
					Method:  vayload.HttpGet,
					Handler: s.Sign,
				},
				{
					Path:    "/folders",
					Method:  vayload.HttpPost,
					Handler: s.CreateFolder,
				},
				{
					Path:    "/folders/:id/contents",
					Method:  vayload.HttpGet,
					Handler: s.GetFolderContents,
				},
				{
					Path:    "/rename",
					Method:  vayload.HttpPatch,
					Handler: s.Rename,
				},
				{
					Path:    "/move",
					Method:  vayload.HttpPatch,
					Handler: s.Move,
				},
				{
					Path:    "/delete/:id",
					Method:  vayload.HttpDelete,
					Handler: s.Delete,
				},
				{
					Path:    "/search",
					Method:  vayload.HttpGet,
					Handler: s.Search,
				},
			},
		},
	}
}
