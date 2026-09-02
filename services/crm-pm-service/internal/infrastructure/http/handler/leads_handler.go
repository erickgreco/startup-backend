package handler

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/dto"
	"github.com/erickgreco/startup-backend/shared/response"
	"github.com/erickgreco/startup-backend/shared/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		mapDomainError(c, "failed_to_create_lead", err)
		return
	}

	response.Created(c, dto.FromLeadDomain(lead))
}

func (h *LeadsHandler) GetByStatus(c *gin.Context) {
	statusParam := c.Query("status")
	if statusParam == "" {
		response.BadRequest(c, ErrMissingStatusParam)
		return
	}

	status := domain.LeadStatus(statusParam)

	leads, err := h.service.GetLeadsByStatus(c.Request.Context(), status)
	if err != nil {
		mapDomainError(c, "failed_to_get_lead_by_status", err)
		return
	}

	response.Ok(c, dto.FromLeadDomainList(leads))
}

func (h *LeadsHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		response.BadRequest(c, ErrMissingIDParam)
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidLeadID)
		return
	}

	lead, err := h.service.GetLeadByID(c.Request.Context(), id)
	if err != nil {
		mapDomainError(c, "failed_to_get_lead_by_id", err)
		return
	}

	response.Ok(c, dto.FromLeadDomain(lead))
}
