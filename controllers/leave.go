package controllers

import "github.com/gin-gonic/gin"

// ApplyLeave - POST /api/leaves/apply
func (s *HandlerFunc) ApplyLeave(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Apply leave"})
}

// AdminAddLeave - POST /api/leaves/admin-add
func (s *HandlerFunc) AdminAddLeave(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Admin add leave"})
}

// ActionLeave - POST /api/leaves/:id/action
func (s *HandlerFunc) ActionLeave(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Approve/Reject leave"})
}
