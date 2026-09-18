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
		return Currency{}, fmt.Errorf("invalid_currency %q: %w", code, err)
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
		return fmt.Errorf("invalid_currency %q from_db: %w", s, err)
	}

	c.unit = u

	return nil
}
