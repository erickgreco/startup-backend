package handler

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/dto"
	"github.com/erickgreco/startup-backend/shared/logger"
	"github.com/erickgreco/startup-backend/shared/response"
	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	service domain.UserProfileService
}

func NewUserProfileHandler(service domain.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		service: service,
	}
}

func (h *UserProfileHandler) Create(c *gin.Context) {
	var req dto.UserProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("user_profile_handler", err)
		response.BadRequest(c, ErrInvalidReqBody)
		return
	}

	created, err := h.service.CreateUserProfile(c.Request.Context(), req.ToDomain())
	if err != nil {
		mapDomainError(c, "failed_to_create_user_profile", err)
		return
	}

	response.Created(c, dto.FromUserProfileDomain(created))
}
