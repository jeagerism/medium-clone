package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/medium-clone/backend/internal/articles/services"
)

type articleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) ArticleHandler {
	return &articleHandler{articleService: articleService}
}

func (h *articleHandler) GetArticlesHandler(c *gin.Context) {
	articles, err := h.articleService.GetArticles()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}
