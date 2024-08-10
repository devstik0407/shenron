package cleaner

import (
	"context"
	"fmt"
	"github.com/devstik0407/shenron/store"
	"time"
)

type Worker interface {
	Run(ctx context.Context)
}

type worker struct {
	cfg     *Config
	storage store.Store
}

func New(cfg *Config, s store.Store) *worker {
	return &worker{
		cfg:     cfg,
		storage: s,
	}
}

func (w *worker) Run(ctx context.Context) {
	for i := 0; i < w.cfg.WorkerCount; i++ {
		go w.process(ctx)
	}
	fmt.Printf("Started %d cleanup workers\n", w.cfg.WorkerCount)
}

func (w *worker) process(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(w.cfg.CleanupInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[cleaner.process] Stopping cleaner process...")
			return
		case <-ticker.C:
			deletedKey := w.storage.DeleteOldestExpiredKey()
			fmt.Printf("[cleaner.process] Deleted key: %s\n", deletedKey)
		}
	}
}
