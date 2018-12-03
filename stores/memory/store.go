package memory

import (
	"context"
	"sync"

	"github.com/nownabe/got/store"
)

type Provider struct {
	llock sync.Mutex
	locks map[string]*sync.RWMutex
	store map[string]interface{}
}

type Store struct {
	*Provider
	namespace string
}

func New() store.Provider {
	return &Provider{
		llock: sync.Mutex{},
		locks: map[string]*sync.RWMutex{},
		store: map[string]interface{}{},
	}
}

func (sp *Provider) Store(namespace string) store.Store {
	return &Store{
		Provider:  sp,
		namespace: namespace,
	}
}

func (s *Store) Get(ctx context.Context, key string) (interface{}, error) {
	s.lock().RLock()
	defer s.lock().RUnlock()

	val, ok := s.store[s.handlerKey(key)]
	if !ok {
		return nil, store.ErrStoreNotFound
	}

	return val, nil
}

func (s *Store) GetString(ctx context.Context, key string) (string, error) {
	v, err := s.Get(ctx, key)
	if err != nil {
		return "", err
	}

	str, ok := v.(string)
	if !ok {
		return "", store.ErrStoreNotStringValue
	}

	return str, nil
}

func (s *Store) Set(ctx context.Context, key string, val interface{}) error {
	s.lock().RLock()
	defer s.lock().RUnlock()

	s.store[s.handlerKey(key)] = val

	return nil
}

func (s *Store) lock() *sync.RWMutex {
	s.llock.Lock()
	defer s.llock.Unlock()
	if _, ok := s.locks[s.namespace]; !ok {
		s.locks[s.namespace] = &sync.RWMutex{}
	}
	return s.locks[s.namespace]
}

func (s *Store) handlerKey(key string) string {
	return s.namespace + ":" + key
}
