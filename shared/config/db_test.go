package config_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/takumi616/go-restapi/shared/config"
)

var dbEnvKeyList = []string{
	"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "DB_DRIVER",
	"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS", "DB_CONN_MAX_LIFETIME", "DB_CONN_MAX_IDLE_TIME",
}

type expectedDatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	Sslmode         string
	Driver          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func TestNewDatabaseConfigNormal(t *testing.T) {
	inputList := []string{
		"host", "5432", "user", "pass", "dbname", "disable", "postgres",
		"25", "10", "5m", "2m",
	}
	expected := expectedDatabaseConfig{
		"host", "5432", "user", "pass", "dbname", "disable", "postgres",
		25, 10, 5 * time.Minute, 2 * time.Minute,
	}

	for i, key := range dbEnvKeyList {
		t.Setenv(key, inputList[i])
	}

	dbCfg, err := config.NewDatabaseConfig()

	assert.NoError(t, err)
	assert.NotNil(t, dbCfg)
	assert.Equal(t, expected.Host, dbCfg.Connection.Host)
	assert.Equal(t, expected.Port, dbCfg.Connection.Port)
	assert.Equal(t, expected.User, dbCfg.Connection.User)
	assert.Equal(t, expected.Password, dbCfg.Connection.Password)
	assert.Equal(t, expected.DbName, dbCfg.Connection.DbName)
	assert.Equal(t, expected.Sslmode, dbCfg.Connection.Sslmode)
	assert.Equal(t, expected.Driver, dbCfg.Driver)
	assert.Equal(t, expected.MaxOpenConns, dbCfg.Pool.MaxOpenConns)
	assert.Equal(t, expected.MaxIdleConns, dbCfg.Pool.MaxIdleConns)
	assert.Equal(t, expected.ConnMaxLifetime, dbCfg.Pool.ConnMaxLifetime)
	assert.Equal(t, expected.ConnMaxIdleTime, dbCfg.Pool.ConnMaxIdleTime)
}

func TestNewDatabaseConfigEmptyDriver(t *testing.T) {
	driverKey := "DB_DRIVER"
	inputList := []string{
		"host", "5432", "user", "pass", "dbname", "disable", "",
		"25", "10", "5m", "2m",
	}

	for i, key := range dbEnvKeyList {
		t.Setenv(key, inputList[i])
	}

	dbCfg, err := config.NewDatabaseConfig()

	assert.Nil(t, dbCfg)
	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Sprintf("environment variable %s must be set", driverKey))
}

func TestNewDatabaseConfigInvalidInt(t *testing.T) {
	invalidInt := "openConn"
	inputList := []string{
		"host", "5432", "user", "pass", "dbname", "disable", "postgres",
		invalidInt, "10", "5m", "2m",
	}

	for i, key := range dbEnvKeyList {
		t.Setenv(key, inputList[i])
	}

	dbCfg, err := config.NewDatabaseConfig()

	assert.Nil(t, dbCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("invalid int format: '%s'", invalidInt))
}
