package store

import (
	"context"
	"github.com/emrgen/authbase/pkg/config"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// provider.go contains the logic to provide the correct store based on the project ID.
// It is used when the project has its own store. If the project does not have its own store,
// The default store is used when the project ID is not provided in the context.

// GetProjectStore returns the store for the given project ID. If the project ID is not provided, the default store is returned.
func GetProjectStore(ctx context.Context, store Provider) (AuthBaseStore, error) {
	if config.GetConfig().Mode == config.ModeSingleStore {
		return store.Default(), nil
	}

	// get project id from the context
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	projectID := md.Get("project_id")
	if len(projectID) == 0 {
		return store.Default(), nil
	}

	projectUUID, err := uuid.Parse(projectID[0])
	if err != nil {
		return nil, err
	}

	projectStore, err := store.Provide(projectUUID)
	if err != nil {
		return nil, err
	}

	return projectStore, nil
}

// Provider is an interface to provide the store based on the project ID.
type Provider interface {
	Provide(projectID uuid.UUID) (AuthBaseStore, error)
	Default() AuthBaseStore
}

// MultiStoreProvider is a provider with a cache.
// It caches the store for the given project ID.
type MultiStoreProvider struct {
	Stores   map[uuid.UUID]AuthBaseStore
	provider Provider
	mu       sync.RWMutex // mu guards the stores map
}

// NewProvider creates a new provider instance.
func NewProvider(provider Provider) *MultiStoreProvider {
	return &MultiStoreProvider{
		Stores:   make(map[uuid.UUID]AuthBaseStore),
		provider: provider,
	}
}

// Provide returns a store for the given project ID.
func (s *MultiStoreProvider) Provide(projectID uuid.UUID) (AuthBaseStore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	store, ok := s.Stores[projectID]
	if ok {
		return store, nil
	}

	store, err := s.provider.Provide(projectID)
	if err != nil {
		return nil, err
	}

	// cache the store for the given project ID
	s.Stores[projectID] = store

	return store, nil
}

// Default returns the default store.
func (s *MultiStoreProvider) Default() AuthBaseStore {
	return s.provider.Default()
}

// DefaultProvider is a default provider implementation.
type DefaultProvider struct {
	Store AuthBaseStore
}

// NewDefaultProvider creates a new default provider instance.
func NewDefaultProvider(store AuthBaseStore) *DefaultProvider {
	return &DefaultProvider{
		Store: store,
	}
}

// Provide returns the default store.
func (d *DefaultProvider) Provide(projectID uuid.UUID) (AuthBaseStore, error) {
	return d.Store, nil
}

// Default returns the default store.
func (d *DefaultProvider) Default() AuthBaseStore {
	return d.Store
}
