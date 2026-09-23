package handler

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/dto"
	"github.com/erickgreco/startup-backend/shared/response"
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

func (h *LeadsHandler) Create(c *gin.Context) {
	var req dto.CreateLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, ErrInvalidReqBody)
		return
	}

	leadReq := req.ToDomain()
	// TODO: por ahora el source es hardcodeado, debera cambiar
	leadReq.Source = domain.LeadSrcManual

	lead, err := h.service.CreateLead(c.Request.Context(), leadReq)
	if err != nil {
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

func (h *LeadsHandler) UpdateStatus(c *gin.Context) {
	var req dto.UpdateLeadStatusRequest

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidLeadID)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, ErrInvalidReqBody)
		return
	}

	status := domain.LeadStatus(req.Status)

	lead, err := h.service.UpdateLeadStatus(c.Request.Context(), id, status)
	if err != nil {
		mapDomainError(c, "failed_to_update_lead_status", err)
		return
	}

	response.Ok(c, dto.FromLeadDomain(lead))
}

func (h *LeadsHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidLeadID)
		return
	}

	if err := h.service.DeleteLead(c.Request.Context(), id); err != nil {
		mapDomainError(c, "failed_to_delete_lead", err)
		return
	}

	response.NoContent(c)
}
