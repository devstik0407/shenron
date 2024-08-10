package store

import (
	"context"
	"github.com/devstik0407/shenron/pkg"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var storeInstance *store

type Store interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) (*Item, error)
	Delete(key string)
	DeleteOldestExpiredKey() string
}

type store struct {
	mu      *sync.RWMutex
	items   map[string]Item
	itemsPQ pkg.PriorityQueue
}

func GetInstance() Store {
	if storeInstance == nil {
		storeInstance = &store{
			mu:      &sync.RWMutex{},
			items:   map[string]Item{},
			itemsPQ: pkg.PriorityQueue{},
		}
	}
	return storeInstance
}

func (s *store) Get(ctx context.Context, key string) (interface{}, error) {
	defer s.releaseReadLock()

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "[store.Get] read timeout")
	case <-s.acquireReadLock():
		item, ok := s.items[key]
		if !ok {
			return nil, ErrNotFound
		}
		if item.ExpiresAt.Before(time.Now()) {
			return nil, ErrExpired
		}
		return item.Value, nil
	}

}

func (s *store) Set(ctx context.Context, key string, value interface{}, duration time.Duration) (*Item, error) {
	defer s.releaseWriteLock()

	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "[store.Set] write timeout")
	case <-s.acquireWriteLock():
		expiryTime := time.Now().Add(duration)
		s.items[key] = Item{
			Key:       key,
			Value:     value,
			ExpiresAt: expiryTime,
		}
		s.itemsPQ.PushItem(key, int(expiryTime.Unix()))
		return &Item{
			Key:       key,
			Value:     value,
			ExpiresAt: expiryTime,
		}, nil
	}
}

func (s *store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
	s.itemsPQ.PopItem()
}

func (s *store) DeleteOldestExpiredKey() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	top := s.itemsPQ.TopItem()
	if top == nil {
		return ""
	}

	key := top.Value()
	if s.items[key].ExpiresAt.Before(time.Now()) {
		delete(s.items, key)
		s.itemsPQ.PopItem()
		return key
	}

	return ""
}

func (s *store) acquireReadLock() chan int {
	s.mu.RLock()
	ch := make(chan int, 1)
	ch <- 1
	return ch
}

func (s *store) acquireWriteLock() chan int {
	s.mu.Lock()
	ch := make(chan int, 1)
	ch <- 1
	return ch
}

func (s *store) releaseReadLock() {
	s.mu.RUnlock()
}

func (s *store) releaseWriteLock() {
	s.mu.Unlock()
}
