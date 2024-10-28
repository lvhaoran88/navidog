package storage

import "sync"

type ConnectionStorage struct {
	mutex   sync.Mutex
	storage *LocalStorage
}

func NewConnectionStorage() *ConnectionStorage {
	return &ConnectionStorage{
		storage: NewLocalStorage("connections.yaml"),
	}
}

// Read
func (s *ConnectionStorage) Read() ([]byte, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.storage.Load()
}

// Write
func (s *ConnectionStorage) Write(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.storage.Store(data)
}
