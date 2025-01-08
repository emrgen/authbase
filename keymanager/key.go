package keymanager

import (
	"context"
	"crypto/rsa"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type staticVerifier struct {
	key []byte
}

func newStaticVerifier(key []byte) *staticVerifier {
	return &staticVerifier{key: key}
}

func (v *staticVerifier) Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return v.key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to get claims")
	}

	return claims, nil
}

type staticSigner struct {
	key []byte
}

func newStaticSigner(key []byte) *staticSigner {
	return &staticSigner{key: []byte(key)}
}

func (s *staticSigner) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type StaticKeyProvider struct {
	key []byte
}

func NewStaticKeyProvider(key string) *StaticKeyProvider {
	return &StaticKeyProvider{key: []byte(key)}
}

func (r *StaticKeyProvider) GetSigner(id string) (x.JWTSigner, error) {
	return newStaticSigner(r.key), nil
}

func (r *StaticKeyProvider) GetVerifier(id string) (x.JWTVerifier, error) {
	return newStaticVerifier(r.key), nil
}

type PublicKey struct {
	key      *rsa.PrivateKey
	ExpireAt time.Time
}

// PublicRegistry struct to store public key and manage it
// When the key is expired, it will be removed and get new key from the server.
type PublicRegistry struct {
	keys   map[string]*PublicKey
	mu     sync.Mutex
	client v1.PublicKeyServiceClient
	slack  time.Duration
}

func NewPublicRegistry(client v1.PublicKeyServiceClient) *PublicRegistry {
	return &PublicRegistry{
		client: client,
		keys:   make(map[string]*PublicKey),
		mu:     sync.Mutex{},
		slack:  time.Minute * 10,
	}
}

func (r *PublicRegistry) GetSignKey(id string) (*PublicKey, error) {
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
			for id, key := range r.keys {
				r.refresh(id, key)
			}
		}
	}
}

// refresh checks if the token has expired
func (r *PublicRegistry) refresh(id string, public *PublicKey) {
	if public.ExpireAt.Add(r.slack).After(time.Now()) {
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
	keys  map[string]*KeyPair
	mu    sync.Mutex
	store store.Provider // save key pair to the store
}

func NewPrivateRegistry(store store.Provider) *PrivateRegistry {
	return &PrivateRegistry{
		keys:  make(map[string]*KeyPair),
		mu:    sync.Mutex{},
		store: store,
	}
}

func (r *PrivateRegistry) GetGetSignKey(id string) (*rsa.PrivateKey, error) {
	key, err := r.GetKey(id)
	if err != nil {
		return nil, errors.New("failed to get private key")
	}

	return key.private, nil
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
