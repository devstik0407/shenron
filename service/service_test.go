package service

import (
	"context"
	"fmt"
	"github.com/devstik0407/shenron/cleaner"
	"github.com/devstik0407/shenron/persistence"
	"github.com/devstik0407/shenron/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_service(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := store.GetInstance()
	p := persistence.New(&persistence.Config{
		AOFWorkerCount:     1,
		AOFFlushIntervalMs: 1000,
	})
	c := cleaner.New(&cleaner.Config{
		CleanupInterval: 5,
		WorkerCount:     2,
	}, s)

	svc := New(&Config{
		ReadTimeoutMs:          1,
		WriteTimeoutMs:         2,
		DefaultExpiryInSeconds: 60,
	}, s, p, c)
	svc.Start(ctx)

	for i := 0; i < 1000; i++ {
		err := svc.Set(ctx, fmt.Sprintf("key#%d", i), fmt.Sprintf("value#%d", i), 0)
		assert.NoError(t, err)
	}

	value, err := svc.Get(ctx, "key#1")
	assert.NoError(t, err)
	assert.Equal(t, "value#1", value)
	time.Sleep(time.Second * 10)
}
