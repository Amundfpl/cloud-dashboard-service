package main

import (
	"github.com/amundfpl/Assignment-2/server"
)

// main is the entry point of the application.
// It logs startup and delegates to the server package to launch the HTTP server.
func main() {
	server.StartServer()
}
