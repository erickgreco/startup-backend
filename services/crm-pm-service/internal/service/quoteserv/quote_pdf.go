package quoteserv

import (
	"fmt"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/assets/fonts"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/projdocument"
	"github.com/gpdf-dev/gpdf"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/template"
)

func RenderQuotePDF(q *domain.QuotesModel, up *domain.UserProfileModel) ([]byte, error) {
	doc := gpdf.NewDocument(gpdf.WithPageSize(gpdf.A4), gpdf.WithFont("Inter", fonts.InterLight))

	title, err := projdocument.SetTitleLang(up.Language, q.QuoteNumber)
	if err != nil {
		return nil, err
	}

	createdLabel, validLabel := projdocument.DateLabels(up.Language)
	createdAtDate := projdocument.FormatDate(&q.CreatedAt, up.Language)
	validUntilDate := projdocument.FormatDate(q.ValidUntil, up.Language)

	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text(title, template.AlignCenter(), template.FontSize(18),
					template.Bold())
				c.Spacer(document.Mm(3))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text(createdLabel+createdAtDate, template.AlignLeft(), template.FontSize(9))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text(validLabel+validUntilDate, template.AlignRight(), template.FontSize(9))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Spacer(document.Mm(5))
			})
		})
	})

	fromLabel, toLabel := projdocument.PartyLabels(up.Language)

	page := doc.AddPage()

	// Parties - filas
	fromData := []string{fromName(up), up.Email, up.PhoneNumber}
	if up.WebsiteURL != nil {
		fromData = append(fromData, *up.WebsiteURL)
	}

	toData := []string{q.ClientName}
	if q.ClientBusiness != nil {
		toData = append(toData, *q.ClientBusiness)
	}
	toData = append(toData, q.ClientEmail, q.ClientPhone)

	for len(fromData) < len(toData) {
		fromData = append(fromData, "")
	}
	for len(toData) < len(fromData) {
		toData = append(toData, "")
	}

	// Filas como pares
	var partyRows [][]string
	for i := range fromData {
		partyRows = append(partyRows, []string{fromData[i], toData[i]})
	}

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Table(
				[]string{fromLabel, toLabel},
				partyRows,
				template.TableHeaderStyle(template.Bold(), template.FontSize(11)),
				template.ColumnWidths(0.5, 0.5),
			)
		})
	})

	// Project
	sectionLabel, descLabel, resumeLabel, durationLabel := projdocument.ProjectLabels(up.Language)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(sectionLabel, template.FontSize(14), template.Bold())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(q.ProjectTitle, template.FontSize(11), template.Bold())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(descLabel, template.FontSize(10), template.Bold())
			c.Text(q.ProjectDescription, template.FontSize(9))
		})
	})

	if q.ProjectResume != nil {
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text(resumeLabel, template.FontSize(10), template.Bold())
				c.Text(*q.ProjectResume, template.FontSize(9))
			})
		})
	}

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(durationLabel+": "+projdocument.DurationText(up.Language, q.ProjectEstimatedDurationDays), template.FontSize(9))
		})
	})

	descH, qtyH, unitH, priceH, subtotalH := projdocument.ItemHeaders(up.Language)
	cur := q.Currency.String()

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(4, func(c *template.ColBuilder) {
			c.Text(descH, template.FontSize(9), template.Bold())
		})
		r.Col(1, func(c *template.ColBuilder) {
			c.Text(qtyH, template.FontSize(9), template.Bold(), template.AlignCenter())
		})
		r.Col(2, func(c *template.ColBuilder) {
			c.Text(unitH, template.FontSize(9), template.Bold(), template.AlignCenter())
		})
		r.Col(2, func(c *template.ColBuilder) {
			c.Text(priceH, template.FontSize(9), template.Bold(), template.AlignRight())
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(subtotalH, template.FontSize(9), template.Bold(), template.AlignRight())
		})
	})

	for _, item := range q.Items {
		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(4, func(c *template.ColBuilder) {
				c.Text(item.Description, template.FontSize(9))
			})
			r.Col(1, func(c *template.ColBuilder) {
				c.Text(item.Quantity.String(), template.FontSize(9), template.AlignCenter())
			})
			r.Col(2, func(c *template.ColBuilder) {
				c.Text(item.Unit, template.FontSize(9), template.AlignCenter())
			})
			r.Col(2, func(c *template.ColBuilder) {
				c.Text(cur+" "+item.UnitPrice.String(), template.FontSize(9), template.AlignRight())
			})
			r.Col(3, func(c *template.ColBuilder) {
				c.Text(cur+" "+item.SubTotal.String(), template.FontSize(9), template.AlignRight())
			})
		})
	}

	subtotalL, totalL := projdocument.TotalLabels(up.Language)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(9, func(c *template.ColBuilder) {})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(subtotalL+": "+cur+" "+q.Subtotal.String(), template.FontSize(9), template.AlignRight())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(9, func(c *template.ColBuilder) {})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(totalL+": "+cur+" "+q.Total.String(), template.FontSize(11), template.Bold(),
				template.AlignRight())
		})
	})

	disclousure := Disclosure(up.Language, validUntilDate, cur, q.PaymentMethod.Translate(up.Language))

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(disclousure, template.FontSize(8), template.AlignCenter())
		})
	})

	return doc.Generate()
}

func fromName(up *domain.UserProfileModel) string {
	if up.BusinessName != nil {
		return *up.BusinessName
	}

	return up.FirstName + " " + up.LastName
}

func Disclosure(lang, validUntil, currency, payMethod string) string {
	switch lang {
	case "es":
		return fmt.Sprintf(
			"Esta cotización es válida hasta el %s. Los precios están expresados en %s. Método de pago: %s. Cambios en el alcance del proyecto pueden generar ajustes en el precio final",
			validUntil, currency, payMethod,
		)
	default:
		return fmt.Sprintf(
			"This quote is valid until %s. Prices are expressed in %s. Payment method: %s. Changes in project scope may result in adjustments to the final price.",
			validUntil, currency, payMethod,
		)
	}
}
