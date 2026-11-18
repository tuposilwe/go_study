package domain

import (
	"banking/errs"
	"banking/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

// NewCustomerRepositoryDb creates a repository using an existing sqlx.DB connection.
func NewCustomerRepositoryDb(db *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{client: db}
}

// customerRow is an internal DTO used by sqlx to scan raw DB types before converting to Customer.
type customerRow struct {
	Id          string    `db:"customer_id"`
	Name        string    `db:"name"`
	City        string    `db:"city"`
	Zipcode     string    `db:"zipcode"`
	DateofBirth time.Time `db:"date_of_birth"`
	Status      int16     `db:"status"`
}

func (r customerRow) toDomain() Customer {
	status := "inactive"
	if r.Status == 1 {
		status = "active"
	}

	return Customer{
		Id:          r.Id,
		Name:        r.Name,
		City:        r.City,
		Zipcode:     r.Zipcode,
		DateofBirth: r.DateofBirth.Format("2006-01-02"),
		Status:      status,
	}
}

// FindAll retrieves all customers, optionally filtered by status ("active", "inactive", or "").
func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		rows []customerRow
		err  error
	)

	baseQuery := `
		SELECT customer_id, name, city, zipcode, date_of_birth, status 
		FROM customers`

	if status == "" {
		err = d.client.SelectContext(ctx, &rows, baseQuery+" ORDER BY customer_id")
	} else {
		statusCode := 1
		if status == "inactive" {
			statusCode = 0
		}
		err = d.client.SelectContext(ctx, &rows, baseQuery+" WHERE status = $1 ORDER BY customer_id", statusCode)
	}

	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error querying customers: %w", err)
	}

	customers := make([]Customer, 0, len(rows))
	for _, r := range rows {
		customers = append(customers, r.toDomain())
	}

	return customers, nil
}

// ById retrieves a single customer by customer_id.
func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT customer_id, name, city, zipcode, date_of_birth, status 
		FROM customers 
		WHERE customer_id = $1`

	var row customerRow
	err := d.client.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error(err.Error())
			return nil, errs.NewNotFoundError("customer not found")
		}

		logger.Error(err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	customer := row.toDomain()
	return &customer, nil
}