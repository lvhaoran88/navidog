package storage

import (
	"os"
	"path"
)

type LocalStorage struct {
	Path string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(filepath string) *LocalStorage {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &LocalStorage{
		Path: path.Join(wd, "Navidog", filepath),
	}
}

// Load
func (s *LocalStorage) Load() ([]byte, error) {
	return os.ReadFile(s.Path)
}

// Store
func (s *LocalStorage) Store(data []byte) error {
	// Create directory if it doesn't exist
	if _, err := os.Stat(path.Dir(s.Path)); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path.Dir(s.Path), os.ModePerm); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return os.WriteFile(s.Path, data, os.ModePerm)
}
