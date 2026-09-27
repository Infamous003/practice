package handlers

import (
	"practice/internal/response"

	"github.com/gin-gonic/gin"
)

func HealthHandler(c *gin.Context) {
	response.OK(c, map[string]string{"status": "ok"})
}
