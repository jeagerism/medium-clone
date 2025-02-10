package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetUserProfileHandler(c *gin.Context)
	AddFollowHandler(c *gin.Context)
	DeleteFollowHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
	RefreshTokenHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
}
