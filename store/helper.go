package store

import "time"

type Item struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
}
