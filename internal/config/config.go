// Package config need for convenient get params.
package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	// ErrHandler need for calling if its happens error.
	ErrHandler = func(err error, msg string) { panic(fmt.Errorf("%s:%w", msg, err)) }
)

// Env returns value from os.Env.
func Env(key string) string {
	return os.Getenv(key)
}

// IntEnv returns inr from os.Env.
// If value not integer calling ErrHandler.
func IntEnv(key string) int {
	val := os.Getenv(key)

	i, err := strconv.Atoi(val)
	if err != nil {
		ErrHandler(err, fmt.Sprintf("nov valid integer: %d", i))
	}

	return i
}

// EnvOrDef returns value from os.Env.
// If value not set returns default value.
func EnvOrDef(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

// IntEnvOrDef returns value from os.Env.
// If value not set returns default value.
func IntEnvOrDef(key string, defaultValue int) int {
	val := os.Getenv(key)

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return i
}
