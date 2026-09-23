package domain

import (
	"errors"
)

var (
	ErrRepoTimeout           = errors.New("repository_operation_timed_out")
	ErrEmailAlreadyExists    = errors.New("email_already_exists")
	ErrInvalidData           = errors.New("invalid_data")
	ErrInvalidSource         = errors.New("invalid_source")
	ErrMissingRequiredFields = errors.New("missing_required_fields")
	ErrReqTimeout            = errors.New("request_timed_out, please_try_again")
	ErrNotFound              = errors.New("data_not_found")
	ErrMissingItems          = errors.New("quote_requires_at_least_one_item")
	ErrPaymethodRequired     = errors.New("payment_method_is_required")
	ErrInvalidQuoteItem      = errors.New("invalid_quote_item")
	ErrInvalidDuration       = errors.New("invalid_duration_days")
	ErrInvalidUnitPrice      = errors.New("unit_price_cannot_have_more_than_2_decimals")
	ErrInternal              = errors.New("internal_error")
)
