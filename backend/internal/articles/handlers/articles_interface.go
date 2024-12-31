package handlers

import "github.com/gin-gonic/gin"

type ArticleHandler interface {
	GetArticlesHandler(c *gin.Context)
}
