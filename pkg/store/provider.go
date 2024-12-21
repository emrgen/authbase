package store

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"sync"
)

// provider.go contains the logic to provide the correct store based on the project ID.
// It is used when the project has its own store. If the project does not have its own store,
// The default store is used when the project ID is not provided in the context.

func GetProjectStore(ctx context.Context, store Provider) (AuthBaseStore, error) {
	// get project id from the context
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	projectID := md.Get("project_id")
	if len(projectID) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing project id")
	}

	projectIDUUID, err := uuid.Parse(projectID[0])
	if err != nil {
		return nil, err
	}

	projectStore, err := store.Provide(projectIDUUID)
	if err != nil {
		// if the project store is not found, return the default store
		if os.Getenv("APP_MODE") != "masterless" {
			return store.Default(), nil
		}

		return nil, err
	}

	return projectStore, nil
}

type Provider interface {
	Provide(projectID uuid.UUID) (AuthBaseStore, error)
	Default() AuthBaseStore
}

type ProviderCache struct {
	Stores   map[uuid.UUID]AuthBaseStore
	provider Provider
	mu       sync.RWMutex // mu guards the stores map
}

// NewProvider creates a new provider instance.
func NewProvider(provider Provider) *ProviderCache {
	return &ProviderCache{
		Stores:   make(map[uuid.UUID]AuthBaseStore),
		provider: provider,
	}
}

// Provide returns a store for the given organization ID.
func (s *ProviderCache) Provide(projectID uuid.UUID) (AuthBaseStore, error) {
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

	s.Stores[projectID] = store

	return store, nil
}

// Default returns the default store.
func (s *ProviderCache) Default() AuthBaseStore {
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

func (d *DefaultProvider) Provide(projectID uuid.UUID) (AuthBaseStore, error) {
	return d.Store, nil
}

func (d *DefaultProvider) Default() AuthBaseStore {
	return d.Store
}
