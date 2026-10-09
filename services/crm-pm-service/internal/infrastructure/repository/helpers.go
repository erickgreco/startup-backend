package repository

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/repository/sqlgen"
)

// * Convierte el modelo de dominio a los parametros esperados
// * por la query generada
func leadDomainToGen(lead *domain.LeadsModel) sqlgen.CreateLeadParams {
	return sqlgen.CreateLeadParams{
		ID:           lead.ID,
		Name:         lead.Name,
		BusinessName: lead.BusinessName,
		PhoneNumber:  lead.PhoneNumber,
		Email:        lead.Email,
		Message:      lead.Message,
		DataQuality:  string(lead.DataQuality),
		QualityScore: lead.QualityScore,
		Status:       string(lead.Status),
		Source:       string(lead.Source),
	}
}

// * Convierte el modelo generado por sqlc al modelo de dominio
// * Es el unico lugar donde sqlgen.Lead cruza hacia domain
func genToLeadDomain(row sqlgen.Lead) *domain.LeadsModel {
	return &domain.LeadsModel{
		ID:           row.ID,
		Name:         row.Name,
		BusinessName: row.BusinessName,
		PhoneNumber:  row.PhoneNumber,
		Email:        row.Email,
		Message:      row.Message,
		DataQuality:  domain.Quality(row.DataQuality),
		QualityScore: row.QualityScore,
		Status:       domain.LeadStatus(row.Status),
		Source:       domain.LeadSource(row.Source),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

func quoteDomainToGen(quote *domain.QuotesModel) sqlgen.CreateQuoteParams {
	return sqlgen.CreateQuoteParams{
		ID:                           quote.ID,
		LeadID:                       quote.LeadID,
		QuoteNumber:                  quote.QuoteNumber,
		ClientName:                   quote.ClientName,
		ClientBusiness:               quote.ClientBusiness,
		ClientEmail:                  quote.ClientEmail,
		ClientPhone:                  quote.ClientPhone,
		ProjectTitle:                 quote.ProjectTitle,
		ProjectDescription:           quote.ProjectDescription,
		ProjectEstimatedDurationDays: int32(quote.ProjectEstimatedDurationDays),
		PaymentMethod:                string(quote.PaymentMethod),
		Status:                       string(quote.Status),
		CreatedAt:                    quote.CreatedAt,
		Subtotal:                     quote.Subtotal.ToCents(),
		Total:                        quote.Total.ToCents(),
		Currency:                     quote.Currency.String(),
		Notes:                        quote.Notes,
	}
}

func quoteItemsDomainToGen(items []domain.QuoteItem) []sqlgen.CreateQuoteItemsParams {
	params := make([]sqlgen.CreateQuoteItemsParams, 0, len(items))
	for i, item := range items {
		params = append(params, sqlgen.CreateQuoteItemsParams{
			ID:          item.ID,
			QuoteID:     item.QuoteID,
			Description: item.Description,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			UnitPrice:   item.UnitPrice.ToCents(),
			Subtotal:    item.SubTotal.ToCents(),
			SortOrder:   int32(i),
			CreatedAt:   item.CreatedAt,
		})
	}

	return params
}

func genToQuoteDomain(row sqlgen.Quote, items []sqlgen.QuoteItem) (*domain.QuotesModel, error) {
	currency, err := domain.ParseCurrency(row.Currency)
	if err != nil {
		return nil, err
	}

	return &domain.QuotesModel{
		ID:                           row.ID,
		LeadID:                       row.LeadID,
		QuoteNumber:                  row.QuoteNumber,
		ClientName:                   row.ClientName,
		ClientBusiness:               row.ClientBusiness,
		ClientEmail:                  row.ClientEmail,
		ClientPhone:                  row.ClientPhone,
		ProjectTitle:                 row.ProjectTitle,
		ProjectDescription:           row.ProjectDescription,
		ProjectResume:                row.ProjectResume,
		ProjectEstimatedDurationDays: row.ProjectEstimatedDurationDays,
		FunctionalityExample:         row.FunctionalityExample,
		Items:                        genItemsToDomain(items),
		PaymentMethod:                domain.PayMethod(row.PaymentMethod),
		Status:                       domain.QuoteStatus(row.Status),
		CreatedAt:                    row.CreatedAt,
		SentAt:                       row.SentAt,
		AcceptedAt:                   row.AcceptedAt,
		RejectedAt:                   row.RejectedAt,
		Subtotal:                     domain.FromCents(row.Subtotal),
		Total:                        domain.FromCents(row.Total),
		Currency:                     currency,
		Notes:                        row.Notes,
		UpdatedAt:                    row.UpdatedAt,
	}, nil
}

func genItemsToDomain(rows []sqlgen.QuoteItem) []domain.QuoteItem {
	items := make([]domain.QuoteItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain.QuoteItem{
			ID:          row.ID,
			QuoteID:     row.QuoteID,
			Description: row.Description,
			Quantity:    row.Quantity,
			Unit:        row.Unit,
			UnitPrice:   domain.FromCents(row.UnitPrice),
			SubTotal:    domain.FromCents(row.Subtotal),
			SortOrder:   int(row.SortOrder),
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	return items
}

func quoteDomainToGenUpdate(quote *domain.QuotesModel) sqlgen.UpdateQuoteParams {
	return sqlgen.UpdateQuoteParams{
		ClientName:                   quote.ClientName,
		ClientBusiness:               quote.ClientBusiness,
		ClientEmail:                  quote.ClientEmail,
		ClientPhone:                  quote.ClientPhone,
		ProjectTitle:                 quote.ProjectTitle,
		ProjectDescription:           quote.ProjectDescription,
		ProjectResume:                quote.ProjectResume,
		ProjectEstimatedDurationDays: int32(quote.ProjectEstimatedDurationDays),
		FunctionalityExample:         quote.FunctionalityExample,
		PaymentMethod:                string(quote.PaymentMethod),
		Subtotal:                     quote.Subtotal.ToCents(),
		Total:                        quote.Total.ToCents(),
		Currency:                     quote.Currency.String(),
		Notes:                        quote.Notes,
		ID:                           quote.ID,
	}
}

func userProfileDomainToGen(up *domain.UserProfileModel) sqlgen.CreateUserProfileParams {
	return sqlgen.CreateUserProfileParams{
		ID:              up.ID,
		FirstName:       up.FirstName,
		LastName:        up.LastName,
		Email:           up.Email,
		PhoneNumber:     up.PhoneNumber,
		ProfileImageUrl: up.ProfileImageURL,
		BusinessName:    up.BusinessName,
		Industry:        up.Industry,
		LogoUrl:         up.LogoURL,
		DefaultCurrency: up.DefaultCurrency.String(),
		Timezone:        up.Timezone,
		Language:        up.Language,
		WebsiteUrl:      up.WebsiteURL,
		CreatedAt:       up.CreatedAt,
	}
}

func genToUserProfileDomain(row sqlgen.UserProfile) (*domain.UserProfileModel, error) {
	currency, err := domain.ParseCurrency(row.DefaultCurrency)
	if err != nil {
		return nil, err
	}

	return &domain.UserProfileModel{
		ID:              row.ID,
		FirstName:       row.FirstName,
		LastName:        row.LastName,
		Email:           row.Email,
		PhoneNumber:     row.PhoneNumber,
		ProfileImageURL: row.ProfileImageUrl,
		BusinessName:    row.BusinessName,
		Industry:        row.Industry,
		LogoURL:         row.LogoUrl,
		DefaultCurrency: currency,
		Timezone:        row.Timezone,
		Language:        row.Language,
		WebsiteURL:      row.WebsiteUrl,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}
