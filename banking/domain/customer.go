package domain

import "banking/errs"

type Customer struct {
	Id          string `json:"customer_id" xml:"id" db:"customer_id"`
	Name        string `json:"full_name" xml:"name" db:"name"`
	City        string `json:"city" xml:"city" db:"city"`
	Zipcode     string `json:"zip_code" xml:"zipcode" db:"zipcode"`
	DateofBirth string `json:"date_of_birth" xml:"date_of_birth" db:"date_of_birth"`
	Status      string `json:"status" xml:"status" db:"status"`
}

type CustomerRepository interface {
	FindAll(status string) ([]Customer, error)
	ById(id string) (*Customer, *errs.AppError)
}