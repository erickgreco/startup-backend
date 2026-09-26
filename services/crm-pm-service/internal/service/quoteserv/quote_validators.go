package quoteserv

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/shopspring/decimal"
)

func validateRequest(req *domain.QuoteRequest) error {
	switch {
	case req.ClientName == "", req.ClientEmail == "", req.ClientPhone == "":
		return domain.ErrMissingClientData
	case req.ProjectTitle == "", req.ProjectDescription == "":
		return domain.ErrMissingProjectData
	case req.ProjectEstimatedDurationDays < 1:
		return domain.ErrInvalidProjectDuration
	}

	if err := validateRequestItems(req.Items); err != nil {
		return err
	}

	return nil
}

func validateRequestItems(reqItems []domain.ItemRequest) error {
	switch {
	case len(reqItems) == 0:
		return domain.ErrMissingItems
	}

	for _, item := range reqItems {
		switch {
		case item.Quantity.LessThanOrEqual(decimal.Zero):
			return domain.ErrInvalidQuantity
		case item.UnitPrice.LessThanOrEqual(domain.ZeroMoney) || item.UnitPrice.HasSubCentPrecision():
			return domain.ErrInvalidUnitPrice
		}
	}

	return nil
}
