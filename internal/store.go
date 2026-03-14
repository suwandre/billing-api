package internal

import (
	"database/sql"

	"github.com/suwandre/billing-api/internal/db/customers"
)

type Store interface {
	Customers() customers.CustomerStore
}

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) Customers() customers.CustomerStore {
	return customers.NewCustomerStore(s.db)
}
