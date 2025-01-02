package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/medium-clone/backend/config"
	"github.com/jeagerism/medium-clone/backend/internal/articles/handlers"
	"github.com/jeagerism/medium-clone/backend/internal/articles/repositories"
	"github.com/jeagerism/medium-clone/backend/internal/articles/services"
	"github.com/jeagerism/medium-clone/backend/pkg/database"
)

type ginServer struct {
	app *gin.Engine
	db  database.Database
	cfg *config.Config
}

func NewGInServer(cfg *config.Config, db database.Database) Server {
	gin.SetMode(gin.ReleaseMode)
	ginApp := gin.New()
	ginApp.Use(gin.Recovery())
	ginApp.Use(gin.Logger())
	return &ginServer{
		app: ginApp,
		db:  db,
		cfg: cfg,
	}
}

func (s *ginServer) Start() {
	s.app.GET("/v1/health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})

	s.articleRoutes()

	serverURL := fmt.Sprintf(":%d", s.cfg.Server.Port)
	if err := s.app.Run(serverURL); err != nil {
		log.Fatalf("Failed to start server: %v", err) // ใช้ log ของ Go
	}
}

func (s *ginServer) articleRoutes() {
	arcRepo := repositories.NewArticleRepository(s.db.GetDb())
	arcServ := services.NewArticleService(arcRepo)
	arcHand := handlers.NewArticleHandler(arcServ)

	routes := s.app.Group("/articles")
	{
		routes.GET("", arcHand.GetArticlesHandler)
		routes.GET("/:id", arcHand.GetArticleByIDHandler)
	}
}
