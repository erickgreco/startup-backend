package handler

import (
	"errors"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/erickgreco/startup-backend/shared/response"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidReqBody     = errors.New("invalid_request_body")
	ErrValidationFailed   = errors.New("validation_failed")
	ErrMissingStatusParam = errors.New("status_query_parameter_is_required")
)

func mapDomainError(c *gin.Context, logMsg string, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		response.Conflict(c, err)
	case errors.Is(err, domain.ErrInvalidSource), errors.Is(err, domain.ErrMissingRequiredFields):
		response.BadRequest(c, err)
	case errors.Is(err, domain.ErrRepoTimeout):
		response.GatewayTimeout(c, err)
	case errors.Is(err, domain.ErrNotFound):
		response.NotFound(c, err)
	default:
		logger.Error(logMsg, err)
		response.InternalServerError(c)
	}
}
