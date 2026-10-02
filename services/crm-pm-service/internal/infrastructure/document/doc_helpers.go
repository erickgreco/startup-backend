package document

import (
	"fmt"
	"time"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
)

func setTitleLang(lang, quoteNm string) (string, error) {
	switch {
	case lang == "es":
		return fmt.Sprintf("Cotización No. %s", quoteNm), nil
	case lang == "en":
		return fmt.Sprintf("Quote No. %s", quoteNm), nil
	default:
		return "", domain.ErrInvalidLanguage
	}
}

func formatDate(t *time.Time, lang string) string {
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

func dateLabels(lang string) (created, valid string) {
	switch lang {
	case "es":
		return "Fecha: ", "Válido hasta: "
	default:
		return "Date: ", "Valid until: "
	}
}
