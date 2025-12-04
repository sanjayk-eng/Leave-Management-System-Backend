package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

type DesignationInput struct {
	DesignationName string  `json:"designation_name" validate:"required"`
	Description     *string `json:"description,omitempty"`
}

// CreateDesignation - POST /api/designations
// Only ADMIN, SUPERADMIN, and HR can create designations
func (h *HandlerFunc) CreateDesignation(c *gin.Context) {
	// 1️⃣ Permission check
	role := c.GetString("role")
	if role != "SUPERADMIN" && role != "ADMIN" && role != "HR" {
		utils.RespondWithError(c, http.StatusForbidden, "only ADMIN, SUPERADMIN, and HR can create designations")
		return
	}

	// 2️⃣ Bind input JSON
	var input DesignationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid input: "+err.Error())
		return
	}

	// 3️⃣ Create designation
	designationID, err := h.Query.CreateDesignation(input.DesignationName, input.Description)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to create designation: "+err.Error())
		return
	}

	// 4️⃣ Response
	c.JSON(http.StatusCreated, gin.H{
		"message":        "designation created successfully",
		"designation_id": designationID,
	})
}

// GetAllDesignations - GET /api/designations
// All authenticated users can view designations
func (h *HandlerFunc) GetAllDesignations(c *gin.Context) {
	// 1️⃣ Fetch all designations
	designations, err := h.Query.GetAllDesignations()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to fetch designations: "+err.Error())
		return
	}

	// 2️⃣ Response
	c.JSON(http.StatusOK, gin.H{
		"message":      "designations fetched successfully",
		"designations": designations,
	})
}

// GetDesignationByID - GET /api/designations/:id
// All authenticated users can view a specific designation
func (h *HandlerFunc) GetDesignationByID(c *gin.Context) {
	// 1️⃣ Parse designation ID
	designationIDStr := c.Param("id")
	designationID, err := uuid.Parse(designationIDStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid designation ID")
		return
	}

	// 2️⃣ Fetch designation
	designation, err := h.Query.GetDesignationByID(designationID)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "designation not found")
		return
	}

	// 3️⃣ Response
	c.JSON(http.StatusOK, gin.H{
		"message":     "designation fetched successfully",
		"designation": designation,
	})
}

// UpdateDesignation - PATCH /api/designations/:id
// Only ADMIN, SUPERADMIN, and HR can update designations
func (h *HandlerFunc) UpdateDesignation(c *gin.Context) {
	// 1️⃣ Permission check
	role := c.GetString("role")
	if role != "SUPERADMIN" && role != "ADMIN" && role != "HR" {
		utils.RespondWithError(c, http.StatusForbidden, "only ADMIN, SUPERADMIN, and HR can update designations")
		return
	}

	// 2️⃣ Parse designation ID
	designationIDStr := c.Param("id")
	designationID, err := uuid.Parse(designationIDStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid designation ID")
		return
	}

	// 3️⃣ Bind input JSON
	var input DesignationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid input: "+err.Error())
		return
	}

	// 4️⃣ Update designation
	err = h.Query.UpdateDesignation(designationID, input.DesignationName, input.Description)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to update designation: "+err.Error())
		return
	}

	// 5️⃣ Response
	c.JSON(http.StatusOK, gin.H{
		"message":        "designation updated successfully",
		"designation_id": designationID,
	})
}

// DeleteDesignation - DELETE /api/designations/:id
// Only ADMIN, SUPERADMIN, and HR can delete designations
func (h *HandlerFunc) DeleteDesignation(c *gin.Context) {
	// 1️⃣ Permission check
	role := c.GetString("role")
	if role != "SUPERADMIN" && role != "ADMIN" && role != "HR" {
		utils.RespondWithError(c, http.StatusForbidden, "only ADMIN, SUPERADMIN, and HR can delete designations")
		return
	}

	// 2️⃣ Parse designation ID
	designationIDStr := c.Param("id")
	designationID, err := uuid.Parse(designationIDStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid designation ID")
		return
	}

	// 3️⃣ Delete designation (will set employee designation_id to NULL due to ON DELETE SET NULL)
	err = h.Query.DeleteDesignation(designationID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to delete designation: "+err.Error())
		return
	}

	// 4️⃣ Response
	c.JSON(http.StatusOK, gin.H{
		"message": "designation deleted successfully. Employee designation_id set to NULL.",
	})
}
