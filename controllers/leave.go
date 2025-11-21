package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

// ApplyLeave - POST /api/leaves/apply
func (h *HandlerFunc) ApplyLeave(c *gin.Context) {
	// 1️⃣ Get employee ID and role from middleware
	employeeIDRaw, exists := c.Get("user_id")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Employee ID not found")
		return
	}

	employeeIDStr, ok := employeeIDRaw.(string)
	if !ok {
		utils.RespondWithError(c, http.StatusInternalServerError, "Invalid employee ID format")
		return
	}

	employeeID, err := uuid.Parse(employeeIDStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to parse employee ID: "+err.Error())
		return
	}

	roleRaw, exists := c.Get("role")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Role not found")
		return
	}
	role := roleRaw.(string)

	// Only employees can apply leave
	if role != "EMPLOYEE" {
		utils.RespondWithError(c, http.StatusForbidden, "Only employees can apply leave")
		return
	}

	// 2️⃣ Bind request body
	var input models.LeaveInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	// Assign middleware UUIDs
	input.EmployeeID = employeeID
	input.AppliedByID = &employeeID

	// 3️⃣ Calculate leave days if not provided
	if input.Days == nil {
		days := input.EndDate.Sub(input.StartDate).Hours()/24 + 1
		if days <= 0 {
			utils.RespondWithError(c, http.StatusBadRequest, "End date must be after start date")
			return
		}
		input.Days = &days
	}

	// 4️⃣ Check leave type exists
	var leaveType struct {
		ID                 int `db:"id"`
		DefaultEntitlement int `db:"default_entitlement"`
	}
	err = h.DB.Get(&leaveType, `
        SELECT id, default_entitlement 
        FROM Tbl_Leave_type 
        WHERE id=$1
    `, input.LeaveTypeID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(c, http.StatusBadRequest,
				fmt.Sprintf("Leave type with ID %d does not exist", input.LeaveTypeID))
			return
		}
		utils.RespondWithError(c, http.StatusInternalServerError,
			fmt.Sprintf("Failed to fetch leave type: %s", err.Error()))
		return
	}

	// 5️⃣ Check leave balance
	var balance struct {
		Closing float64 `db:"closing"`
	}
	err = h.DB.Get(&balance, `
        SELECT closing
        FROM Tbl_Leave_balance
        WHERE employee_id=$1 AND leave_type_id=$2 AND year=EXTRACT(YEAR FROM CURRENT_DATE)
    `, employeeID, input.LeaveTypeID)

	if err != nil {
		// First-time leave: create leave balance record
		_, err2 := h.DB.Exec(`
            INSERT INTO Tbl_Leave_balance 
                (employee_id, leave_type_id, year, opening, accrued, used, adjusted, closing, created_at)
            VALUES ($1, $2, EXTRACT(YEAR FROM CURRENT_DATE), $3, 0, 0, 0, $3, NOW())
        `, employeeID, input.LeaveTypeID, leaveType.DefaultEntitlement)
		if err2 != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create leave balance: "+err2.Error())
			return
		}
		balance.Closing = float64(leaveType.DefaultEntitlement)
	}

	// Check if employee has enough leave
	if balance.Closing < *input.Days {
		utils.RespondWithError(c, http.StatusBadRequest, "Insufficient leave balance")
		return
	}

	// 6️⃣ Check overlapping leaves
	var overlapCount int
	err = h.DB.Get(&overlapCount, `
        SELECT COUNT(1)
        FROM Tbl_Leave
        WHERE employee_id=$1
        AND status IN ('Pending','Approved')
        AND start_date <= $2 AND end_date >= $3
    `, employeeID, input.EndDate, input.StartDate)

	if err != nil || overlapCount > 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "Overlapping leave exists")
		return
	}

	// 7️⃣ Insert leave request
	var leaveID uuid.UUID
	status := "Pending"
	err = h.DB.QueryRow(`
        INSERT INTO Tbl_Leave 
            (employee_id, leave_type_id, start_date, end_date, days, status, applied_by, created_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
        RETURNING id
    `, employeeID, input.LeaveTypeID, input.StartDate, input.EndDate, *input.Days, status, employeeID).Scan(&leaveID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to apply leave: "+err.Error())
		return
	}

	// 8️⃣ Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message":  "Leave applied successfully",
		"leave_id": leaveID,
		"status":   status,
		"days":     *input.Days,
	})
}

// AdminAddLeave - POST /api/leaves/admin-add
func (s *HandlerFunc) AdminAddLeave(c *gin.Context) {
	roleValue, exists := c.Get("role")
	if !exists {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to get role")
		return
	}
	userRole := roleValue.(string)
	if userRole != "SUPERADMIN" {
		utils.RespondWithError(c, http.StatusUnauthorized, "not permitted to assign manager")
		return
	}
	var input models.LeaveTypeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid input: "+err.Error())
		return
	}

	// Set defaults
	if input.IsPaid == nil {
		defaultPaid := false
		input.IsPaid = &defaultPaid
	}
	if input.DefaultEntitlement == nil {
		defaultEntitlement := 0
		input.DefaultEntitlement = &defaultEntitlement
	}
	if input.LeaveCount == nil {
		defaultCount := 2
		input.LeaveCount = &defaultCount
	}

	if *input.LeaveCount <= 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "leave_count must be greater than 0")
		return
	}

	query := `
		INSERT INTO Tbl_Leave_type (name, is_paid, default_entitlement)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	var leave models.LeaveType
	err := s.DB.QueryRow(query, input.Name, *input.IsPaid, *input.DefaultEntitlement).
		Scan(&leave.ID, &leave.CreatedAt, &leave.UpdatedAt)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to insert leave type: "+err.Error())
		return
	}

	leave.Name = input.Name
	leave.IsPaid = *input.IsPaid
	leave.DefaultEntitlement = *input.DefaultEntitlement

	c.JSON(http.StatusOK, leave)
}

func (s *HandlerFunc) GetAllLeavePolicies(c *gin.Context) {
	var leaves []models.LeaveType

	query := `SELECT id, name, is_paid, default_entitlement,  created_at, updated_at FROM Tbl_Leave_type ORDER BY id`
	err := s.DB.Select(&leaves, query)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch leave types: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, leaves) // send models directly
}

// ActionLeave - POST /api/leaves/:id/action
func (s *HandlerFunc) ActionLeave(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Approve/Reject leave"})
}
