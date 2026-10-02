package document

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/gpdf-dev/gpdf"
	"github.com/gpdf-dev/gpdf/template"
)

func RenderQuotePDF(q *domain.QuotesModel, up *domain.UserProfileModel) ([]byte, error) {
	doc := gpdf.NewDocument(gpdf.WithPageSize(gpdf.A4))

	title, err := setTitleLang(up.Language, q.QuoteNumber)
	if err != nil {
		return nil, err
	}

	createdLabel, validLabel := dateLabels(up.Language)

	createdAt := formatDate(&q.CreatedAt, up.Language)
	validUntil := formatDate(q.ValidUntil, up.Language)

	doc.Header(func(p *template.PageBuilder) {
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(12, func(c *template.ColBuilder) {
				c.Text(title, template.AlignCenter(), template.FontSize(18),
					template.Bold(), template.FontFamily("Inter"))
			})
		})
		p.AutoRow(func(r *template.RowBuilder) {
			r.Col(6, func(c *template.ColBuilder) {
				c.Text(createdLabel+createdAt, template.AlignLeft(), template.FontSize(9),
					template.FontFamily("Inter"))
			})
			r.Col(6, func(c *template.ColBuilder) {
				c.Text(validLabel+validUntil, template.AlignRight(), template.FontSize(9),
					template.FontFamily("Inter"))
			})
		})
	})
}
