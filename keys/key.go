package keys

import (
	"context"
	"crypto/rsa"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/x"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type PublicKey struct {
	key      *rsa.PrivateKey
	ExpireAt time.Time
}

// PublicRegistry struct to store public key and manage it
// When the key is expired, it will be removed and get new key from the server.
type PublicRegistry struct {
	keys   map[string]*PublicKey
	client v1.PublicKeyServiceClient
	mu     sync.Mutex
}

func NewPublicRegistry(client v1.PublicKeyServiceClient) *PublicRegistry {
	return &PublicRegistry{
		client: client,
		keys:   make(map[string]*PublicKey),
		mu:     sync.Mutex{},
	}
}

func (r *PublicRegistry) GetKey(id string) (*PublicKey, error) {
	if key, ok := r.keys[id]; ok {
		return key, nil
	}

	res, err := r.client.GetPublicKey(context.TODO(), &v1.GetPublicKeyRequest{Id: id})
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(res.Public.Key))
	if err != nil {
		return nil, err
	}
	r.keys[id] = &PublicKey{
		key:      key,
		ExpireAt: res.ExpireAt.AsTime(),
	}

	return r.keys[id], nil
}

func (r *PublicRegistry) AddKey(id string, key *rsa.PrivateKey) {
	r.keys[id] = &PublicKey{
		key:      key,
		ExpireAt: time.Now().Add(time.Hour),
	}
}

func (r *PublicRegistry) RemoveKey(id string) {
	delete(r.keys, id)
}

func (r *PublicRegistry) Reset() {
	r.keys = make(map[string]*PublicKey)
}

func (r *PublicRegistry) Size() int {
	return len(r.keys)
}

func (r *PublicRegistry) Run() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			r.mu.Lock()
			for id, key := range r.keys {
				r.refresh(id, key)
			}
		}
	}
}

func (r *PublicRegistry) refresh(id string, public *PublicKey) {
	if public.ExpireAt.After(time.Now()) {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.keys, id)

	res, err := r.client.GetPublicKey(context.TODO(), &v1.GetPublicKeyRequest{Id: id})
	if err != nil {
		logrus.Errorf("Failed to get public key: %v", err)
		// schedule exponential backoff
		return
	}

	publicKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(res.Public.Key))
	if err != nil {
		logrus.Errorf("Failed to parse public key: %v", err)
		// schedule exponential backoff
		return
	}

	r.keys[id] = &PublicKey{
		key:      publicKey,
		ExpireAt: res.ExpireAt.AsTime(),
	}
}

type KeyPair struct {
	private  *rsa.PrivateKey
	public   *rsa.PublicKey
	ExpireAt time.Time
}

// PrivateRegistry struct to store private key pair and manage it
// When the key pair is expired, it will be removed and generate new key pair
type PrivateRegistry struct {
	keys map[string]*KeyPair
	mu   sync.Mutex
}

func NewPrivateRegistry() *PrivateRegistry {
	return &PrivateRegistry{
		keys: make(map[string]*KeyPair),
		mu:   sync.Mutex{},
	}
}

// GetKey function to get key pair
func (r *PrivateRegistry) GetKey(id string) (*KeyPair, error) {
	key, ok := r.keys[id]
	if ok {
		return key, nil
	}

	err := r.GenerateKeyPair(id)
	if err != nil {
		return nil, err
	}

	return r.keys[id], nil
}

func (r *PrivateRegistry) AddKey(id string, key *KeyPair) {
	r.keys[id] = key
}

func (r *PrivateRegistry) RemoveKey(id string) {
	delete(r.keys, id)
}

func (r *PrivateRegistry) Reset() {
	r.keys = make(map[string]*KeyPair)
}

func (r *PrivateRegistry) Size() int {
	return len(r.keys)
}

func (r *PrivateRegistry) GenerateKeyPair(id string) error {
	private, public, err := x.GenerateKeyPair(2048)
	if err != nil {
		return err
	}

	r.AddKey(id, &KeyPair{
		private: private,
		public:  public,
	})

	return nil
}

func (r *PrivateRegistry) Run() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			for id, key := range r.keys {
				if key.ExpireAt.Before(time.Now()) {
					delete(r.keys, id)
					err := r.GenerateKeyPair(id)
					if err != nil {
						continue
					}
				}
			}
		}
	}
}

func (r *PrivateRegistry) refresh(id string, key *KeyPair) {
	if key.ExpireAt.After(time.Now()) {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.keys, id)
	err := r.GenerateKeyPair(id)
	if err != nil {
		logrus.Errorf("Failed to generate key pair: %v", err)
		// schedule exponential backoff
	}
}
