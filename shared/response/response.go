package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{Error: err.Error()})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Data: data})
}
