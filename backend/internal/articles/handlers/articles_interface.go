package handlers

import "github.com/gin-gonic/gin"

type ArticleHandler interface {
	GetArticlesHandler(c *gin.Context)
	GetArticleByIDHandler(c *gin.Context)
	AddArticleHandler(c *gin.Context)
	UpdateArticleHandler(c *gin.Context)
	AddCommentHandler(c *gin.Context)
	DeleteArticleHandler(c *gin.Context)
	DeleteCommentHandler(c *gin.Context)
	GetArticleCommentsHandler(c *gin.Context)
}
