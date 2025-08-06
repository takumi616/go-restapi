package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/takumi616/go-restapi/infrastructure/db"
	"github.com/takumi616/go-restapi/infrastructure/web"
	"github.com/takumi616/go-restapi/shared/config"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "REST API sample in Go")
}

func testDbConnection(w http.ResponseWriter, r *http.Request) {
	dbCfg, err := config.NewDatabaseConfig()
	if err != nil {
		log.Printf("failed to get database config: %v", err)
		return
	}

	ctx := r.Context()
	db, err := db.NewDBConnection(ctx, dbCfg)
	if err != nil {
		log.Printf("failed to open database: %v", err)
		return
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' -- 通常のスキーマ名
			AND table_name = $1
		);
	`

	var tableExists bool
	err = db.QueryRowContext(ctx, query, "tasks").Scan(&tableExists)
	if err != nil {
		log.Printf("no rows error: %v", err)
		return
	}

	if tableExists {
		fmt.Fprintf(w, "DB migration success")
	} else {
		log.Println("DB migration fail")
	}
}

func run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", testHandler)
	mux.HandleFunc("GET /test/db", testDbConnection)

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
