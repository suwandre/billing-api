package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suwandre/billing-api/internal/db/customers"
	"github.com/suwandre/billing-api/internal/db/plans"
)

type Store interface {
	Customers() customers.CustomerStore
	Subscriptions() plans.SubscriptionStore
}

type store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &store{pool: pool}
}

func (s *store) Customers() customers.CustomerStore {
	return customers.NewCustomerStore(s.pool)
}

func (s *store) Subscriptions() plans.SubscriptionStore {
	return plans.NewSubscriptionStore(s.pool)
}
