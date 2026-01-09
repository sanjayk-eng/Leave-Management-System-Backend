package repositories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// LeaveTypeData represents leave type information for balance calculation
type LeaveTypeData struct {
	LeaveTypeID       int     `db:"leave_type_id"`
	LeaveTypeName     string  `db:"leave_type_name"`
	DefaultEntitlement float64 `db:"default_entitlement"`
}

// BalanceData represents raw balance data from database
type BalanceData struct {
	LeaveTypeID int     `db:"leave_type_id"`
	Opening     float64 `db:"opening"`
	Accrued     float64 `db:"accrued"`
	Used        float64 `db:"used"`
	Adjusted    float64 `db:"adjusted"`
	Closing     float64 `db:"closing"`
}

// LeaveBalanceForAdjustment represents leave balance structure for adjustment operations
type LeaveBalanceForAdjustment struct {
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

// GetAllLeaveTypesWithEntitlements fetches all leave types with their default entitlements
func (r *Repository) GetAllLeaveTypesWithEntitlements() ([]LeaveTypeData, error) {
	var leaveTypes []LeaveTypeData
	query := `
		SELECT 
			lt.id AS leave_type_id,
			lt.name AS leave_type_name,
			COALESCE(lt.default_entitlement, 0) AS default_entitlement
		FROM Tbl_Leave_Type lt
		ORDER BY lt.id
	`
	err := r.DB.Select(&leaveTypes, query)
	return leaveTypes, err
}

// GetLeaveBalancesByEmployeeAndYear fetches leave balances for a specific employee and year
func (r *Repository) GetLeaveBalancesByEmployeeAndYear(employeeID uuid.UUID, year int) ([]BalanceData, error) {
	var balanceRecords []BalanceData
	query := `
		SELECT 
			leave_type_id,
			COALESCE(opening, 0) AS opening,
			COALESCE(accrued, 0) AS accrued,
			COALESCE(used, 0) AS used,
			COALESCE(adjusted, 0) AS adjusted,
			COALESCE(closing, 0) AS closing
		FROM Tbl_Leave_balance
		WHERE employee_id = $1 AND year = $2
	`
	err := r.DB.Select(&balanceRecords, query, employeeID, year)
	return balanceRecords, err
}

// GetLeaveBalanceForAdjustment fetches leave balance for adjustment with FOR UPDATE lock
func (r *Repository) GetLeaveBalanceForAdjustment(tx *sqlx.Tx, employeeID uuid.UUID, leaveTypeID int, year int) (LeaveBalanceForAdjustment, error) {
	var balance LeaveBalanceForAdjustment
	query := `
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
	`
	err := tx.Get(&balance, query, employeeID, leaveTypeID, year)
	return balance, err
}

// GetDefaultEntitlementByLeaveTypeID fetches default entitlement for a leave type
func (r *Repository) GetDefaultEntitlementByLeaveTypeID(tx *sqlx.Tx, leaveTypeID int) (float64, error) {
	var defaultEntitlement float64
	err := tx.Get(&defaultEntitlement, `SELECT default_entitlement FROM Tbl_Leave_Type WHERE id=$1`, leaveTypeID)
	return defaultEntitlement, err
}

// CreateLeaveBalanceForAdjustment creates a new leave balance record
func (r *Repository) CreateLeaveBalanceForAdjustment(tx *sqlx.Tx, employeeID uuid.UUID, leaveTypeID int, year int, defaultEntitlement float64) (LeaveBalanceForAdjustment, error) {
	var balance LeaveBalanceForAdjustment
	err := tx.QueryRow(`
		INSERT INTO Tbl_Leave_balance
		(employee_id, leave_type_id, year, opening, accrued, used, adjusted, closing, created_at, updated_at)
		VALUES ($1,$2,$3,$4,0,0,0,$4,NOW(),NOW())
		RETURNING id, opening, accrued, used, adjusted, closing, employee_id, leave_type_id, year
	`, employeeID, leaveTypeID, year, defaultEntitlement).
		Scan(&balance.ID, &balance.Opening, &balance.Accrued, &balance.Used, &balance.Adjusted, &balance.Closing, &balance.EmployeeID, &balance.LeaveTypeID, &balance.Year)
	return balance, err
}

// UpdateLeaveBalanceAdjustment updates adjusted and closing values for leave balance
func (r *Repository) UpdateLeaveBalanceAdjustment(tx *sqlx.Tx, balanceID uuid.UUID, newAdjusted, newClosing float64) error {
	query := `
		UPDATE Tbl_Leave_balance
		SET adjusted=$1, closing=$2, updated_at=NOW()
		WHERE id=$3
	`
	_, err := tx.Exec(query, newAdjusted, newClosing, balanceID)
	return err
}

// InsertLeaveAdjustment inserts a record into leave adjustment log
func (r *Repository) InsertLeaveAdjustment(tx *sqlx.Tx, employeeID uuid.UUID, leaveTypeID int, quantity float64, reason string, createdBy string, year int) error {
	query := `
		INSERT INTO Tbl_Leave_adjustment
		(employee_id, leave_type_id, quantity, reason, created_by, created_at, year)
		VALUES ($1,$2,$3,$4,$5,NOW(),$6)
	`
	_, err := tx.Exec(query, employeeID, leaveTypeID, quantity, reason, createdBy, year)
	return err
}
