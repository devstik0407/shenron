package persistence

import (
	"context"
	"fmt"
	"github.com/devstik0407/shenron/persistence/aof"
	"github.com/pkg/errors"
)

const (
	SetOperation    = "SET"
	DeleteOperation = "DELETE"
)

type Client interface {
	Apply(op *Operation) error
	StartWorkers(ctx context.Context)
}

type client struct {
	cfg *Config
	aof aof.Client
}

func New(cfg *Config) *client {
	aofClient := aof.New()
	return &client{
		cfg: cfg,
		aof: aofClient,
	}
}

func (c *client) Apply(op *Operation) error {
	if err := c.aof.Apply(toAOFOperation(op)); err != nil {
		return errors.Wrap(err, "[persistence.Apply] failed to persist operation in AOF")
	}
	return nil
}

func (c *client) StartWorkers(ctx context.Context) {
	for i := 0; i < c.cfg.AOFWorkerCount; i++ {
		go c.aof.Process(ctx, c.cfg.AOFFlushIntervalMs)
	}
	fmt.Printf("Started %d AOF workers\n", c.cfg.AOFWorkerCount)
}

//
//func (ps *PersistentStore) restorePreviousWrites() {
//	wd, _ := os.Getwd()
//	filePath := wd + "/" + appendOnlyFileName + "." + fileTypeText
//
//	aof, err := os.Open(filePath)
//	if err != nil {
//		fmt.Println("Error reading AOF file: ", err)
//		return
//	}
//	defer aof.Close()
//
//	scanner := bufio.NewScanner(aof)
//	for scanner.Scan() {
//		w := &Write{}
//		if err = json.Unmarshal([]byte(scanner.Text()), w); err != nil {
//			fmt.Println("Error reading AOF file: ", err)
//			return
//		}
//
//		expiresAt, err := time.Parse(time.RFC3339, w.ExpiresAt)
//		if err != nil {
//			fmt.Println("Failed to convert timestamp:", err)
//			return
//		}
//
//		if expiresAt.Before(time.Now()) {
//			continue
//		}
//
//		if _, err = ps.storage.Set(w.Key, w.Value, expiresAt.Sub(time.Now())); err != nil {
//			fmt.Printf("Failed to insert key %s in store: %s", w.Key, err)
//			return
//		}
//	}
//}
