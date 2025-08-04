package web_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/takumi616/go-restapi/infrastructure/web"
	"github.com/takumi616/go-restapi/shared/config"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func TestRunShutdownGracefully(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", dummyHandler)

	appConf := &config.AppConfig{
		Port: "0",
		Timeout: config.TimeoutConfig{
			ReadTimeout:       1 * time.Second,
			ReadHeaderTimeout: 1 * time.Second,
			WriteTimeout:      1 * time.Second,
			IdleTimeout:       1 * time.Second,
		},
	}

	server := web.NewServer(appConf, mux)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := server.Run(ctx)
		assert.Nil(t, err)
	}()

	time.Sleep(200 * time.Millisecond)

	cancel()

	time.Sleep(300 * time.Millisecond)
}

func TestRunListenerError(t *testing.T) {
	appConf := &config.AppConfig{
		Port:    "invalid-port",
		Timeout: config.TimeoutConfig{},
	}

	server := web.NewServer(appConf, http.NewServeMux())
	err := server.Run(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create http listener")
}
