package domain

import (
	"database/sql/driver"
	"fmt"

	"golang.org/x/text/currency"
)

type Currency struct {
	unit currency.Unit
}

func ParseCurrency(code string) (Currency, error) {
	u, err := currency.ParseISO(code)
	if err != nil {
		return Currency{}, fmt.Errorf("%w - %q: %w", ErrInvalidCurrency, code, err)
	}

	return Currency{unit: u}, nil
}

func (c Currency) String() string {
	return c.unit.String()
}

func (c Currency) Value() (driver.Value, error) {
	return c.unit.String(), nil
}

func (c *Currency) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot_scan %T into_currency", src)
	}

	u, err := currency.ParseISO(s)
	if err != nil {
		return fmt.Errorf("%w - %q from_db: %w", ErrInvalidCurrency, s, err)
	}

	c.unit = u

	return nil
}
