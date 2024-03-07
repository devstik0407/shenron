package config

import "github.com/spf13/viper"

func GetString(key string, defaultValue string) string {
	if !viper.IsSet(key) {
		return defaultValue
	}

	return viper.GetString(key)
}

func GetInt(key string, defaultValue int) int {
	if !viper.IsSet(key) {
		return defaultValue
	}

	return viper.GetInt(key)
}
