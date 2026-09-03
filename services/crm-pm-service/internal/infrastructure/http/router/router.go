package router

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/handler"
	"github.com/gin-gonic/gin"
)

// ! Router nunca importa service ni repository directamente, solo conoce el tipo *handler.(Handlers)

// Agrupa todos los handlers del servicio. Cuando se agregue un recurso nuevo
// (contracts, deals, etc) solo se agrega un campo aqui
type Handlers struct {
	Health *handler.HealthHandler
	Leads  *handler.LeadsHandler
	// Contracts *handler.ContractsHandler	ejemplo
}

// Setup registra todas las rutas del servicio sobre el engine ya configurado
// NO construye nada: solo cablea rutas -> metodos del handler
func Setup(r *gin.Engine, h *Handlers) {
	r.GET("/health", h.Health.Check)

	api := r.Group("api/v1")
	{
		registerLeadsRoutes(api, h.Leads)
		// registerContractsRoutes(api, h.Contracts)	ejemplo
	}
}

func registerLeadsRoutes(rg *gin.RouterGroup, h *handler.LeadsHandler) {
	leads := rg.Group("/leads")
	{
		leads.POST("", h.Create)
		leads.GET("", h.GetByStatus)
		leads.GET("/:id", h.GetByID)
		leads.PATCH("/:id/status", h.UpdateStatus)
		leads.DELETE("/:id", h.Delete)
	}
}
