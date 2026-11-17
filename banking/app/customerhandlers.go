package app

import (
	"banking/service"
	"encoding/json"
	"encoding/xml"
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

type CustomerHandlers struct {
	service service.CustomerService
}

func (ch *CustomerHandlers) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	customers, err := ch.service.GetAllCustomers(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(customers); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(customers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ch *CustomerHandlers) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customer_id"]

	customer, err := ch.service.GetCustomer(id)
	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	if r.Header.Get("Accept") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		if err := xml.NewEncoder(w).Encode(customer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	writeResponse(w, http.StatusOK, customer)
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

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}

}
