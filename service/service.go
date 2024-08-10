package service

import (
	"context"
	"fmt"
	"github.com/devstik0407/shenron/cleaner"
	"time"

	"github.com/devstik0407/shenron/persistence"
	"github.com/devstik0407/shenron/store"
	"github.com/pkg/errors"
)

type Service interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttlInSeconds int) error
	Delete(key string)
}

type service struct {
	cfg         *Config
	storage     store.Store
	persistence persistence.Client
	cleaner     cleaner.Worker
}

func New(cfg *Config, s store.Store, p persistence.Client, c cleaner.Worker) *service {
	return &service{
		cfg:         cfg,
		storage:     s,
		persistence: p,
		cleaner:     c,
	}
}

func (svc *service) Start(ctx context.Context) {
	svc.cleaner.Run(ctx)
	svc.persistence.StartWorkers(ctx)
}

func (svc *service) Get(ctx context.Context, key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.cfg.ReadTimeoutMs)*time.Millisecond)
	defer cancel()

	value, err := svc.storage.Get(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "[service.Get] failed to get value from store")
	}
	return value, nil
}

func (svc *service) Set(ctx context.Context, key string, value interface{}, ttlInSeconds int) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(svc.cfg.WriteTimeoutMs)*time.Millisecond)
	defer cancel()

	if ttlInSeconds == 0 {
		ttlInSeconds = svc.cfg.DefaultExpiryInSeconds
	}

	item, err := svc.storage.Set(ctx, key, value, time.Duration(ttlInSeconds)*time.Second)
	if err != nil || item == nil {
		return errors.Wrap(err, "[service.Set] failed to create new key-value pair")
	}

	if err = svc.persistence.Apply(&persistence.Operation{
		Type:         persistence.SetOperation,
		Key:          key,
		Value:        value,
		TTLInSeconds: ttlInSeconds,
		CreatedAt:    time.Now(),
	}); err != nil {
		svc.storage.Delete(key)
		return errors.Wrap(err, "[service.Set] failed to persist write operation")
	}

	return nil
}

func (svc *service) Delete(key string) {
	svc.storage.Delete(key)
	if err := svc.persistence.Apply(&persistence.Operation{
		Type: persistence.DeleteOperation,
		Key:  key,
	}); err != nil {
		fmt.Printf("[service.Delete] failed to persist delete operation: %s", err)
	}
	return
}
