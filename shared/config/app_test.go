package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/takumi616/go-restapi/shared/config"
)

var keyList = []string{
	"APP_PORT", "SERVER_READ_TIMEOUT", "SERVER_READ_HEADER_TIMEOUT",
	"SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
}

type expected struct {
	port              string
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
}

func TestNewAppConfigNormal(t *testing.T) {
	inputList := []string{"8080", "2s", "3s", "4s", "5s"}
	expected := expected{
		port: "8080", readTimeout: 2 * time.Second, readHeaderTimeout: 3 * time.Second,
		writeTimeout: 4 * time.Second, idleTimeout: 5 * time.Second,
	}

	for i, key := range keyList {
		t.Setenv(key, inputList[i])
	}

	appCfg, err := config.NewAppConfig()

	assert.NoError(t, err)
	assert.NotNil(t, appCfg)
	assert.Equal(t, expected.port, appCfg.Port)
	assert.Equal(t, expected.readTimeout, appCfg.Timeout.ReadTimeout)
	assert.Equal(t, expected.readHeaderTimeout, appCfg.Timeout.ReadHeaderTimeout)
	assert.Equal(t, expected.writeTimeout, appCfg.Timeout.WriteTimeout)
	assert.Equal(t, expected.idleTimeout, appCfg.Timeout.IdleTimeout)
}

func TestNewAppConfigEmptyPort(t *testing.T) {
	portKey := "APP_PORT"
	inputList := []string{"", "2s", "3s", "4s", "5s"}

	for i, key := range keyList {
		t.Setenv(key, inputList[i])
	}

	appCfg, err := config.NewAppConfig()

	assert.Nil(t, appCfg)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("environment variable %s must be set", portKey))
}

func TestNewAppConfigInvalidDuration(t *testing.T) {
	invalidDuration := "s2"
	inputList := []string{"8080", invalidDuration, "3s", "4s", "5s"}

	for i, key := range keyList {
		t.Setenv(key, inputList[i])
	}

	appCfg, err := config.NewAppConfig()

	assert.Nil(t, appCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("invalid duration format: '%s'", invalidDuration))
}
