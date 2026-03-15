package plans

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlanPricing struct {
	ID        uuid.UUID   `json:"id"`
	PlanID    uuid.UUID   `json:"plan_id"`
	Type      PricingType `json:"type"` // Monthly or yearly
	Price     float64     `json:"price"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type PlanResponse struct {
	ID        uuid.UUID     `json:"id"`
	Name      string        `json:"name"`
	Pricings  []PlanPricing `json:"pricings"` // will be populated by JOIN
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type PricingType uint8

const (
	Monthly PricingType = iota
	Yearly
)

func (p PricingType) String() string {
	switch p {
	case Monthly:
		return "monthly"
	case Yearly:
		return "yearly"
	default:
		return "unknown"
	}
}
