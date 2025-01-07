package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StatusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Server is up and running"})
}
