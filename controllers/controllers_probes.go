package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Readiness e Liveness probe do kubernetes
func HealthGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
	  "status": "UP",
	})
}
