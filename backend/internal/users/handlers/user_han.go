package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jeagerism/medium-clone/backend/internal/users/entities"
	"github.com/jeagerism/medium-clone/backend/internal/users/services"
)

type userHandler struct {
	userServ services.UserService
}

func NewUserHandler(userServ services.UserService) UserHandler {
	return &userHandler{userServ: userServ}
}

func (h *userHandler) GetUserProfileHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	user, err := h.userServ.GetUserProfile(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) AddFollowHandler(c *gin.Context) {
	var req entities.UserAddFollowingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	if err := h.userServ.AddFollowing(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Following added"})
}

func (h *userHandler) DeleteFollowHandler(c *gin.Context) {
	var req entities.UserAddFollowingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request"})
		return
	}

	if err := h.userServ.DeleteFollowing(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Follow deleted"})
}

func (h *userHandler) LoginHandler(c *gin.Context) {
	var req entities.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// เรียก Service Layer
	response, err := h.userServ.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ส่ง Response
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) RegisterHandler(c *gin.Context) {
	var req entities.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userServ.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *userHandler) RefreshTokenHandler(c *gin.Context) {
	// รับค่า Refresh Token จาก Body หรือ Header
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is required"})
		return
	}

	// เรียก Service Layer เพื่อรีเฟรช Token
	token, err := h.userServ.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// ส่งคืน Token ใหม่
	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
