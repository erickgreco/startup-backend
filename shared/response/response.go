package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// Creado para no filtrar errores especificos
var (
	errInternal = errors.New("internal_server_error")
)

func Error(c *gin.Context, status int, err error) {
	c.JSON(status, Response{Error: err.Error()})
}

func BadRequest(c *gin.Context, err error) {
	Error(c, http.StatusConflict, err)
}

// ! No recibe error para evitar filtrar detalles internos
func InternalServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, errInternal)
}

func GatewayTimeout(c *gin.Context, err error) {
	Error(c, http.StatusGatewayTimeout, err)
}

func Conflict(c *gin.Context, err error) {
	Error(c, http.StatusConflict, err)
}

func NotFound(c *gin.Context, err error) {
	Error(c, http.StatusNotFound, err)
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Data: data})
}
