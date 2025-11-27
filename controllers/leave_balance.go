package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

// GetLeaveBalances - GET /api/employees/:id/leave-balances
// GetLeaveBalances - GET /api/employees/:id/leave-balances
func (s *HandlerFunc) GetLeaveBalances(c *gin.Context) {
	// 1. Parse employee ID from path
	employeeIDParam := c.Param("id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	// 2. Role check (employees can only view their own balances)
	roleRaw, _ := c.Get("role")
	role := roleRaw.(string)
	userIDRaw, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDRaw.(string))

	if role == "EMPLOYEE" && userID != employeeID {
		utils.RespondWithError(c, http.StatusForbidden, "Employees can only view their own balances")
		return
	}

	// 3. Query leave balances
	type Balance struct {
		LeaveType string  `db:"leave_type" json:"leave_type"`
		Used      float64 `db:"used" json:"used"`
		Total     float64 `db:"total" json:"total"`
		Available float64 `db:"available" json:"available"`
	}

	var balances []Balance

	query := `
		SELECT 
			lt.name AS leave_type,
			COALESCE(b.used, 0) AS used,
			lt.default_entitlement AS total,
			COALESCE(b.closing, lt.default_entitlement) AS available
		FROM Tbl_Leave_Type lt
		LEFT JOIN Tbl_Leave_balance b 
			ON lt.id = b.leave_type_id AND b.employee_id = $1
		ORDER BY lt.id
	`

	// 4. Prepare the statement explicitly
	stmt, err := s.Query.DB.Preparex(query)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to prepare statement: "+err.Error())
		return
	}
	defer stmt.Close()

	if err := stmt.Select(&balances, employeeID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch leave balances: "+err.Error())
		return
	}

	// 5. Send response
	c.JSON(http.StatusOK, gin.H{
		"employee_id": employeeID,
		"balances":    balances,
	})
}

// AdjustLeaveBalance - POST /api/leave-balances/:id/adjust
// AdjustLeaveBalance - POST /api/leave-balances/adjust
// AdjustLeaveBalance - POST /api/leave-balances/:id/adjust
func (s *HandlerFunc) AdjustLeaveBalance(c *gin.Context) {
	// 1️ Role check - Only ADMIN/HR allowed
	roleRaw, _ := c.Get("role")
	role := roleRaw.(string)
	if role != "ADMIN" && role != "SUPERADMIN" {
		utils.RespondWithError(c, 403, "Not authorized to adjust leave balances")
		return
	}

	// 2️ Get employee ID from params
	employeeIDParam := c.Param("id")
	employeeID, err := uuid.Parse(employeeIDParam)
	if err != nil {
		utils.RespondWithError(c, 400, "Invalid employee ID")
		return
	}

	// 3️ Parse JSON input
	var input struct {
		LeaveTypeID int     `json:"leave_type_id" validate:"required"`
		Quantity    float64 `json:"quantity" validate:"required"` // +ve or -ve
		Reason      string  `json:"reason" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, 400, "Invalid input: "+err.Error())
		return
	}

	currentYear := time.Now().Year()

	// 4️ Start transaction
	tx, err := s.Query.DB.Beginx()
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	// 5️ Fetch or create leave balance
	var balance struct {
		ID          uuid.UUID `db:"id"`
		Opening     float64   `db:"opening"`
		Accrued     float64   `db:"accrued"`
		Used        float64   `db:"used"`
		Adjusted    float64   `db:"adjusted"`
		Closing     float64   `db:"closing"`
		EmployeeID  uuid.UUID `db:"employee_id"`
		LeaveTypeID int       `db:"leave_type_id"`
		Year        int       `db:"year"`
	}

	// ✔ FIX: Only fetch the fields your struct has
	err = tx.Get(&balance, `
        SELECT 
            id,
            opening,
            accrued,
            used,
            adjusted,
            closing,
            employee_id,
            leave_type_id,
            year
        FROM Tbl_Leave_balance
        WHERE employee_id=$1 AND leave_type_id=$2 AND year=$3
        FOR UPDATE
    `, employeeID, input.LeaveTypeID, currentYear)

	if err == sql.ErrNoRows {
		// 5A: Fetch default entitlement
		var defaultEntitlement float64
		err = tx.Get(&defaultEntitlement, `SELECT default_entitlement FROM Tbl_Leave_Type WHERE id=$1`, input.LeaveTypeID)
		if err != nil {
			utils.RespondWithError(c, 500, "Failed to fetch leave type: "+err.Error())
			return
		}

		// 5B: Create balance row
		err = tx.QueryRow(`
            INSERT INTO Tbl_Leave_balance
            (employee_id, leave_type_id, year, opening, accrued, used, adjusted, closing, created_at, updated_at)
            VALUES ($1,$2,$3,$4,0,0,0,$4,NOW(),NOW())
            RETURNING id, opening, accrued, used, adjusted, closing
        `, employeeID, input.LeaveTypeID, currentYear, defaultEntitlement).
			Scan(&balance.ID, &balance.Opening, &balance.Accrued, &balance.Used, &balance.Adjusted, &balance.Closing)

		if err != nil {
			utils.RespondWithError(c, 500, "Failed to create leave balance: "+err.Error())
			return
		}

	} else if err != nil {
		utils.RespondWithError(c, 500, "Failed to fetch leave balance: "+err.Error())
		return
	}

	// 6️ Apply adjustment
	newAdjusted := balance.Adjusted + input.Quantity
	newClosing := balance.Opening + balance.Accrued - balance.Used + newAdjusted

	_, err = tx.Exec(`
        UPDATE Tbl_Leave_balance
        SET adjusted=$1, closing=$2, updated_at=NOW()
        WHERE id=$3
    `, newAdjusted, newClosing, balance.ID)

	if err != nil {
		utils.RespondWithError(c, 500, "Failed to update leave balance: "+err.Error())
		return
	}

	// 7️ Insert into adjustment log
	_, err = tx.Exec(`
        INSERT INTO Tbl_Leave_adjustment
        (employee_id, leave_type_id, quantity, reason, created_by, created_at, year)
        VALUES ($1,$2,$3,$4,$5,NOW(),$6)
    `, employeeID, input.LeaveTypeID, input.Quantity, input.Reason, c.GetString("user_id"), currentYear)

	if err != nil {
		utils.RespondWithError(c, 500, "Failed to record leave adjustment: "+err.Error())
		return
	}

	// 8️ Commit
	if err := tx.Commit(); err != nil {
		utils.RespondWithError(c, 500, "Transaction commit failed")
		return
	}

	c.JSON(200, gin.H{
		"message":      "Leave balance adjusted successfully",
		"new_adjusted": newAdjusted,
		"new_closing":  newClosing,
		"year":         currentYear,
	})
}
