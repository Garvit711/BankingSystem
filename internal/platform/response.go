package platform

import "github.com/gin-gonic/gin"

func RespondJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func RespondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}
