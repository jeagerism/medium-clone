package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/medium-clone/backend/internal/articles/services"
	"github.com/jeagerism/medium-clone/backend/internal/entities"
)

type articleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(articleService services.ArticleService) ArticleHandler {
	return &articleHandler{articleService: articleService}
}

func (h *articleHandler) GetArticlesHandler(c *gin.Context) {
	var req entities.GetArticlesParams
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Call the service to fetch articles
	articlesResponse, err := h.articleService.GetArticles(req) // Adjust service method to accept search and tags
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the articles in JSON format
	c.JSON(http.StatusOK, articlesResponse)
}

func (h *articleHandler) GetArticleByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}
