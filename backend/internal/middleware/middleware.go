package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UsersRole string

const (
	UserRole    UsersRole = "user"
	ManagerRole UsersRole = "manager"
	AdminRole   UsersRole = "admin"
)

// JwtAuthentication: ตรวจสอบ JWT Token จาก Header
func JwtAuthentication(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ดึง JWT Token จาก Authorization Header
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Token is missing"})
			c.Abort()
			return
		}

		// ถอดรหัส JWT Token
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			// ตรวจสอบ Signing Method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Invalid token"})
			c.Abort()
			return
		}

		// ตรวจสอบ Claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// บันทึกข้อมูลลงใน Context
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: Invalid claims"})
			c.Abort()
		}
	}
}

// Authorization: ตรวจสอบสิทธิ์การเข้าถึง API
func Authorization(allowedRoles ...UsersRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ดึงบทบาทของผู้ใช้จาก Context
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: No role found"})
			c.Abort()
			return
		}

		// ตรวจสอบว่าบทบาทของผู้ใช้อยู่ในกลุ่มที่อนุญาตหรือไม่
		for _, role := range allowedRoles {
			if userRole == string(role) {
				c.Next()
				return
			}
		}

		// หากบทบาทไม่ตรงกับที่อนุญาต
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: You do not have permission to access this resource"})
		c.Abort()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireRole(roles ...UsersRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: No role found"})
			c.Abort()
			return
		}

		// แปลง roleValue เป็น string
		role, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: Invalid role type"})
			c.Abort()
			return
		}

		// ตรวจสอบบทบาท
		for _, allowedRole := range roles {
			if role == string(allowedRole) {
				c.Next()
				return
			}
		}

		// หากบทบาทไม่ตรงกับที่อนุญาต
		c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: You do not have permission"})
		c.Abort()
	}
}
