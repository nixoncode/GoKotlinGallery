package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/nixoncode/gallery_api/internal/database/sqlc"
	"github.com/nixoncode/gallery_api/internal/models"
	"github.com/sqlc-dev/pqtype"
)

type Storage struct {
	baseDir string
	db      *sql.DB
	queries *sqlc.Queries
}

func NewStorage(baseDir string, db *sql.DB) *Storage {
	return &Storage{
		baseDir: baseDir,
		db:      db,
		queries: sqlc.New(db),
	}
}

func (s *Storage) SaveImage(filename string, description string, metadata map[string]interface{}, file io.Reader) error {
	filePath := filepath.Join(s.baseDir, filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = s.queries.CreateImage(context.Background(), sqlc.CreateImageParams{
		Filename:    filename,
		Description: sql.NullString{String: description, Valid: description != ""},
		Metadata:    pqtype.NullRawMessage{RawMessage: json.RawMessage(metadataJSON), Valid: true},
	})

	return err
}

func (s *Storage) GetImage(filename string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.baseDir, filename)
	file, err := os.Open(filePath)
	return file, err
}

func (s *Storage) GetImageDetails(filename string) (*models.Image, error) {
	image, err := s.queries.GetImageDetails(context.Background(), filename)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	if image.Metadata.Valid {
		err := json.Unmarshal(image.Metadata.RawMessage, &metadata)
		if err != nil {
			return nil, err
		}
	}

	return &models.Image{
		Filename:    image.Filename,
		Description: image.Description.String,
		Metadata:    metadata,
	}, nil
}

func (s *Storage) BaseDir() string {
	return s.baseDir
}
