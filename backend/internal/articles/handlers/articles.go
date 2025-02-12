package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/medium-clone/backend/internal/articles/entities"
	"github.com/jeagerism/medium-clone/backend/internal/articles/services"
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
	var req entities.GetArticlesByUserIDParams
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	articles, err := h.articleService.GetArticleByUserID(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (h *articleHandler) GetArticleByUserIDHandler(c *gin.Context) {
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

func (h *articleHandler) AddArticleHandler(c *gin.Context) {
	// ดึง user_id จาก Context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// แปลง userIDValue ให้เป็น float64 ก่อน
	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// แปลง float64 เป็น int
	userID := int(userIDFloat)

	// ดึง role
	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: No role found"})
		return
	}

	// แปลง role เป็น string
	role, ok := roleValue.(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid role type"})
		return
	}

	if role != "user" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: Only users and admins can add articles"})
		return
	}

	var req entities.AddArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID // กำหนด userID ให้บทความ

	if err := h.articleService.AddArticle(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Article created successfully"})
}

func (h *articleHandler) UpdateArticleHandler(c *gin.Context) {
	// ดึง user_id
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// แปลง user_id เป็น int
	userID, ok := userIDValue.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var req entities.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := h.articleService.GetArticleByID(req.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// ตรวจสอบว่า user เป็นเจ้าของบทความหรือไม่
	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: You are not the owner of this article"})
		return
	}

	if err := h.articleService.UpdateArticle(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

func (h *articleHandler) DeleteArticleHandler(c *gin.Context) {
	// ดึง user_id จาก Context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// แปลง userIDValue ให้เป็น float64 ก่อน
	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// แปลง float64 เป็น int
	userID := int(userIDFloat)

	// ดึง role
	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: No role found"})
		return
	}

	// แปลง role เป็น string
	role, ok := roleValue.(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid role type"})
		return
	}

	// รับค่า request
	var req entities.DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ดึงข้อมูลบทความจากฐานข้อมูล
	article, err := h.articleService.GetArticleByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// ตรวจสอบสิทธิ์: เจ้าของบทความ หรือ Admin เท่านั้นที่สามารถลบได้
	if article.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: You can only delete your own articles or must be an admin"})
		return
	}

	// ลบบทความ
	if err := h.articleService.DeleteArticle(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (h *articleHandler) AddCommentHandler(c *gin.Context) {
	var req entities.AddCommentRequest
	// Bind and validate the JSON payload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Proceed with business logic if validation passes
	err := h.articleService.AddComment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment posted"})
}

func (h *articleHandler) DeleteCommentHandler(c *gin.Context) {
	// ดึง user_id จาก Context
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// แปลง userIDValue ให้เป็น float64 ก่อน
	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	// แปลง float64 เป็น int
	userID := int(userIDFloat)

	// ดึง role
	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: No role found"})
		return
	}

	// แปลง role เป็น string
	role, ok := roleValue.(string)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Invalid role type"})
		return
	}

	// รับค่า request
	var req entities.DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ดึงข้อมูลคอมเมนต์จากฐานข้อมูล
	comment, err := h.articleService.GetCommentByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// ตรวจสอบสิทธิ์: เจ้าของคอมเมนต์หรือต้องเป็น Admin เท่านั้น
	if comment.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized: You can only delete your own comments or must be an admin"})
		return
	}

	// ลบคอมเมนต์
	if err := h.articleService.DeleteComment(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (h *articleHandler) GetArticleCommentsHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	comments, err := h.articleService.GetArticleComments(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}
