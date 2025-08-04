package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/takumi616/go-restapi/infrastructure/web"
	"github.com/takumi616/go-restapi/shared/config"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "REST API sample in Go")
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", testHandler)

	appCfg, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	server := web.NewServer(appCfg, mux)
	return server.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("Golang API server does not work correctly: %v", err)
	}
}
