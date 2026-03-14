package plans

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SubscriptionStore interface {
	CreateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error)
	CreateSubscriptionPricing(ctx context.Context, pricing *SubscriptionPricing) (*SubscriptionPricing, error)
	List(ctx context.Context) ([]SubscriptionResponse, error) // Includes subscription pricing
}

type subscriptionStore struct {
	db *sql.DB
}

func NewSubscriptionStore(db *sql.DB) SubscriptionStore {
	return &subscriptionStore{db: db}
}

func (s *subscriptionStore) CreateSubscription(ctx context.Context, subscription *Subscription) (*Subscription, error) {
	query := `
		INSERT INTO subscriptions(name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id, name, created_at, updated_at
	`

	if subscription.CreatedAt.IsZero() {
		subscription.CreatedAt = time.Now()
	}
	if subscription.UpdatedAt.IsZero() {
		subscription.UpdatedAt = time.Now()
	}

	row := s.db.QueryRowContext(ctx, query,
		subscription.Name,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	var created Subscription
	err := row.Scan(
		&created.ID,
		&created.Name,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return &created, nil
}

func (s *subscriptionStore) CreateSubscriptionPricing(ctx context.Context, pricing *SubscriptionPricing) (*SubscriptionPricing, error) {
	query := `
		INSERT INTO subscription_pricings(subscription_id, type, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, subscription_id, type, price, created_at, updated_at
	`

	row := s.db.QueryRowContext(ctx, query,
		pricing.SubscriptionID,
		pricing.Type,
		pricing.Price,
		pricing.CreatedAt,
		pricing.UpdatedAt,
	)

	var created SubscriptionPricing
	err := row.Scan(
		&created.ID,
		&created.SubscriptionID,
		&created.Type,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription pricing: %w", err)
	}

	return &created, nil
}

func (s *subscriptionStore) List(ctx context.Context) ([]SubscriptionResponse, error) {
	query := `
    SELECT 
				s.id, s.name, s.created_at, s.updated_at,
				sp.id, sp.subscription_id, sp.type, sp.price, 
				sp.created_at, sp.updated_at
		FROM subscriptions s
		LEFT JOIN subscription_pricings sp ON s.id = sp.subscription_id
		ORDER BY s.id, sp.id
  `

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	subs := make(map[uuid.UUID]SubscriptionResponse)

	for rows.Next() {
		var (
			subID          uuid.UUID
			subName        string
			subCreatedAt   time.Time
			subUpdatedAt   time.Time
			pricingID      sql.NullString
			pricingSubID   sql.NullString
			pricingType    sql.NullInt16
			pricingPrice   sql.NullFloat64
			pricingCreated sql.NullTime
			pricingUpdated sql.NullTime
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
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		if _, exists := subs[subID]; !exists {
			subs[subID] = SubscriptionResponse{
				ID:        subID,
				Name:      subName,
				CreatedAt: subCreatedAt,
				UpdatedAt: subUpdatedAt,
				Pricings:  []SubscriptionPricing{},
			}
		}

		if pricingID.Valid {
			pricing := SubscriptionPricing{
				ID:             uuid.MustParse(pricingID.String),
				SubscriptionID: uuid.MustParse(pricingSubID.String),
				Type:           PricingType(pricingType.Int16),
				Price:          pricingPrice.Float64,
				CreatedAt:      pricingCreated.Time,
				UpdatedAt:      pricingUpdated.Time,
			}
			sub := subs[subID]
			sub.Pricings = append(sub.Pricings, pricing)
			subs[subID] = sub
		}
	}

	result := make([]SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		result = append(result, sub)
	}

	return result, nil
}
