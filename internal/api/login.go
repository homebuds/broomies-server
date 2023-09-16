package api

import "github.com/gin-gonic/gin"

type LoginRequest struct {
	Email string `json:"email"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
}
