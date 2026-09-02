package handler

import (
	"errors"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/erickgreco/startup-backend/shared/response"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidReqBody   = errors.New("invalid_request_body")
	ErrValidationFailed = errors.New("validation_failed")
)

func mapCreateLeadError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		response.Conflict(c, err)
	case errors.Is(err, domain.ErrInvalidSource):
		response.BadRequest(c, err)
	case errors.Is(err, domain.ErrMissingRequiredFields):
		response.BadRequest(c, err)
	case errors.Is(err, domain.ErrRepoTimeout):
		response.GatewayTimeout(c, err)
	default:
		logger.Error("failed_to_create_lead", err)
		response.InternalServerError(c)
	}
}
