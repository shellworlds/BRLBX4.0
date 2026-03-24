package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load enables reading from environment variables.
func Load() {
	viper.AutomaticEnv()
}

// MustGet returns a required string env or panics during startup.
func MustGet(key string) string {
	v := viper.GetString(key)
	if strings.TrimSpace(v) == "" {
		panic(fmt.Errorf("config: required %s is empty", key))
	}
	return v
}

// GetString returns viper string (may be empty).
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns viper int with default.
func GetInt(key string, def int) int {
	if !viper.IsSet(key) {
		return def
	}
	return viper.GetInt(key)
}

// GetBool returns viper bool with default.
func GetBool(key string, def bool) bool {
	if !viper.IsSet(key) {
		return def
	}
	return viper.GetBool(key)
}

// GetFloat64 returns viper float64 with default.
func GetFloat64(key string, def float64) float64 {
	if !viper.IsSet(key) {
		return def
	}
	return viper.GetFloat64(key)
}
