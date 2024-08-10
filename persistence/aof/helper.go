package aof

import "time"

type Operation interface {
	Name() string
	ExecAt() string
}

type Set struct {
	Key          string      `json:"key"`
	Value        interface{} `json:"value"`
	CreatedAt    time.Time   `json:"created_at"`
	TTLInSeconds int         `json:"ttl_in_seconds"`
	Metadata     *Metadata   `json:"metadata"`
}

type Delete struct {
	Key      string    `json:"key"`
	Metadata *Metadata `json:"metadata"`
}

type Metadata struct {
	ExecutedAt time.Time `json:"executed_at"`
}

func (s *Set) Name() string {
	return "SET"
}

func (s *Set) ExecAt() string {
	return s.Metadata.ExecutedAt.Format(time.RFC3339)
}

func (d *Delete) Name() string {
	return "DELETE"
}

func (s *Delete) ExecAt() string {
	return s.Metadata.ExecutedAt.Format(time.RFC3339)
}
