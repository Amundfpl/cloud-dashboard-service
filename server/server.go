package server

import (
	"github.com/amundfpl/Assignment-2/cache"
	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/utils"
	"log"
	"net/http"
	"os"
)

// StartServer initializes services, sets up routes, and runs the HTTP server.
func StartServer() {
	// Initialize Firebase and Firestore clients
	if dbInitErr := DatabaseInitialization(); dbInitErr != nil {
		log.Fatalf(utils.ErrMsgInitDB, dbInitErr)
	}

	// Ensure Firestore client closes on shutdown
	defer func() {
		if closeErr := db.CloseFirestore(); closeErr != nil {
			log.Printf(utils.ErrMsgCloseFirestore, closeErr)
		}
	}()

	// Start cache purge loop in background
	go cache.StartCachePurgeLoop()

	// Determine port from environment variable
	port := os.Getenv(utils.EnvPort)
	if port == "" {
		log.Println("$PORT not set. Defaulting to", utils.DefaultPort)
		port = utils.DefaultPort
	}

	// Initialize route handlers
	router := InitializeRoutes()

	log.Println(utils.MsgServerStart, port)

	// Start the HTTP server
	if listenErr := http.ListenAndServe(utils.AddrPrefix+port, router); listenErr != nil {
		log.Fatalf(utils.ErrMsgServerStart, listenErr)
	}
}
