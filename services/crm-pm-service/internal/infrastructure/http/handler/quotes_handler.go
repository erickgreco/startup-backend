package handler

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/dto"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/erickgreco/startup-backend/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ! Handler nunca importa service ni repository directamente, solo conoce el tipo domain.(Service)

type QuotesHandler struct {
	service domain.QuoteService
}

func NewQuotesHandler(service domain.QuoteService) *QuotesHandler {
	return &QuotesHandler{
		service: service,
	}
}

func (h *QuotesHandler) ManualCreate(c *gin.Context) {
	var req dto.CreateQuoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("quote_handler_err", err)
		response.BadRequest(c, ErrInvalidReqBody)
		return
	}

	quoteReq, err := req.ToDomain()
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	created, err := h.service.CreateManualQuote(c.Request.Context(), quoteReq)
	if err != nil {
		mapDomainError(c, "failed_to_create_quote", err)
		return
	}

	response.Created(c, dto.FromQuoteDomain(created))
}

func (h *QuotesHandler) CreateFromLead(c *gin.Context) {
	var req dto.CreateQuoteRequest

	idParam := c.Param("leadID")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidLeadID)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("quote_handler_err", err)
		response.BadRequest(c, ErrInvalidReqBody)
		return
	}

	quoteReq, err := req.ToDomain()
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	quote, err := h.service.CreateQuoteFromLead(c.Request.Context(), id, quoteReq)
	if err != nil {
		mapDomainError(c, "failed_to_create_quote", err)
		return
	}

	response.Created(c, dto.FromQuoteDomain(quote))
}

func (h *QuotesHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidQuoteID)
		return
	}

	quote, err := h.service.GetQuoteByID(c.Request.Context(), id)
	if err != nil {
		mapDomainError(c, "failed_to_get_quote_by_id", err)
		return
	}

	response.Ok(c, dto.FromQuoteDomain(quote))
}

func (h *QuotesHandler) UpdateDraft(c *gin.Context) {
	var req dto.CreateQuoteRequest

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidQuoteID)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("quote_handler_err", err)
		return
	}

	quoteReq, err := req.ToDomain()
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	updated, err := h.service.UpdateQuote(c.Request.Context(), id, quoteReq)
	if err != nil {
		mapDomainError(c, "failed_to_update_quote", err)
		return
	}

	response.Ok(c, dto.FromQuoteDomain(updated))
}

func (h *QuotesHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, ErrInvalidQuoteID)
		return
	}

	if err := h.service.DeleteQuote(c.Request.Context(), id); err != nil {
		mapDomainError(c, "failed_to_delete_lead", err)
		return
	}

	response.NoContent(c)
}
