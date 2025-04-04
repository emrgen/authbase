package token

import (
	"container/heap"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/sirupsen/logrus"
	"time"
)

// Token is a wrapper for a JWT token
type Token struct {
	id    string
	token *v1.Tokens
	index int
}

// NewToken creates a new JWT token wrapper with an ID and token
func NewToken(id string, token *v1.Tokens) *Token {
	return &Token{
		id:    id,
		token: token,
	}
}

type TokenQueue []*Token

// Len returns the length of the token queue
func (tq TokenQueue) Len() int { return len(tq) }

// Less returns true if the token at index i expires before the token at index j
func (tq TokenQueue) Less(i, j int) bool {
	return tq[i].token.ExpiresAt.AsTime().Before(tq[j].token.ExpiresAt.AsTime())
}

// Swap swaps the tokens at index i and j
func (tq TokenQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
	tq[i].index = i
	tq[j].index = j
}

// Push adds a token to the queue
func (tq *TokenQueue) Push(x interface{}) {
	n := len(*tq)
	item := x.(*Token)
	item.index = n
	*tq = append(*tq, item)
}

// Pop removes a token from the queue
func (tq *TokenQueue) Pop() interface{} {
	old := *tq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*tq = old[0 : n-1]
	return item
}

// Registry is a token registry
type Registry struct {
	keys   map[string]string
	tokens map[string]*v1.Tokens
	queue  TokenQueue
	done   chan struct{}
	rotate Rotate
}

// NewRegistry creates a new token registry
func NewRegistry() *Registry {
	return &Registry{
		keys:   make(map[string]string),
		tokens: make(map[string]*v1.Tokens),
		queue:  TokenQueue{},
		done:   make(chan struct{}),
	}
}

func (r *Registry) Add(id, key string) error {
	r.keys[id] = key

	return nil
}

func (r *Registry) Remove(key string) {
	delete(r.tokens, key)
	delete(r.keys, key)
}

func (r *Registry) Expire() {
	for r.queue.Len() > 0 {
		item := heap.Pop(&r.queue).(*Token)
		if item.token.ExpiresAt.AsTime().After(time.Now()) {
			heap.Push(&r.queue, item)
			break
		}
		delete(r.tokens, item.id)
	}
}

func (r *Registry) Get(token string) *v1.Tokens {
	return r.tokens[token]
}

func (r *Registry) Size() int {
	return len(r.tokens)
}

func (r *Registry) Reset() {
	r.tokens = make(map[string]*v1.Tokens)
	r.queue = TokenQueue{}
	r.keys = make(map[string]string)
	r.done = make(chan struct{})
}

func (r *Registry) getTokens() {
	for id, key := range r.keys {
		token, err := r.rotate.GetToken(key)
		if err != nil {
			logrus.Errorf("failed to get token: %v", err.Error())
			continue
		}
		r.tokens[id] = token
		r.queue.Push(NewToken(id, token))
	}
}

// Start starts the token rotation
func (r *Registry) Start() {
	go r.Run()
}

func (r *Registry) Stop() {
	close(r.done)
}

func (r *Registry) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			item := r.queue.Pop()
			if item == nil {
				break
			}

			token := item.(*Token)
			if token.token.ExpiresAt.AsTime().After(time.Now()) {
				heap.Push(&r.queue, token)
			} else {
				delete(r.tokens, token.id)
				newToken, err := r.rotate.RefreshToken(token.token.RefreshToken)
				if err != nil {
					logrus.Infof("failed to refresh token")
					continue
				}
				r.tokens[token.id] = newToken
			}
		case <-r.done:
			return
		}
	}
}

// Rotate is an interface for token refresh
type Rotate interface {
	GetToken(key string) (*v1.Tokens, error)
	RefreshToken(token string) (*v1.Tokens, error)
}
