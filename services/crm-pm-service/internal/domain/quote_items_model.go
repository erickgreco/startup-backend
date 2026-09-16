package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

/*
Ejemplos:

		-- Horas de trabajo
		Item {
			Description: 	"Desarrollo-backend",
			Quantity:		40,		horas
			UnitPrice:		50000, 	$500 USD/hora
			SubTotal:		2000000	$20,000
		}

		-- Numero de entregables
		Item {
			Description:	"Diseño-UI-mockups",
			Quantity:		10,		10 mockups
			UnitPrice:		5000,	$50 USD/mockup
			SubTotal:		50000	$500
		}

		-- Tiempo en meses
		Item {
	    	Description:		"Hosting-anual",
	    	Quantity: 12,		12 meses
	    	UnitPrice: 10000,	$100 USD/mes
	    	SubTotal: 120000	$1,200
		}

		-- Lump sum (cantidad = 1)
		Item {
			Description:		"Consultoria-estrategica",
	   		Quantity: 1,		1 proyecto completo
	   		UnitPrice: 500000,	$5,000 USD todo
	   		SubTotal: 500000	$5,000
		}
*/
type QuoteItem struct {
	ID          uuid.UUID
	QuoteID     uuid.UUID
	Description string
	Quantity    decimal.Decimal // * Puede ser duracion o cantidad
	UnitPrice   Money
	SubTotal    Money
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ItemRequest struct {
	Description string
	Quantity    decimal.Decimal
	UnitPrice   decimal.Decimal
}
