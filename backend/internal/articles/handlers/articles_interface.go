package handlers

import "github.com/gin-gonic/gin"

type ArticleHandler interface {
	GetArticlesHandler(c *gin.Context)
	GetArticleByIDHandler(c *gin.Context)
	AddArticleHandler(c *gin.Context)
	UpdateArticleHandler(c *gin.Context)
}
