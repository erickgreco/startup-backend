package domain

import (
	"database/sql/driver"
	"fmt"

	"github.com/shopspring/decimal"
)

/*
Representa dinero en unidad monetaria mayor (10.50 = $10.50)
Se persiste como BIGINT en centavos; la conversion decimal <-> centavos
vive en Value/Scan
*/
type Money decimal.Decimal

var ZeroMoney = Money(decimal.Zero)

// * WRITE
func (m Money) ToCents() int64 {
	return decimal.Decimal(m).Mul(decimal.NewFromInt(100)).Round(0).IntPart()
}

// * READ
func FromCents(cents int64) Money {
	return Money(decimal.NewFromInt(cents).Div(decimal.NewFromInt(100)))
}

// * Money -> BIGINT(cents) al ESCRIBIR, pgx lo llama solo
func (m Money) Value() (driver.Value, error) {
	return m.ToCents(), nil
}

// * BIGINT(cents) -> Money al LEER, pgx lo llama solo
// * Receiver PUNTERO: modifica el valor; con receiver de valor se perderia
func (m *Money) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*m = FromCents(v)
		return nil
	case nil:
		*m = ZeroMoney
		return nil
	default:
		return fmt.Errorf("cannot_scan %T into_money", src)
	}
}

// * Money -> JSON
// * Delega en decimal que serializa como numero ("20000.50")
func (m Money) MarshalJSON() ([]byte, error) {
	return decimal.Decimal(m).MarshalJSON()
}

// * JSON -> Money
// * Receiver PUNTERO: modifica el valor
func (m *Money) UnmarshalJSON(data []byte) error {
	var d decimal.Decimal

	if err := d.UnmarshalJSON(data); err != nil {
		return err
	}

	*m = Money(d)

	return nil
}

func (m Money) LessThanOrEqual(m2 Money) bool {
	return decimal.Decimal(m).LessThanOrEqual(decimal.Decimal(m2))
}

func (m Money) Mul(m2 decimal.Decimal) Money {
	return Money(decimal.Decimal(m).Mul(m2))
}

func (m Money) Add(m2 Money) Money {
	return Money(decimal.Decimal(m).Add(decimal.Decimal(m2)))
}

func (m Money) HasSubCentPrecision() bool {
	d := decimal.Decimal(m)

	return !d.Equal(d.Truncate(2))
}
