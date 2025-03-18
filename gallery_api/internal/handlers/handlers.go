package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/nixoncode/gallery_api/internal/storage"
	"github.com/nixoncode/gallery_api/internal/utils"
)

type Handlers struct {
	storage     *storage.Storage
	maxFileSize int64
}

func NewHandlers(storage *storage.Storage, maxFileSize int64) *Handlers {
	return &Handlers{storage: storage, maxFileSize: maxFileSize}
}

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (h *Handlers) UploadImage(c echo.Context) error {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Error: "Unable to get file"})
	}
	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Unable to open file"})
	}
	defer file.Close()

	if fileHeader.Size > h.maxFileSize {
		return c.JSON(http.StatusBadRequest, Response{Error: "File size exceeds limit"})
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Could not read file"})
	}

	fileType := http.DetectContentType(buffer)
	if fileType != "image/jpeg" && fileType != "image/png" {
		return c.JSON(http.StatusBadRequest, Response{Error: "File type not supported"})
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Could not reset file reader"})
	}

	filename := fileHeader.Filename
	description := c.FormValue("description")
	newFilename := c.FormValue("newFilename")

	tempFilePath := filepath.Join(h.storage.BaseDir(), filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Failed to create temp file"})
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Failed to write to temp file"})
	}

	metadata, err := utils.ExtractMetadata(tempFilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Failed to extract metadata"})
	}

	if newFilename != "" {
		newFilePath := filepath.Join(h.storage.BaseDir(), newFilename)
		err = utils.RenameFile(tempFilePath, newFilePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Response{Error: "Failed to rename file"})
		}
		filename = newFilename
	}

	err = h.storage.SaveImage(filename, description, metadata, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Error: "Unable to save file"})
	}

	return c.JSON(http.StatusOK, Response{Message: "File uploaded successfully"})
}

func (h *Handlers) GetImage(c echo.Context) error {
	filename := c.QueryParam("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, Response{Error: "Filename is required"})
	}

	file, err := h.storage.GetImage(filename)
	if err != nil {
		return c.JSON(http.StatusNotFound, Response{Error: "File not found"})
	}
	defer file.Close()

	return c.Stream(http.StatusOK, "image/jpeg", file)
}

func (h *Handlers) GetImageDetails(c echo.Context) error {
	filename := c.QueryParam("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, Response{Error: "Filename is required"})
	}

	image, err := h.storage.GetImageDetails()
	if err != nil {
		return c.JSON(http.StatusNotFound, Response{Error: "File not found"})
	}

	return c.JSON(http.StatusOK, Response{Data: image})
}
