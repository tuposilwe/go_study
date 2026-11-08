package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Customer defines the structure for our data
type Customer struct {
	Name    string `json:"full_name" xml:"name"`
	City    string `json:"city" xml:"city"`
	Zipcode string `json:"zip_code" xml:"zipcode"`
}

// CustomersResponse wraps the slice to provide a valid XML root element (<customers>)
type CustomersResponse struct {
	XMLName   xml.Name   `xml:"customers"`
	Customers []Customer `xml:"customer"`
}

func main() {
	// 1. Create a clean request router
	mux := http.NewServeMux()
	mux.HandleFunc("/greet", greet)
	mux.HandleFunc("/customers", getAllCustomers)

	// 2. Configure the HTTP server parameters explicitly
	srv := &http.Server{
		Addr:    "localhost:8000",
		Handler: mux,
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

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	// 10 random Tanzanian records for dummy data
	customers := []Customer{
		{"Juma", "New Mbeya", "56276"},
		{"Mizo", "New Dar", "889276"},
		{"Neema", "Arusha", "23101"},
		{"Baraka", "Dodoma", "41101"},
		{"Subira", "Mwanza", "33102"},
		{"Bakari", "Tanga", "21103"},
		{"Faraja", "Zanzibar", "71101"},
		{"Amani", "Morogoro", "61104"},
		{"Rehema", "Iringa", "51102"},
		{"Chacha", "Musoma", "31105"},
	}

	// Use Content Negotiation by checking what format the client is asking for (Accept)
	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		
		response := CustomersResponse{Customers: customers}
		if err := xml.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Default fallback is JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(customers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
