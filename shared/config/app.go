package config

import (
	"time"
)

type AppConfig struct {
	Port    string
	Timeout TimeoutConfig
}

type TimeoutConfig struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func NewAppConfig() (*AppConfig, error) {
	port, err := getEnvValue("APP_PORT")
	if err != nil {
		return nil, err
	}

	readTimeout, err := getDurationEnvValue("SERVER_READ_TIMEOUT")
	if err != nil {
		return nil, err
	}
	readHeaderTimeout, err := getDurationEnvValue("SERVER_READ_HEADER_TIMEOUT")
	if err != nil {
		return nil, err
	}
	writeTimeout, err := getDurationEnvValue("SERVER_WRITE_TIMEOUT")
	if err != nil {
		return nil, err
	}
	idleTimeout, err := getDurationEnvValue("SERVER_IDLE_TIMEOUT")
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Port: port,
		Timeout: TimeoutConfig{
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: readHeaderTimeout,
			WriteTimeout:      writeTimeout,
			IdleTimeout:       idleTimeout,
		},
	}, nil
}
