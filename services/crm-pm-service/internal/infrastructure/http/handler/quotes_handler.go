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

func (h *QuotesHandler) Create(c *gin.Context) {
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

	quoteReq := req.ToDomain()

	quote, err := h.service.CreateQuote(c.Request.Context(), id, quoteReq)
	if err != nil {
		mapDomainError(c, "failed_to_create_quote", err)
		return
	}

	response.Created(c, dto.FromQuoteDomain(quote))
}
