package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

// CompanySettings struct mapping the DB table
type CompanySettings struct {
	ID                   uuid.UUID `db:"id" json:"id"`
	WorkingDaysPerMonth  int       `db:"working_days_per_month" json:"working_days_per_month"`
	AllowManagerAddLeave bool      `db:"allow_manager_add_leave" json:"allow_manager_add_leave"`
	CreatedAt            string    `db:"created_at" json:"created_at"`
	UpdatedAt            string    `db:"updated_at" json:"updated_at"`
}

// GetCompanySettings - GET /api/settings/company
func (h *HandlerFunc) GetCompanySettings(c *gin.Context) {
	// Only SUPERADMIN and ADMIN allowed
	roleRaw, _ := c.Get("role")
	role := roleRaw.(string)
	if role != "SUPERADMIN" && role != "ADMIN" && role != "HR" {
		utils.RespondWithError(c, 403, "Not authorized to view settings")
		return
	}

	var settings CompanySettings
	err := h.Query.DB.Get(&settings, `SELECT * FROM Tbl_Company_Settings LIMIT 1`)
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to fetch settings: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"settings": settings,
	})
}

// UpdateCompanySettings - PUT /api/settings/company
func (h *HandlerFunc) UpdateCompanySettings(c *gin.Context) {
	// Only SUPERADMIN and ADMIN allowed
	roleRaw, _ := c.Get("role")
	role := roleRaw.(string)
	if role != "SUPERADMIN" && role != "ADMIN" && role != "HR" {
		utils.RespondWithError(c, 403, "Not authorized to update settings")
		return
	}

	var input struct {
		WorkingDaysPerMonth  int  `json:"working_days_per_month" binding:"required"`
		AllowManagerAddLeave bool `json:"allow_manager_add_leave"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, 400, "Invalid input: "+err.Error())
		return
	}

	_, err := h.Query.DB.Exec(`
        UPDATE Tbl_Company_Settings
        SET working_days_per_month=$1, allow_manager_add_leave=$2, updated_at=NOW()
    `, input.WorkingDaysPerMonth, input.AllowManagerAddLeave)

	if err != nil {
		utils.RespondWithError(c, 500, "Failed to update settings: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company settings updated successfully",
	})
}
