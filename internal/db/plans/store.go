package plans

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlanStore interface {
	Create(ctx context.Context, plan *Plan) (*Plan, error)
	CreatePricing(ctx context.Context, pricing *PlanPricing) (*PlanPricing, error)
	List(ctx context.Context) ([]PlanResponse, error) // Includes plan pricing
}

type planStore struct {
	pool *pgxpool.Pool
}

func NewPlanStore(pool *pgxpool.Pool) PlanStore {
	return &planStore{pool: pool}
}

func (s *planStore) Create(ctx context.Context, plan *Plan) (*Plan, error) {
	query := `
		INSERT INTO plans(name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, name, created_at, updated_at
	`

	if plan.CreatedAt.IsZero() {
		plan.CreatedAt = time.Now()
	}
	if plan.UpdatedAt.IsZero() {
		plan.UpdatedAt = time.Now()
	}

	row := s.pool.QueryRow(ctx, query,
		plan.Name,
		plan.CreatedAt,
		plan.UpdatedAt,
	)

	var created Plan
	err := row.Scan(
		&created.ID,
		&created.Name,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	return &created, nil
}

func (s *planStore) CreatePricing(ctx context.Context, pricing *PlanPricing) (*PlanPricing, error) {
	query := `
		INSERT INTO plan_pricings(plan_id, type, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, plan_id, type, price, created_at, updated_at
	`

	if pricing.CreatedAt.IsZero() {
		pricing.CreatedAt = time.Now()
	}
	if pricing.UpdatedAt.IsZero() {
		pricing.UpdatedAt = time.Now()
	}

	row := s.pool.QueryRow(ctx, query,
		pricing.PlanID,
		pricing.Type,
		pricing.Price,
		pricing.CreatedAt,
		pricing.UpdatedAt,
	)

	var created PlanPricing
	err := row.Scan(
		&created.ID,
		&created.PlanID,
		&created.Type,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create plan pricing: %w", err)
	}

	return &created, nil
}

func (s *planStore) List(ctx context.Context) ([]PlanResponse, error) {
	query := `
    SELECT 
				s.id, s.name, s.created_at, s.updated_at,
				sp.id, sp.plan_id, sp.type, sp.price,
				sp.created_at, sp.updated_at
		FROM plans s
		LEFT JOIN plan_pricings sp ON s.id = sp.plan_id
		ORDER BY s.id, sp.id
  `

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	defer rows.Close()

	subs := make(map[uuid.UUID]PlanResponse)

	for rows.Next() {
		var (
			subID          uuid.UUID
			subName        string
			subCreatedAt   time.Time
			subUpdatedAt   time.Time
			pricingID      *string
			pricingSubID   *string
			pricingType    *int16
			pricingPrice   *float64
			pricingCreated *time.Time
			pricingUpdated *time.Time
		)

		err := rows.Scan(
			&subID,
			&subName,
			&subCreatedAt,
			&subUpdatedAt,
			&pricingID,
			&pricingSubID,
			&pricingType,
			&pricingPrice,
			&pricingCreated,
			&pricingUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan: %w", err)
		}

		if _, exists := subs[subID]; !exists {
			subs[subID] = PlanResponse{
				ID:        subID,
				Name:      subName,
				CreatedAt: subCreatedAt,
				UpdatedAt: subUpdatedAt,
				Pricings:  []PlanPricing{},
			}
		}

		if pricingID != nil {
			pricing := PlanPricing{
				ID:        uuid.MustParse(*pricingID),
				PlanID:    uuid.MustParse(*pricingSubID),
				Type:      PricingType(*pricingType),
				Price:     *pricingPrice,
				CreatedAt: *pricingCreated,
				UpdatedAt: *pricingUpdated,
			}
			sub := subs[subID]
			sub.Pricings = append(sub.Pricings, pricing)
			subs[subID] = sub
		}
	}

	result := make([]PlanResponse, 0, len(subs))
	for _, sub := range subs {
		result = append(result, sub)
	}

	return result, nil
}
