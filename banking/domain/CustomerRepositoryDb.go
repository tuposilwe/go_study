package domain

import (
	"banking/errs"
	"banking/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepositoryDb struct {
	client *pgxpool.Pool
}

// NewCustomerRepositoryDb creates a repository using an existing pgxpool connection pool.
func NewCustomerRepositoryDb(pool *pgxpool.Pool) CustomerRepositoryDb {
	return CustomerRepositoryDb{client: pool}
}

// FindAll retrieves all customers, optionally filtered by status ("active", "inactive", or "").
func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		rows pgx.Rows
		err  error
	)

	baseQuery := `
		SELECT customer_id, name, city, zipcode, date_of_birth, status 
		FROM customers`

	if status == "" {
		rows, err = d.client.Query(ctx, baseQuery+" ORDER BY customer_id")
	} else {
		statusCode := 1
		if status == "inactive" {
			statusCode = 0
		}
		rows, err = d.client.Query(ctx, baseQuery+" WHERE status = $1 ORDER BY customer_id", statusCode)
	}

	if err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error querying customers: %w", err)
	}
	defer rows.Close()

	customers := make([]Customer, 0)
	for rows.Next() {
		var (
			c          Customer
			dob        time.Time
			statusCode int16
		)

		err := rows.Scan(
			&c.Id,
			&c.Name,
			&c.City,
			&c.Zipcode,
			&dob,
			&statusCode,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning customer row: %w", err)
		}

		c.DateofBirth = dob.Format("2006-01-02")
		if statusCode == 1 {
			c.Status = "active"
		} else {
			c.Status = "inactive"
		}

		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		logger.Error(err.Error())
		return nil, fmt.Errorf("error iterating customer rows: %w", err)
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

	var (
		c          Customer
		dob        time.Time
		statusCode int16
	)

	err := d.client.QueryRow(ctx, query, id).Scan(
		&c.Id,
		&c.Name,
		&c.City,
		&c.Zipcode,
		&dob,
		&statusCode,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error(err.Error())
			return nil, errs.NewNotFoundError("customer not found")
		}

		logger.Error(err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	c.DateofBirth = dob.Format("2006-01-02")
	if statusCode == 1 {
		c.Status = "active"
	} else {
		c.Status = "inactive"
	}

	return &c, nil
}
