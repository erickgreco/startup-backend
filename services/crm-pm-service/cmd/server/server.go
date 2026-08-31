package server

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/cmd/config"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/handler"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/http/router"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/infrastructure/repository"
	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// const version string = "0.0.1"

type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
	cfg    *config.Config
}

func New(db *pgxpool.Pool, cfg *config.Config) *Server {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	s := &Server{
		router: r,
		db:     db,
		cfg:    cfg,
	}

	handlers := s.buildHandlers()
	router.Setup(s.router, handlers)

	return s
}

// Ensambla cada cadena repository -> service -> handler
// Cada capa se asigna a su tipo interfaz (domain.LeadRepository, domain.LeadService), nunca
// al struct concreto. Esto es lo unico que permite que handler y service dependan de
// abstracciones, no de implementaciones (Dependency Inversion)
func (s *Server) buildHandlers() *router.Handlers {
	var leadsRepo domain.LeadRepository = repository.NewLeadRepository(s.db)
	var leadsSrv domain.LeadService = service.NewLeadService(leadsRepo)
	leadsHandler := handler.NewLeadsHandler(leadsSrv)

	return &router.Handlers{
		Leads: leadsHandler,
	}
}

func (s *Server) Run() error {
	return s.router.Run(s.cfg.Port)
}
