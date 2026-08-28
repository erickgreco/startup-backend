package server

import (
	"github.com/erickgreco/startup-backend/services/crm-pm-service/cmd/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version string = "0.0.1"

type Server struct {
	router  *gin.Engine
	db      *pgxpool.Pool
	cfg     *config.Config
	version string
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

	return s
}

func (s *Server) Run() error {
	return s.router.Run(s.cfg.Port)
}
