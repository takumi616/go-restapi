package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"go-restapi/config"
	"go-restapi/routing"
)

func RunHTTPServer(ctx context.Context) {
	//get config file
	config := config.GetConfig()
	//get router
	router := routing.RegistHandler()
	//set port and router to server struct
	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	//create ctx with signal
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM)
	defer stop()

	//run http server in groutine
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Printf("Http server error happened while calling ListenAndServe method.")
			stop()
		}
	}()

	//wait for ctx is canceled
	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println("Failed to shutdown server.")
	} else {
		log.Println("Successfully shutdown server.")
	}
}
