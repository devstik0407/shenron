package aof

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Client interface {
	Apply(operation Operation) error
	Process(ctx context.Context, flushIntervalMs int)
}

type client struct {
	buffer chan Operation
}

func New() *client {
	return &client{
		buffer: make(chan Operation, 1000),
	}
}

func (c *client) Apply(operation Operation) error {
	if operation != nil {
		switch operation.(type) {
		case *Set:
			c.buffer <- operation
		case *Delete:
			c.buffer <- operation
		default:
			return ErrUndefinedOperation
		}
	}
	return nil
}

func (c *client) Process(ctx context.Context, flushIntervalMs int) {
	ticker := time.NewTicker(time.Duration(flushIntervalMs) * time.Millisecond)
	var ops []Operation
	for {
		select {
		case <-ctx.Done():
			return
		case op := <-c.buffer:
			ops = append(ops, op)
			if len(ops) == 10 {
				c.appendOperation(ops)
				ops = nil
			}
			break
		case <-ticker.C:
			c.appendOperation(ops)
			ops = nil
			break
		}
	}
}
func (c *client) appendOperation(ops []Operation) {
	file, err := os.OpenFile("./resources/append-only-file.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, op := range ops {
		data, err := json.Marshal(op)
		if err != nil {
			continue
		}
		if _, err = file.WriteString(fmt.Sprintf("Timestamp:%s Operation:%s Data:%s\n", op.ExecAt(), op.Name(), data)); err != nil {
			panic(err)
		}
	}
}
