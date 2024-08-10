package config

import (
	"github.com/devstik0407/shenron/cleaner"
	"github.com/devstik0407/shenron/persistence"
	"github.com/devstik0407/shenron/service"
	"github.com/spf13/viper"
)

var App *Config

func init() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resources")
	viper.AddConfigPath("./../resources")

	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	Load()
}

type Config struct {
	Port          int
	Address       string
	Service       *service.Config
	CleanupWorker *cleaner.Config
	Persistence   *persistence.Config
}

func Load() {
	App = &Config{
		Port:    GetInt("PORT", 8080),
		Address: GetString("ADDRESS", "localhost"),
		Service: &service.Config{
			ReadTimeoutMs:          GetInt("READ_TIMEOUT_MS", 10),
			WriteTimeoutMs:         GetInt("WRITE_TIMEOUT_MS", 20),
			DefaultExpiryInSeconds: GetInt("DEFAULT_EXPIRY_SECONDS", 3600),
		},
		CleanupWorker: &cleaner.Config{
			CleanupInterval: GetInt("CLEANUP_INTERVAL_SECONDS", 5),
			WorkerCount:     GetInt("CLEANUP_WORKER_COUNT", 2),
		},
		Persistence: &persistence.Config{
			AOFWorkerCount:     GetInt("AOF_PERSISTENCE_WORKER_COUNT", 2),
			AOFFlushIntervalMs: GetInt("AOF_FLUSH_INTERVAL_MS", 1000),
		},
	}
}
