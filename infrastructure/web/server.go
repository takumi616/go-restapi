package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/takumi616/go-restapi/shared/config"
)

type Server struct {
	Port       string
	HttpServer *http.Server
}

func NewServer(appConf *config.AppConfig, mux http.Handler) *Server {
	return &Server{
		Port: appConf.Port,
		HttpServer: &http.Server{
			Handler:           mux,
			ReadTimeout:       appConf.Timeout.ReadTimeout,
			ReadHeaderTimeout: appConf.Timeout.ReadHeaderTimeout,
			WriteTimeout:      appConf.Timeout.WriteTimeout,
			IdleTimeout:       appConf.Timeout.IdleTimeout,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return fmt.Errorf("failed to create http listener: %w", err)
	}

	// Run http server in another goroutine
	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if err := s.HttpServer.Serve(listener); err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Execute graceful shutdown
		shutdownErr := s.HttpServer.Shutdown(shutdownCtx)

		var serveErr error
		select {
		case serveErr = <-serverErrCh:
		case <-time.After(2 * time.Second):
		}

		switch {
		case shutdownErr != nil && serveErr != nil:
			return fmt.Errorf("shutdown error: %w; serve error: %w", shutdownErr, serveErr)
		case shutdownErr != nil:
			return fmt.Errorf("shutdown error: %w", shutdownErr)
		case serveErr != nil:
			return fmt.Errorf("serve error: %w", serveErr)
		default:
			return nil
		}

	case err := <-serverErrCh:
		return err
	}
}
