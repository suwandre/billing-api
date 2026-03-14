package customers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerStore interface {
	Create(ctx context.Context, customer *Customer) (*Customer, error)
	GetByEmail(ctx context.Context, email string) (*CustomerResponse, error)
}

type customerStore struct {
	pool *pgxpool.Pool
}

func NewCustomerStore(pool *pgxpool.Pool) CustomerStore {
	return &customerStore{pool: pool}
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

	row := c.pool.QueryRow(ctx, query,
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

func (c *customerStore) GetByEmail(ctx context.Context, email string) (*CustomerResponse, error) {
	query := `
		SELECT id, email, username, created_at, updated_at 
		FROM customers 
		WHERE email = $1
	`

	row := c.pool.QueryRow(ctx, query, email)

	var customer CustomerResponse
	err := row.Scan(
		&customer.ID,
		&customer.Email,
		&customer.Username,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer by email: %w", err)
	}

	return &customer, nil
}
