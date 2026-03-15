package internal

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suwandre/billing-api/internal/db/customers"
	"github.com/suwandre/billing-api/internal/db/plans"
)

type Store interface {
	Customers() customers.CustomerStore
	Plans() plans.PlanStore
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

func (s *store) Plans() plans.PlanStore {
	return plans.NewPlanStore(s.pool)
}
