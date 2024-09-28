package env

import (
	"os"
	"strconv"
	"time"
)

// GetBoolOrDefault will return the value of the environment variable name provided if it exists.
// If the environment variable can't be found, or an error occurs, returns defaultV.
func GetBoolOrDefault(name string, defaultV bool) bool {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsBool, err := strconv.ParseBool(v)
	if err != nil {
		return defaultV
	}
	return vAsBool
}

// GetDurationOrDefault will return the value of the environment variable name provided if it exists.
// If the environment variable can't be found, or an error occurs, returns defaultV.
func GetDurationOrDefault(name string, defaultV time.Duration) time.Duration {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsDuration, err := time.ParseDuration(v)
	if err != nil {
		return defaultV
	}
	return vAsDuration
}

// GetIntOrDefault will return the value of the environment variable name provided if it exists.
// If the environment variable can't be found, or an error occurs, returns defaultV.
func GetIntOrDefault(name string, defaultV int) int {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	vAsInt, err := strconv.Atoi(v)
	if err != nil {
		return defaultV
	}
	return vAsInt
}

// GetStringOrDefault will return the value of the environment variable name provided if it exists.
// If the environment variable can't be found, returns defaultV.
func GetStringOrDefault(name, defaultV string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultV
	}
	return v
}
