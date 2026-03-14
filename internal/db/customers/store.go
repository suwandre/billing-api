package customers

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CustomerStore interface {
	Create(ctx context.Context, customer *Customer) (*Customer, error)
}

type customerStore struct {
	db *sql.DB
}

func (c *customerStore) Create(ctx context.Context, customer *Customer) (*Customer, error) {
	query := `
		INSERT INTO customers(email, username, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, username, created_at, updated_at
	`

	if customer.CreatedAt.IsZero() {
		customer.CreatedAt = time.Now()
	}
	if customer.UpdatedAt.IsZero() {
		customer.UpdatedAt = time.Now()
	}

	row := c.db.QueryRowContext(ctx, query,
		customer.Email,
		customer.Username,
		customer.PasswordHash,
		customer.CreatedAt,
		customer.UpdatedAt,
	)

	var created Customer
	err := row.Scan(
		&created.ID,
		&created.Email,
		&created.Username,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &created, nil
}
