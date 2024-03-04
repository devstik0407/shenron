package cleaner

import (
	"fmt"
	"github.com/devstik0407/shenron/store"
	"time"
)

type Config struct {
	CleanupInterval int
}

func Clean(cfg *Config, s store.Store) {
	for true {
		key, err := s.GetOldestExpiredKey()
		if err != nil {
			continue
		}
		s.Delete(key)

		fmt.Printf("Deleted key: %s\n", key)

		time.Sleep(time.Duration(cfg.CleanupInterval) * time.Second)
	}
}
