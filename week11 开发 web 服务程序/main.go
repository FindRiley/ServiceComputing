package main

import (
	"os"

	"github.com/github.user/CloudGo/server"
)

func main() {
	port := os.Getenv("PORT")
	// set default listening port
	if len(port) == 0 {
		port = "8088"
	}
	server.Start(port)
}
