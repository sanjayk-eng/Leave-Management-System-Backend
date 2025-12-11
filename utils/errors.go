package utils

import (
	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func CustomErr(c *gin.Context, code int, message string) error {
	return &AppError{Code: code, Message: message}

}
