package quoteserv

import (
	"fmt"
	"strings"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/assets/fonts"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/projdocument"
	"github.com/gpdf-dev/gpdf"
	"github.com/gpdf-dev/gpdf/document"
	"github.com/gpdf-dev/gpdf/pdf"
	"github.com/gpdf-dev/gpdf/template"
)

// Paleta del documento
var (
	colorPrimary = pdf.RGBHex(0x1F2A44)
	colorMuted   = pdf.Gray(0.45)
	colorBorder  = pdf.Gray(0.85)
	colorPanel   = pdf.RGBHex(0xF5F6F8)

	cellPadding = document.Edges{
		Top: document.Mm(2), Bottom: document.Mm(2), Left: document.Mm(2), Right: document.Mm(2),
	}
)

// TODO: agregar opcion para escoger font, color, etc
func RenderQuotePDF(q *domain.QuotesModel, up *domain.UserProfileModel) ([]byte, error) {
	doc := gpdf.NewDocument(
		gpdf.WithPageSize(gpdf.A4),
		template.WithMargins(document.UniformEdges(document.Mm(18))),
		gpdf.WithFont("Inter", fonts.InterRegular),
		gpdf.WithFont("Inter-Bold", fonts.InterSemiBold),
		template.WithDefaultFont("Inter", 9),
		template.WithMetadata(document.DocumentMetadata{
			Title:  q.QuoteNumber,
			Author: fromName(up),
		}),
	)

	title, err := projdocument.SetTitleLang(up.Language, q.QuoteNumber)
	if err != nil {
		return nil, err
	}

	createdLabel, validLabel := projdocument.DateLabels(up.Language)
	createdAtDate := projdocument.FormatDate(&q.CreatedAt, up.Language)
	validUntilDate := projdocument.FormatDate(q.ValidUntil, up.Language)

	// Header - se repite en cada pagina
	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(7, func(c *template.ColBuilder) {
				c.Text(fromName(up), template.FontSize(16), template.Bold(),
					template.TextColor(colorPrimary))
				if up.WebsiteURL != nil {
					c.Text(*up.WebsiteURL, template.FontSize(8), template.TextColor(colorMuted))
				}
			})
			r.Col(5, func(c *template.ColBuilder) {
				c.Text(title, template.FontSize(13), template.Bold(), template.AlignRight(),
					template.TextColor(colorPrimary))
				c.Spacer(document.Mm(1))
				c.Text(createdLabel+createdAtDate, template.FontSize(8), template.AlignRight(),
					template.TextColor(colorMuted))
				c.Text(validLabel+validUntilDate, template.FontSize(8), template.AlignRight(),
					template.TextColor(colorMuted))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Spacer(document.Mm(3))
				c.Line(template.LineColor(colorPrimary), template.LineThickness(document.Pt(1.5)))
				c.Spacer(document.Mm(6))
			})
		})
	})

	// Footer - numero de pagina
	doc.Footer(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Line(template.LineColor(colorBorder))
				c.Spacer(document.Mm(2))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(8, func(c *template.ColBuilder) {
				c.Text(q.QuoteNumber, template.FontSize(7), template.TextColor(colorMuted))
			})
			r.Col(4, func(c *template.ColBuilder) {
				c.PageNumber(template.FontSize(7), template.AlignRight(),
					template.TextColor(colorMuted))
			})
		})
	})

	page := doc.AddPage()

	// Parties - dos paneles lado a lado
	fromLabel, toLabel := projdocument.PartyLabels(up.Language)

	fromData := []string{up.Email, up.PhoneNumber}

	toData := []string{}
	if q.ClientBusiness != nil {
		toData = append(toData, *q.ClientBusiness)
	}
	toData = append(toData, q.ClientEmail, q.ClientPhone)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(6, func(c *template.ColBuilder) {
			partyPanel(c, fromLabel, fromName(up), fromData, document.Edges{Right: document.Mm(2)})
		})
		r.Col(6, func(c *template.ColBuilder) {
			partyPanel(c, toLabel, q.ClientName, toData, document.Edges{Left: document.Mm(2)})
		})
	})

	sectionSpacer(page, 8)

	// Project
	sectionLabel, descLabel, resumeLabel, durationLabel := projdocument.ProjectLabels(up.Language)

	sectionTitle(page, sectionLabel)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(9, func(c *template.ColBuilder) {
			c.Text(q.ProjectTitle, template.FontSize(11), template.Bold())
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(durationLabel, template.FontSize(7), template.AlignRight(),
				template.TextColor(colorMuted))
			c.Text(projdocument.DurationText(up.Language, q.ProjectEstimatedDurationDays),
				template.FontSize(10), template.Bold(), template.AlignRight())
		})
	})

	sectionSpacer(page, 3)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			labeledParagraph(c, descLabel, q.ProjectDescription)
			if q.ProjectResume != nil {
				c.Spacer(document.Mm(3))
				labeledParagraph(c, resumeLabel, *q.ProjectResume)
			}
		})
	})

	sectionSpacer(page, 8)

	// Items
	descH, qtyH, unitH, priceH, subtotalH := projdocument.ItemHeaders(up.Language)
	cur := q.Currency.String()

	headerOpts := []template.TextOption{template.FontSize(8), template.Bold(),
		template.TextColor(pdf.White), template.BgColor(colorPrimary), template.TextPadding(cellPadding)}

	itemRow(page,
		[5]string{descH, qtyH, unitH, priceH, subtotalH},
		headerOpts,
	)

	for _, item := range q.Items {
		itemRow(page,
			[5]string{
				item.Description,
				item.Quantity.String(),
				item.Unit,
				moneyText(cur, item.UnitPrice),
				moneyText(cur, item.SubTotal),
			},
			[]template.TextOption{template.FontSize(9), template.TextPadding(cellPadding)},
		)

		page.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Line(template.LineColor(colorBorder), template.LineThickness(document.Pt(0.5)))
			})
		})
	}

	sectionSpacer(page, 4)

	// Totals - alineados a la derecha
	subtotalL, totalL := projdocument.TotalLabels(up.Language)

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(7, func(c *template.ColBuilder) {})
		r.Col(2, func(c *template.ColBuilder) {
			c.Text(subtotalL, template.AlignRight(), template.TextColor(colorMuted))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(moneyText(cur, q.Subtotal), template.AlignRight())
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(7, func(c *template.ColBuilder) {})
		r.Col(5, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(2))
			c.Line(template.LineColor(colorPrimary), template.LineThickness(document.Pt(1)))
			c.Spacer(document.Mm(2))
		})
	})

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(7, func(c *template.ColBuilder) {})
		r.Col(2, func(c *template.ColBuilder) {
			c.Text(totalL, template.FontSize(12), template.Bold(), template.AlignRight(),
				template.TextColor(colorPrimary))
		})
		r.Col(3, func(c *template.ColBuilder) {
			c.Text(moneyText(cur, q.Total), template.FontSize(12), template.Bold(),
				template.AlignRight(), template.TextColor(colorPrimary))
		})
	})

	sectionSpacer(page, 12)

	// Terms
	disclosure := Disclosure(up.Language, validUntilDate, cur, q.PaymentMethod.Translate(up.Language))

	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Box(func(b *template.ColBuilder) {
				b.Text(termsLabel(up.Language), template.FontSize(8), template.Bold(),
					template.TextColor(colorPrimary))
				b.Spacer(document.Mm(1.5))
				b.Text(disclosure, template.FontSize(8), template.TextColor(colorMuted))
			},
				template.WithBoxBackground(colorPanel),
				template.WithBoxPadding(document.UniformEdges(document.Mm(4))),
				template.WithBoxBorder(template.Border(
					template.BorderWidths(document.Pt(0), document.Pt(0), document.Pt(0), document.Pt(2)),
					template.BorderColor(colorPrimary),
				)),
			)
		})
	})

	return doc.Generate()
}

// itemRow dibuja una fila de la tabla de conceptos: Desc(4) Cant(2) Unidad(2) Precio(2) Subtotal(2)
func itemRow(page *template.PageBuilder, cells [5]string, opts []template.TextOption) {
	spans := [5]int{4, 2, 2, 2, 2}
	aligns := [5]template.TextOption{
		template.AlignLeft(),
		template.AlignCenter(),
		template.AlignCenter(),
		template.AlignRight(),
		template.AlignRight(),
	}

	page.AutoRow(func(r *template.RowBuilder) {
		for i := range cells {
			r.Col(spans[i], func(c *template.ColBuilder) {
				c.Text(cells[i], append(append([]template.TextOption{}, opts...), aligns[i])...)
			})
		}
	})
}

// partyPanel dibuja el bloque De/Para: etiqueta, nombre y datos de contacto
func partyPanel(c *template.ColBuilder, label, name string, lines []string, margin document.Edges) {
	c.Box(func(b *template.ColBuilder) {
		b.Text(strings.ToUpper(strings.TrimSuffix(label, ":")), template.FontSize(7),
			template.Bold(), template.LetterSpacing(0.8), template.TextColor(colorMuted))
		b.Spacer(document.Mm(1.5))
		b.Text(name, template.FontSize(10), template.Bold())
		for _, line := range lines {
			if line == "" {
				continue
			}
			b.Text(line, template.FontSize(9))
		}
	},
		template.WithBoxBackground(colorPanel),
		template.WithBoxPadding(document.UniformEdges(document.Mm(4))),
		template.WithBoxMargin(margin),
	)
}

// sectionTitle titulo de seccion con linea inferior
func sectionTitle(page *template.PageBuilder, title string) {
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Text(strings.ToUpper(title), template.FontSize(9), template.Bold(),
				template.LetterSpacing(0.8), template.TextColor(colorPrimary))
			c.Spacer(document.Mm(1.5))
			c.Line(template.LineColor(colorBorder), template.LineThickness(document.Pt(0.75)))
			c.Spacer(document.Mm(3))
		})
	})
}

func sectionSpacer(page *template.PageBuilder, mm float64) {
	page.AutoRow(func(r *template.RowBuilder) {
		r.Col(12, func(c *template.ColBuilder) {
			c.Spacer(document.Mm(mm))
		})
	})
}

func labeledParagraph(c *template.ColBuilder, label, text string) {
	c.Text(label, template.FontSize(8), template.Bold(), template.TextColor(colorMuted))
	c.Spacer(document.Mm(0.5))
	c.Text(text, template.FontSize(9))
}

// moneyText formatea "MXN 15,000.00"
func moneyText(cur string, m domain.Money) string {
	s := m.String()
	sign := ""
	if strings.HasPrefix(s, "-") {
		sign, s = "-", s[1:]
	}

	intPart, decPart, _ := strings.Cut(s, ".")

	var b strings.Builder
	for i, d := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(d)
	}

	return cur + " " + sign + b.String() + "." + decPart
}

func termsLabel(lang string) string {
	switch lang {
	case "es":
		return "Términos y condiciones"
	default:
		return "Terms and conditions"
	}
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
