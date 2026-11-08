package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Import the Gorilla Mux package
	"github.com/gorilla/mux"
)

func Start() {
	// 1. Initialize the Gorilla Mux router instead of http.NewServeMux()
	router := mux.NewRouter()
	
	// Define your routes using the new router
	router.HandleFunc("/greet", greet).Methods("GET")
	router.HandleFunc("/customers", getAllCustomers).Methods("GET")
	router.HandleFunc("/customers", createCustomer).Methods("POST") // Added POST route
	router.HandleFunc("/customers/{customer_id:[0-9]+}", getCustomer).Methods("GET")

	// 2. Configure the HTTP server parameters explicitly, passing the Gorilla router as the Handler
	srv := &http.Server{
		Addr:    "localhost:8000",
		Handler: router, // Gorilla Mux router implements the http.Handler interface
	}

	// 3. Start the server inside a goroutine so it runs in the background
	go func() {
		log.Println("Server is starting on http://localhost:8000 ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	// 4. Create a channel to catch OS interrupt/termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// This line blocks execution until a signal drops into the 'stop' channel
	sig := <-stop
	log.Printf("Received signal: %v. Initiating graceful shutdown...\n", sig)

	// 5. Enforce a maximum timeout deadline for in-flight requests to finish up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 6. Tell the server to reject new incoming connections and finish active jobs
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown with active connections: %v", err)
	}

	log.Println("Server successfully exited clean.")
}