package http

import (
	"POSFlowBackend/internal/infrastructure/http/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	port   string
}

func NewServer(port string) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	return &Server{
		router: router,
		port:   port,
	}
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	log.Printf("🚀 Server starting on http://localhost:%s\n", s.port)
	return s.router.Run(":" + s.port)
}

func (s *Server) RegisterRoutes(registerFunc func(*gin.Engine)) {
	registerFunc(s.router)
}
