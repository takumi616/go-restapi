package main

import (
	"context"
	"log/slog"

	"github.com/takumi616/go-restapi/application/usecase"
	"github.com/takumi616/go-restapi/infrastructure/db"
	"github.com/takumi616/go-restapi/infrastructure/db/repository"
	"github.com/takumi616/go-restapi/infrastructure/web"
	"github.com/takumi616/go-restapi/interface/gateway"
	"github.com/takumi616/go-restapi/interface/handler"
	"github.com/takumi616/go-restapi/shared/config"
	"github.com/takumi616/go-restapi/shared/logger"
)

func run(ctx context.Context) error {
	logger.Init()

	dbCfg, err := config.NewDatabaseConfig()
	if err != nil {
		return err
	}

	db, err := db.NewDBConnection(ctx, dbCfg)
	if err != nil {
		return err
	}

	appCfg, err := config.NewAppConfig()
	if err != nil {
		return err
	}

	repository := repository.NewTaskRepository(db)
	gateway := gateway.NewTaskGateway(repository)
	usecase := usecase.NewTaskUsecase(gateway)
	handler := handler.NewTaskHandler(usecase)

	serveMux := web.NewServeMux(handler)

	server := web.NewServer(appCfg, serveMux.RegisterHandler())
	return server.Run(ctx)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.ErrorContext(
			ctx,
			"golang API server does not work correctly",
			slog.String("err", err.Error()),
		)
	}
}
