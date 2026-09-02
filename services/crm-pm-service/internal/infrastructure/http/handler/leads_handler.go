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

// TODO: por ahora el source es hardcodeado, debera cambiar
func (h *LeadsHandler) Create(c *gin.Context) {
	var req dto.CreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, ErrInvalidReqBody)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		response.BadRequest(c, ErrValidationFailed)
		return
	}

	lead := req.ToDomain()
	lead.Source = domain.LeadSrcManual

	if err := h.service.CreateLead(c.Request.Context(), lead); err != nil {
		mapCreateLeadError(c, err)
		return
	}

	response.Created(c, dto.FromLeadDomain(lead))
}
