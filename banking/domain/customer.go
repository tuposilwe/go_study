package domain

import (
	"banking/dto"
	"banking/errs"
)

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string `db:"name"`
	City        string `db:"city"`
	Zipcode     string `db:"zipcode"`
	DateofBirth string `db:"date_of_birth"`
	Status      string `db:"status"`
}

func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateofBirth: c.DateofBirth,
		Status:      c.Status,
	}
}

type CustomerRepository interface {
	FindAll(status string) ([]Customer, error)
	ById(id string) (*Customer, *errs.AppError)
}