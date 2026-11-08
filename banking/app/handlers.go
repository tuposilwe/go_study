package app

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
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

// CustomerResponse wraps a single customer to provide a root element (<customer>) for XML
type CustomerResponse struct {
	XMLName xml.Name `xml:"customer"`
	Customer
}

// 10 random Tanzanian records for dummy data
var customers = []Customer{
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

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {

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

func getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["customer_id"]
	
	// Convert ID string to integer to simulate database lookup
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	
	// Hardcoded dummy check against our 10 items
	if err != nil || id < 1 || id > 10 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Fetch specific customer (adjusting for 0-indexed slice)
	customer := customers[id-1]

	// Content negotiation for individual customer
	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		response := CustomerResponse{Customer: customer}
		if err := xml.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(customer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	var newCustomer Customer

	// Content negotiation for incoming payload type
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/xml" {
		if err := xml.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
			http.Error(w, "Invalid XML payload: "+err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Default fallback to JSON
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // Optional: rejects malformed/unexpected fields
		if err := decoder.Decode(&newCustomer); err != nil {
			http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Validate required fields
	if newCustomer.Name == "" || newCustomer.City == "" {
		http.Error(w, "Name and City are required fields", http.StatusUnprocessableEntity)
		return
	}

	// Simulate saving to a database
	log.Printf("Created new customer successfully: %+v\n", newCustomer)

	// Respond with 201 Created and return the newly created resource
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	
	if err := json.NewEncoder(w).Encode(newCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}