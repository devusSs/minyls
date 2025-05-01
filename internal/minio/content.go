package minio

import (
	"fmt"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

// does this contain enough content types?
func findContentType(filePath string) (string, error) {
	mime, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", err
	}

	return mime.String(), nil
}

// randomizeFileName takes in a full file path and returns a randomized
// file name (uuid) plus the extension of the file for the purpose of serving
// them via browsers (content detection).
func randomizeFileName(filePath string) (string, error) {
	file := filepath.Base(filePath)
	ext := filepath.Ext(file)

	uid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("could not create uuid: %w", err)
	}

	return uid.String() + ext, nil
}
