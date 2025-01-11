package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondWithError sends an error response
func RespondWithInternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   true,
		"message": message,
	})
}

func RespondWithBadRequestError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   true,
		"message": message,
	})
}

func RespondWithUnauthorizedError(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error":   true,
		"message": message,
	})
	c.Abort()
}

// RespondWithSuccess sends a success response
func RespondWithSuccessData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "success",
		"data":    data,
	})
}

func RespondWithSuccessMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": message,
	})
}

// ==========================================
// utils.RespondWithBadRequestError(c, err.Error())
// utils.RespondWithBadRequestError(c, constant.EMAIL_VALIDATION_FAILED)
// utils.RespondWithSuccessMsg(c, "OTP sent successfully")
