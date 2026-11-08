package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Customer struct{
	Name string `json:"full_name"`
	City string `json:"city"`
	Zipcode string `json:"zip_code"`
}

func main(){
	// define routes
	http.HandleFunc("/greet",greet)
	http.HandleFunc("/customers",getAllCustomers)

	// starting server
	log.Fatal(http.ListenAndServe("localhost:8000",nil))
}

func greet(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w,"Hello World!")
}

func getAllCustomers(w http.ResponseWriter, r *http.Request){
	customers:= []Customer{
		{"Juma","New Mbeya","56276"},
		{"Mizo","New Dar","889276"},
	}

	w.Header().Add("Content-Type","application/json")

	json.NewEncoder(w).Encode(customers)
}