package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetUserProfileHandler(c *gin.Context)
	AddFollowHandler(c *gin.Context)
	DeleteFollowHandler(c *gin.Context)
}
