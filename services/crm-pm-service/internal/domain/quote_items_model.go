package domain

import (
	"time"

	"github.com/google/uuid"
)

type QuoteItem struct {
	ID          uuid.UUID
	QuoteID     uuid.UUID
	Description string
	Quantity    int
	UnitPrice   Money
	SubTotal    Money
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
