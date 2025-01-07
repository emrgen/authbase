package token

import (
	"container/heap"
	v1 "github.com/emrgent/authbase/api/v1"
}
type Token struct {
	token *v1.Tokens
	index int
}

type TokenQueue []*Token

func (tq TokenQueue) Len() int { return len(tq) }

func (tq TokenQueue) Less(i, j int) bool {
	return tq[i].token.Expiry < tq[j].token.Expiry
}

func (tq TokenQueue) Swap(i, j int) {
	tq[i], tq[j] = tq[j], tq[i]
	tq[i].index = i
	tq[j].index = j
}

func (tq *TokenQueue) Push(x interface{}) {
	n := len(*tq)
	item := x.(*Token)
	item.index = n
	*tq = append(*tq, item)
}

func (tq *TokenQueue) Pop() interface{} {
	old := *tq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*tq = old[0 : n-1]
	return item
}

type Registry struct {
	map[string]*v1.Tokens
	queue TokenQueue
	done chan struct{}
}

func NewRegistry() *Registry {
	return &Registry{
		map[string]*v1.Tokens{},
		TokenQueue{},
		make(chan struct{}),
	}
}

func (r *Registry) Add(token *v1.Tokens) {
	r.map[token.Token] = token
	heap.Push(&r.queue, &Token{token, len(r.queue)})
}

func (r *Registry) Remove(token *v1.Tokens) {
	delete(r.map, token.Token)
}

func (r *Registry) Expire() {
	for r.queue.Len() > 0 {
		item := heap.Pop(&r.queue).(*Token)
		if item.token.Expiry > time.Now().Unix() {
			heap.Push(&r.queue, item)
			break
		}
		delete(r.map, item.token.Token)
	}
}

func (r *Registry) Get(token string) *v1.Tokens {
	if token, ok := r.map[token]; ok {
		return token
	}
	return nil
}

func (r *Registry) List() []*v1.Tokens {
	tokens := []*v1.Tokens{}
	for _, token := range r.map{
		tokens = append(tokens, token)
	}
	return tokens
}

func (r *Registry) Size() int {
	return len(r.
	map)
}

func (r *Registry) Reset() {
	r.map = map[string]*v1.Tokens{}
	r.queue = TokenQueue{}
}

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
			token := r.queue.Pop()
			if token == nil {
				break
			}

			if token.token.Expiry > time.Now().Unix() {
				heap.Push(&r.queue, token)
				break
			}

			delete(r.map, token.token.Token)
		}

		case <-r.done:
			return
	}
}

// TokenRefresh is an interface for token refresh
type TokenRefresh interface {
	Refresh(token string) (*v1.Tokens, error)
}