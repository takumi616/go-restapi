package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func getEnvValue(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s must be set", key)
	}
	return val, nil
}

func getDurationEnvValue(key string) (time.Duration, error) {
	strVal, err := getEnvValue(key)
	if err != nil {
		return 0, err
	}

	val, err := time.ParseDuration(strVal)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: '%s': %w", strVal, err)
	}
	return val, nil
}

func getIntEnvValue(key string) (int, error) {
	strVal, err := getEnvValue(key)
	if err != nil {
		return 0, err
	}

	val, err := strconv.Atoi(strVal)
	if err != nil {
		return 0, fmt.Errorf("invalid int format: '%s': %w", strVal, err)
	}
	return val, nil
}
