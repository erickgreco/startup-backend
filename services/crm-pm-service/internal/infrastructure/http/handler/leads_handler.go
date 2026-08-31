package handler

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/dto"
	"github.com/erickgreco/startup-backend/shared/response"
	"github.com/erickgreco/startup-backend/shared/validator"
	"github.com/gin-gonic/gin"
)

// ! Handler nunca importa service ni repository directamente, solo conoce el tipo domain.(Service)

type LeadsHandler struct {
	service domain.LeadService
}

func NewLeadsHandler(service domain.LeadService) *LeadsHandler {
	return &LeadsHandler{
		service: service,
	}
}

func (h *LeadsHandler) Create(c *gin.Context) {
	var req dto.CreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, response.ErrInvalidReqBody)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.BadRequest(c, response.ErrValidationFailed)
		return
	}

	lead := req.ToDomain()

	if err := h.service.CreateLead(c.Request.Context(), lead); err != nil {
		// TODO: mapear dependiendo el error para evitar filtrar info
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, dto.FromLeadDomain(lead))
}
