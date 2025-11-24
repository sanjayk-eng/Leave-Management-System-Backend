package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

type UpdateRoleInput struct {
	Role string `json:"role" validate:"required"` // Only valid roles
}
type UpdateManagerInput struct {
	ManagerID string `json:"manager_id" validate:"required"` // UUID of new manager
}

// GetEmployees - GET /api/employees
func (h *HandlerFunc) GetEmployee(c *gin.Context) {
	role, _ := c.Get("role")
	r := role.(string)

	if r != "SUPERADMIN" && r != "Admin" {
		utils.RespondWithError(c, http.StatusUnauthorized, "not permitted")
		return
	}

	rows, err := h.Query.GetAllEmployees()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var employees []models.EmployeeInput

	for rows.Next() {
		var emp models.EmployeeInput
		if err := rows.Scan(
			&emp.ID, &emp.FullName, &emp.Email, &emp.Status,
			&emp.Role, &emp.Password, &emp.ManagerID,
			&emp.Salary, &emp.JoiningDate,
			&emp.CreatedAt, &emp.UpdatedAt, &emp.DeletedAt,
		); err != nil {
			utils.RespondWithError(c, 500, err.Error())
			return
		}
		employees = append(employees, emp)
	}

	c.JSON(200, gin.H{
		"message":   "Employees fetched",
		"employees": employees,
	})
}
func (h *HandlerFunc) GetEmployeeById(c *gin.Context) {

}

func (h *HandlerFunc) CreateEmployee(c *gin.Context) {
	role := c.GetString("role")
	if role != "SUPERADMIN" && role != "Admin" {
		utils.RespondWithError(c, http.StatusUnauthorized, "not permitted")
		return
	}

	var input models.EmployeeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !strings.HasSuffix(input.Email, "@zenithive.com") {
		utils.RespondWithError(c, 400, "email must end with @zenithive.com")
		return
	}

	// EMAIL EXIST CHECK
	exists, err := h.Query.CheckEmailExists(input.Email)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error())
		return
	}
	if exists {
		utils.RespondWithError(c, 400, "email already exists")
		return
	}

	// GET ROLE ID
	roleID, err := h.Query.GetRoleID(input.Role)
	if err != nil {
		utils.RespondWithError(c, 400, "role not found")
		return
	}

	// HASH PASSWORD
	hash, _ := utils.HashPassword(input.Password)

	// INSERT
	err = h.Query.InsertEmployee(
		input.FullName, input.Email,
		roleID, hash,
		input.Salary, input.JoiningDate,
	)
	if err != nil {
		utils.RespondWithError(c, 500, "failed to create employee")
		return
	}

	c.JSON(201, gin.H{"message": "employee created"})
}
func (h *HandlerFunc) UpdateEmployeeRole(c *gin.Context) {
	role := c.GetString("role")
	if role != "SUPERADMIN" && role != "ADMIN" {
		utils.RespondWithError(c, 401, "not permitted")
		return
	}

	empID := c.Param("id")
	var input UpdateRoleInput
	c.ShouldBindJSON(&input)

	currentRole, err := h.Query.GetEmployeeCurrentRole(empID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if currentRole == input.Role {
		c.JSON(200, gin.H{"message": "already same role"})
		return
	}

	updatedID, err := h.Query.UpdateEmployeeRole(empID, input.Role)
	if err != nil {
		utils.RespondWithError(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message":     "role updated",
		"employee_id": updatedID,
	})
}

func (h *HandlerFunc) UpdateEmployeeManager(c *gin.Context) {
	role := c.GetString("role")
	if role != "SUPERADMIN" && role != "ADMIN" && role != "HR" {
		utils.RespondWithError(c, 401, "not permitted")
		return
	}

	empID, _ := uuid.Parse(c.Param("id"))

	var input UpdateManagerInput
	c.ShouldBindJSON(&input)
	managerID, _ := uuid.Parse(input.ManagerID)

	if empID == managerID {
		utils.RespondWithError(c, 400, "cannot assign self")
		return
	}

	// CHECK MANAGER EXISTS
	exists, err := h.Query.ManagerExists(managerID)
	if err != nil || !exists {
		utils.RespondWithError(c, 404, "manager not found")
		return
	}

	// UPDATE MANAGER
	err = h.Query.UpdateManager(empID, managerID)
	if err != nil {
		utils.RespondWithError(c, 500, "failed")
		return
	}

	c.JSON(200, gin.H{
		"message":     "manager updated",
		"employee_id": empID,
		"manager_id":  managerID,
	})
}

// func (h *HandlerFunc) UpdateEmployeeInfo(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"message": "Employee info updated"})
// }

// GetEmployeeReports - GET /api/employees/:id/reports
func (s *HandlerFunc) GetEmployeeReports(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get employee reports"})
}
