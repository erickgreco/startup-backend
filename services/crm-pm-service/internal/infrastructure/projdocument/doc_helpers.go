package projdocument

import (
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
)

func SetTitleLang(lang, quoteNm string) (string, error) {
	switch {
	case lang == "es":
		return fmt.Sprintf("Cotización No. %s", quoteNm), nil
	case lang == "en":
		return fmt.Sprintf("Quote No. %s", quoteNm), nil
	default:
		return "", domain.ErrInvalidLanguage
	}
}

func FormatDate(t *time.Time, lang string) string {
	if t == nil {
		return "N/A"
	}

	switch lang {
	case "es":
		return fmt.Sprintf("%d/%s/%d", t.Day(), monthES(t.Month()), t.Year())
	default:
		return fmt.Sprintf("%d/%s/%d", t.Day(), t.Month().String(), t.Year())
	}
}

func monthES(m time.Month) string {
	months := [...]string{
		"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio",
		"agosto", "septiembre", "octubre", "noviembre", "diciembre",
	}

	return months[m-1]
}

func DateLabels(lang string) (created, valid string) {
	switch lang {
	case "es":
		return "Fecha: ", "Válido hasta: "
	default:
		return "Date: ", "Valid until: "
	}
}

func PartyLabels(lang string) (from, to string) {
	switch lang {
	case "es":
		return "De:", "Para:"
	default:
		return "From:", "To:"
	}
}

func ProjectLabels(lang string) (section, desc, resume, duration string) {
	switch lang {
	case "es":
		return "Proyecto", "Descripción", "Propuesta", "Duración estimada"
	default:
		return "Project", "Description", "Proposal", "Estimated duration"
	}
}

func ItemHeaders(lang string) (desc, qty, unit, price, subtotal string) {
	switch lang {
	case "es":
		return "Descripción", "Cantidad", "Unidad", "Precio unitario", "Subtotal"
	default:
		return "Description", "Quantity", "Unit", "Unit price", "Subtotal"
	}
}

func TotalLabels(lang string) (subtotal, total string) {
	switch lang {
	case "es":
		return "Subtotal", "Total"
	default:
		return "Subtotal", "Total"
	}
}

func DurationText(lang string, days int32) string {
	switch lang {
	case "es":
		return fmt.Sprintf("%d días", days)
	default:
		return fmt.Sprintf("%d days", days)
	}
}
