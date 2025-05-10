package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Data struct {
	Entries []*DataEntry `json:"entries"`
}

type DataEntry struct {
	ID         int           `json:"id"`
	Timestamp  time.Time     `json:"timestamp"`
	MinioLink  string        `json:"minio_link"`
	YOURLSLink string        `json:"yourls_link"`
	Expiry     time.Duration `json:"expiry"`
}

var (
	storageDir  = ".data"
	storagePath string
	storageFile *os.File
	currentData *Data
)

func Init(expiry time.Duration) error {
	if err := prepareStorage(); err != nil {
		return fmt.Errorf("failed to prepare storage: %w", err)
	}

	if err := loadData(); err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	removeExpiredEntries(expiry)

	return nil
}

func WriteEntry(entry *DataEntry) error {
	if entry == nil {
		return errors.New("entry cannot be nil")
	}

	if err := entry.validate(); err != nil {
		return fmt.Errorf("entry validation failed: %w", err)
	}

	entry.ID = findLatestID() + 1
	entry.Timestamp = time.Now()
	currentData.Entries = append(currentData.Entries, entry)

	return writeData()
}

func Read() (*Data, error) {
	return currentData, nil
}

func findLatestID() int {
	latest := 0
	for _, e := range currentData.Entries {
		if e.ID > latest {
			latest = e.ID
		}
	}
	return latest
}

func (e *DataEntry) validate() error {
	if e.MinioLink == "" {
		return errors.New("minio link cannot be empty")
	}
	if e.YOURLSLink == "" {
		return errors.New("yourls link cannot be empty")
	}
	return nil
}

func removeExpiredEntries(expiry time.Duration) {
	valid := make([]*DataEntry, 0, len(currentData.Entries))
	for _, e := range currentData.Entries {
		if expiry == -1 || time.Since(e.Timestamp) <= expiry {
			valid = append(valid, e)
		}
	}
	currentData.Entries = valid
}

func prepareStorage() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to locate executable: %w", err)
	}

	storageDir = filepath.Join(filepath.Dir(exePath), storageDir)
	storagePath = filepath.Join(storageDir, ".minyls.data.json")

	err = os.MkdirAll(storageDir, 0700)
	if err != nil {
		return fmt.Errorf("failed to create storage dir: %w", err)
	}

	storageFile, err = os.OpenFile(storagePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open storage file: %w", err)
	}

	return nil
}

func loadData() error {
	if _, err := storageFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek failed: %w", err)
	}

	decoder := json.NewDecoder(storageFile)
	d := &Data{}

	err := decoder.Decode(d)
	if errors.Is(err, io.EOF) {
		d.Entries = []*DataEntry{}
	} else if err != nil {
		return fmt.Errorf("decode failed: %w", err)
	}

	currentData = d
	return nil
}

func writeData() error {
	if err := storageFile.Truncate(0); err != nil {
		return fmt.Errorf("truncate failed: %w", err)
	}
	if _, err := storageFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek failed: %w", err)
	}

	encoder := json.NewEncoder(storageFile)
	if err := encoder.Encode(currentData); err != nil {
		return fmt.Errorf("encode failed: %w", err)
	}

	if err := storageFile.Sync(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	return nil
}
