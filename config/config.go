package config

import (
	"github.com/devstik0407/shenron/cleaner"
	"github.com/devstik0407/shenron/store"
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
	Port    int
	Address string
	Cleaner *cleaner.Config
	Store   *store.Config
}

func Load() {
	App = &Config{
		Port:    GetInt("PORT", 8080),
		Address: GetString("ADDRESS", "localhost"),
		Cleaner: &cleaner.Config{
			CleanupInterval: GetInt("CLEANUP_INTERVAL_IN_SECONDS", 2),
		},
		Store: &store.Config{
			DefaultExpiry: GetInt("DEFAULT_EXPIRY_IN_SECONDS", 1000),
		},
	}
}
