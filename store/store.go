package store

import (
	"github.com/devstik0407/shenron/pkg"
	"sync"
	"time"
)

type Config struct {
	DefaultExpiry int
}

type Store interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, duration time.Duration) error
	GetOldestExpiredKey() (string, error)
	Delete(key string)
}

type store struct {
	mu      sync.Mutex
	items   map[string]Item
	itemsPQ pkg.PriorityQueue
	config  *Config
}

type Item struct {
	Key    string
	Value  interface{}
	Expiry time.Time
}

func New(cfg *Config) *store {
	return &store{
		mu:      sync.Mutex{},
		items:   map[string]Item{},
		itemsPQ: pkg.PriorityQueue{},
		config:  cfg,
	}
}

func (s *store) Get(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return nil, ErrNotFound
	}

	if item.Expiry.Before(time.Now()) {
		return nil, ErrExpired
	}

	return item.Value, nil
}

func (s *store) Set(key string, value interface{}, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiryTime := time.Now().Add(duration)

	s.items[key] = Item{
		Key:    key,
		Value:  value,
		Expiry: expiryTime,
	}
	s.itemsPQ.PushItem(key, int(expiryTime.Unix()))

	return nil
}

func (s *store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
	s.itemsPQ.PopItem()
}

func (s *store) GetOldestExpiredKey() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	top := s.itemsPQ.TopItem()
	if top == nil || top.Priority() > int(time.Now().Unix()) {
		return "", ErrNoKeys
	}

	return top.Value(), nil
}
