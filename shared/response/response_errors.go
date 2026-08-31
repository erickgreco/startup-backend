package response

import "errors"

var (
	ErrInvalidReqBody   = errors.New("invalid_request_body")
	ErrValidationFailed = errors.New("validation_failed")
)
