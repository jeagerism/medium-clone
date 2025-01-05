package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
