package handlers

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetUserProfileHandler(c *gin.Context)
}
