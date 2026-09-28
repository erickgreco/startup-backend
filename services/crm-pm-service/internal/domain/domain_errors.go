package domain

import (
	"errors"
)

var (
	ErrRepoTimeout            = errors.New("repository_operation_timed_out")
	ErrAlreadyExists          = errors.New("data_already_exists")
	ErrInvalidData            = errors.New("invalid_data")
	ErrInvalidSource          = errors.New("invalid_source")
	ErrMissingRequiredFields  = errors.New("missing_required_fields")
	ErrReqTimeout             = errors.New("request_timed_out, please_try_again")
	ErrNotFound               = errors.New("data_not_found")
	ErrMissingItems           = errors.New("quote_requires_at_least_one_item")
	ErrPaymethodRequired      = errors.New("payment_method_is_required")
	ErrInvalidQuoteItem       = errors.New("invalid_quote_item")
	ErrInvalidDuration        = errors.New("invalid_duration_days")
	ErrInvalidUnitPrice       = errors.New("unit_price_cannot_have_more_than_2_decimals")
	ErrInternal               = errors.New("internal_error")
	ErrMissingClientData      = errors.New("missing_required_client_data")
	ErrMissingProjectData     = errors.New("missing_required_project_data")
	ErrInvalidProjectDuration = errors.New("invalid_project_duration")
	ErrInvalidQuantity        = errors.New("item_quantity_required")
	ErrInvalidStatusForUpdate = errors.New("quote_status_must_be_draft")
	ErrUnableToDeleteQuote    = errors.New("unable_to_delete_quote")
)
