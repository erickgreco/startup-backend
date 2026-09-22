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
	Quotes *handler.QuotesHandler
	// Contracts *handler.ContractsHandler	ejemplo
}

// Setup registra todas las rutas del servicio sobre el engine ya configurado
// NO construye nada: solo cablea rutas -> metodos del handler
func Setup(r *gin.Engine, h *Handlers) {
	r.GET("/health", h.Health.Check)

	api := r.Group("api/v1")
	{
		registerLeadsRoutes(api, h.Leads)
		registerQuotesRoutes(api, h.Quotes)
		// registerContractsRoutes(api, h.Contracts)	ejemplo
	}
}

func registerLeadsRoutes(rg *gin.RouterGroup, lh *handler.LeadsHandler) {
	leads := rg.Group("/leads")
	{
		leads.POST("", lh.Create)
		leads.GET("", lh.GetByStatus)
		leads.GET("/:id", lh.GetByID)
		leads.PATCH("/:id/status", lh.UpdateStatus)
		leads.DELETE("/:id", lh.Delete)
	}
}

func registerQuotesRoutes(rg *gin.RouterGroup, qh *handler.QuotesHandler) {
	quotes := rg.Group("/quotes")
	{
		quotes.POST("/:leadID", qh.Create)
		quotes.GET("/:id", qh.GetByID)
	}
}
