package config

import (
	"time"
)

type DatabaseConfig struct {
	Connection ConnectionConfig
	Driver     string
	Pool       PoolConfig
}

type ConnectionConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	Sslmode  string
}

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	host, err := getEnvValue("DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getEnvValue("DB_PORT")
	if err != nil {
		return nil, err
	}

	user, err := getEnvValue("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnvValue("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnvValue("DB_NAME")
	if err != nil {
		return nil, err
	}

	sslmode, err := getEnvValue("DB_SSLMODE")
	if err != nil {
		return nil, err
	}

	connection := ConnectionConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
		Sslmode:  sslmode,
	}

	driver, err := getEnvValue("DB_DRIVER")
	if err != nil {
		return nil, err
	}

	maxOpenConns, err := getIntEnvValue("DB_MAX_OPEN_CONNS")
	if err != nil {
		return nil, err
	}

	maxIdleConns, err := getIntEnvValue("DB_MAX_IDLE_CONNS")
	if err != nil {
		return nil, err
	}

	connMaxLifetime, err := getDurationEnvValue("DB_CONN_MAX_LIFETIME")
	if err != nil {
		return nil, err
	}

	connMaxIdleTime, err := getDurationEnvValue("DB_CONN_MAX_IDLE_TIME")
	if err != nil {
		return nil, err
	}

	pool := PoolConfig{
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
	}

	return &DatabaseConfig{
		Connection: connection,
		Driver:     driver,
		Pool:       pool,
	}, nil
}
