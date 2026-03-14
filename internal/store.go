package internal

import (
	"database/sql"

	"github.com/suwandre/billing-api/internal/db/customers"
	"github.com/suwandre/billing-api/internal/db/plans"
)

type Store interface {
	Customers() customers.CustomerStore
	Subscriptions() plans.SubscriptionStore
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

func (s *store) Subscriptions() plans.SubscriptionStore {
	return plans.NewSubscriptionStore(s.db)
}
