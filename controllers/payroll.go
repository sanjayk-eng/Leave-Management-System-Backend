package controllers

import "github.com/gin-gonic/gin"

// RunPayroll - POST /api/payroll/run
func (s *HandlerFunc) RunPayroll(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Run payroll"})
}

// FinalizePayroll - POST /api/payroll/:id/finalize
func (s *HandlerFunc) FinalizePayroll(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Finalize payroll"})
}

// GetPayslipPDF - GET /api/payslips/:id/pdf
func (s *HandlerFunc) GetPayslipPDF(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get payslip PDF"})
}
