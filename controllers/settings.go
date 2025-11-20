package controllers

import "github.com/gin-gonic/gin"

// GetCompanySettings - GET /api/settings/company
func (s *HandlerFunc) GetCompanySettings(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get company settings"})
}

// UpdateCompanySettings - POST /api/settings/company
func (s *HandlerFunc) UpdateCompanySettings(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update company settings"})
}

// GetPermissions - GET /api/settings/permissions
func (s *HandlerFunc) GetPermissions(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get permissions"})
}

// UpdatePermissions - POST /api/settings/permissions
func (s *HandlerFunc) UpdatePermissions(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update permissions"})
}
