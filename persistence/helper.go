package persistence

import (
	"github.com/devstik0407/shenron/persistence/aof"
	"time"
)

type Config struct {
	AOFWorkerCount     int
	AOFFlushIntervalMs int
}

type Operation struct {
	Type         string      `json:"type"`
	Key          string      `json:"key"`
	Value        interface{} `json:"value"`
	TTLInSeconds int         `json:"ttl_in_seconds"`
	CreatedAt    time.Time   `json:"created_at"`
}

func toAOFOperation(op *Operation) aof.Operation {
	if op == nil {
		return nil
	}
	if op.Type == SetOperation {
		return &aof.Set{
			Key:          op.Key,
			Value:        op.Value,
			CreatedAt:    op.CreatedAt,
			TTLInSeconds: op.TTLInSeconds,
			Metadata:     &aof.Metadata{ExecutedAt: time.Now()},
		}
	}
	if op.Type == DeleteOperation {
		return &aof.Delete{
			Key:      op.Key,
			Metadata: &aof.Metadata{ExecutedAt: time.Now()},
		}
	}
	return nil
}
