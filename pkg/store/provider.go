package store

import (
	"errors"
	"github.com/google/uuid"
)

type AuthBaseStoreProvider interface {
	Provide(orgID uuid.UUID) (AuthBaseStore, error)
}

type Provider struct {
	Stores map[uuid.UUID]AuthBaseStore
}

func (s *Provider) Provide(orgID uuid.UUID) (AuthBaseStore, error) {
	store, ok := s.Stores[orgID]
	if !ok {
		return nil, errors.New("store not found")
	}

	return store, nil
}

func NewProvider() *Provider {
	return &Provider{
		Stores: make(map[uuid.UUID]AuthBaseStore),
	}
}

type DefaultProvider struct {
	Store AuthBaseStore
}

func NewDefaultProvider(store AuthBaseStore) *DefaultProvider {
	return &DefaultProvider{
		Store: store,
	}
}

func (d *DefaultProvider) Provide(orgID uuid.UUID) (AuthBaseStore, error) {
	return d.Store, nil
}
