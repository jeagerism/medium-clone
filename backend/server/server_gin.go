package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/medium-clone/backend/config"
	article_han "github.com/jeagerism/medium-clone/backend/internal/articles/handlers"
	article_repo "github.com/jeagerism/medium-clone/backend/internal/articles/repositories"
	article_svc "github.com/jeagerism/medium-clone/backend/internal/articles/services"

	user_han "github.com/jeagerism/medium-clone/backend/internal/users/handlers"
	user_repo "github.com/jeagerism/medium-clone/backend/internal/users/repositories"
	user_svc "github.com/jeagerism/medium-clone/backend/internal/users/services"

	"github.com/jeagerism/medium-clone/backend/pkg/database"
)

type ginServer struct {
	app *gin.Engine
	db  database.Database
	cfg *config.Config
}

func NewGinServer(cfg *config.Config, db database.Database) Server {
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
	s.userRoutes()

	serverURL := fmt.Sprintf(":%d", s.cfg.Server().GetPort())
	if err := s.app.Run(serverURL); err != nil {
		log.Fatalf("Failed to start server: %v", err) // ใช้ log ของ Go
	}
}

func (s *ginServer) articleRoutes() {
	arcRepo := article_repo.NewArticleRepository(s.db.GetDb())
	arcServ := article_svc.NewArticleService(arcRepo)
	arcHand := article_han.NewArticleHandler(arcServ)

	routes := s.app.Group("/articles")
	{
		// Articles API
		routes.GET("", arcHand.GetArticlesHandler)
		routes.GET("/:id", arcHand.GetArticleByIDHandler)
		routes.POST("", arcHand.AddArticleHandler)
		routes.PUT("", arcHand.UpdateArticleHandler)
		routes.DELETE("", arcHand.DeleteArticleHandler)
		routes.GET("list", arcHand.GetArticleByUserIDHandler)

		// Comment API
		routes.POST("/comment", arcHand.AddCommentHandler)
		routes.DELETE("/comment", arcHand.DeleteCommentHandler)
		routes.GET("/:id/comments", arcHand.GetArticleCommentsHandler)
	}
}

func (s *ginServer) userRoutes() {
	userRepo := user_repo.NewUserRepository(s.db.GetDb())
	userServ := user_svc.NewUserService(userRepo, s.cfg) // ส่ง s.cfg เข้าไปใน NewUserService
	userHand := user_han.NewUserHandler(userServ)

	routes := s.app.Group("/user")
	{
		routes.GET("/@:id", userHand.GetUserProfileHandler)
		routes.POST("/following", userHand.AddFollowHandler)
		routes.DELETE("/following", userHand.DeleteFollowHandler)
		routes.POST("/login", userHand.LoginHandler)
		routes.POST("/register", userHand.RegisterHandler)
		routes.POST("/retoken", userHand.RefreshTokenHandler)

	}
}
