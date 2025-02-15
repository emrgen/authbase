package secret

import (
	"github.com/emrgen/authbase/pkg/model"
	"gorm.io/gorm"
)

// Store is an interface for getting secrets
type Store interface {
	// GetSecret returns the secret for the given key
	GetSecret(key string) (string, error)
	// SetSecret sets the secret for the given key
	SetSecret(key string, value string) error
}

// MemStore is an in-memory implementation of Store
type MemStore struct {
	Secrets map[string]string
}

// NewMemStore creates a new MemStore
func NewMemStore() *MemStore {
	return &MemStore{
		Secrets: make(map[string]string),
	}
}

func (s *MemStore) GetSecret(key string) (string, error) {
	secret, ok := s.Secrets[key]
	if !ok {
		return "", nil
	}

	return secret, nil
}

func (s *MemStore) SetSecret(key string, value string) error {
	s.Secrets[key] = value
	return nil
}

type GormStore struct {
	DB *gorm.DB
}

// NewGormStore creates a new GormStore
func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{
		DB: db,
	}
}

func (s *GormStore) GetSecret(key string) (string, error) {
	var secret model.Secret
	err := s.DB.Where("key = ?", key).First(&secret).Error
	if err != nil {
		return "", err
	}

	return secret.Value, nil
}

func (s *GormStore) SetSecret(key string, value string) error {
	secret := model.Secret{
		ID:    key,
		Value: value,
	}

	return s.DB.Create(&secret).Error
}
