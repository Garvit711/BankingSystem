package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userIDKey = "userID"

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		unauthorized := func() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}

		tokenString, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || tokenString == "" {
			unauthorized()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !token.Valid {
			unauthorized()
			return
		}

		sub, err := token.Claims.GetSubject()
		if err != nil {
			unauthorized()
			return
		}
		id, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			unauthorized()
			return
		}

		c.Set(userIDKey, id)
		c.Next()
	}
}

func UserID(c *gin.Context) int64 {
	v, _ := c.Get(userIDKey)
	id, _ := v.(int64)
	return id
}