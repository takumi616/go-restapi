package main

import (
	"context"
	"go-restapi/server"

)

func main() {
	//run http server
	server.RunHTTPServer(context.Background())
}