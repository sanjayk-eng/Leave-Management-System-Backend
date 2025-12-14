package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
)

// 1. Get leave type entitlement
func (r *Repository) GetLeaveTypeByIdTx(tx *sqlx.Tx, leaveTypeID int) (models.LeaveType, error) {
	var leaves models.LeaveType
	query := `SELECT id, name, is_paid, default_entitlement,  created_at, updated_at FROM Tbl_Leave_type WHERE id=$1`
	err := tx.Get(&leaves,
		query,
		leaveTypeID,
	)
	return leaves, err
}

// 1. Get leave type entitlement
func (r *Repository) GetLeaveTypeById(leaveTypeID int) (models.LeaveType, error) {
	var leaves models.LeaveType
	query := `SELECT id, name, is_paid, default_entitlement,  created_at, updated_at FROM Tbl_Leave_type WHERE id=$1`
	err := r.DB.Get(&leaves,
		query,
		leaveTypeID,
	)
	return leaves, err
}

func (q *Repository) GetLeaveTypeByLeaveID(leaveID uuid.UUID) (int, error) {
	var leaveTypeID int
	err := q.DB.Get(&leaveTypeID, `
        SELECT leave_type_id 
        FROM Tbl_Leave 
        WHERE id = $1
    `, leaveID)

	if err != nil {
		return 0, err
	}

	return leaveTypeID, nil
}

func (r *Repository) GetAllLeaveType() ([]models.LeaveType, error) {
	var leaveType []models.LeaveType
	query := `SELECT id, name, is_paid, default_entitlement,  created_at, updated_at FROM Tbl_Leave_type ORDER BY id`
	err := r.DB.Select(&leaveType, query)
	return leaveType, err
}

// Admin add leave type

func (r *Repository) AddLeaveType(tx *sqlx.Tx, input models.LeaveTypeInput) (models.LeaveType, error) {
	var leave models.LeaveType
	query := `
		INSERT INTO Tbl_Leave_type (name, is_paid, default_entitlement)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := tx.QueryRow(query, input.Name, *input.IsPaid, *input.DefaultEntitlement).
		Scan(&leave.ID, &leave.CreatedAt, &leave.UpdatedAt)
	return leave, err
}

// 3. Get leave balance (inside TX)
func (r *Repository) GetLeaveBalance(tx *sqlx.Tx, employeeID uuid.UUID, leaveTypeID int) (float64, error) {
	var balance float64
	err := tx.Get(&balance, `
		SELECT closing 
		FROM Tbl_Leave_balance 
		WHERE employee_id=$1 AND leave_type_id=$2 
		AND year = EXTRACT(YEAR FROM CURRENT_DATE)
	`, employeeID, leaveTypeID)
	return balance, err
}

// create leave balance
func (r *Repository) CreateLeaveBalance(tx *sqlx.Tx, employeeID uuid.UUID, leaveTypeID int, entitlement int) error {
	_, err := tx.Exec(`
		INSERT INTO Tbl_Leave_balance 
			(employee_id, leave_type_id, year, opening, accrued, used, adjusted, closing)
		VALUES ($1, $2, EXTRACT(YEAR FROM CURRENT_DATE), $3, 0, 0, 0, $3)
	`, employeeID, leaveTypeID, entitlement)
	return err
}

// 5. Check overlapping leaves
func (r *Repository) GetOverlappingLeaves(
	tx *sqlx.Tx,
	employeeID uuid.UUID,
	startDate, endDate time.Time,
) ([]struct {
	ID        uuid.UUID `db:"id"`
	LeaveType string    `db:"leave_type"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	Status    string    `db:"status"`
}, error) {

	var result []struct {
		ID        uuid.UUID `db:"id"`
		LeaveType string    `db:"leave_type"`
		StartDate time.Time `db:"start_date"`
		EndDate   time.Time `db:"end_date"`
		Status    string    `db:"status"`
	}

	err := tx.Select(&result, `
		SELECT l.id, lt.name as leave_type, l.start_date, l.end_date, l.status
		FROM Tbl_Leave l
		JOIN Tbl_Leave_type lt ON l.leave_type_id = lt.id
		WHERE l.employee_id=$1 
		AND l.status IN ('Pending','APPROVED')
		AND l.start_date <= $2 
		AND l.end_date >= $3
	`, employeeID, endDate, startDate)

	return result, err
}

func (r *Repository) InsertLeave(
	tx *sqlx.Tx,
	employeeID uuid.UUID,
	leaveTypeID int,
	leaveTimingID int,
	startDate, endDate time.Time,
	days float64,
	reason string,
) (uuid.UUID, error) {

	var leaveID uuid.UUID

	err := tx.QueryRow(`
		INSERT INTO Tbl_Leave 
		(employee_id, leave_type_id, half_id, start_date, end_date, days, status, reason)
		VALUES ($1,$2,$3,$4,$5,$6,'Pending',$7)
		RETURNING id
	`,
		employeeID,
		leaveTypeID,
		leaveTimingID,
		startDate,
		endDate,
		days,
		reason,
	).Scan(&leaveID)

	return leaveID, err
}

func (r *Repository) GetLeaveById(tx *sqlx.Tx, leaveID uuid.UUID) (models.Leave, error) {
	var leave models.Leave
	query := `SELECT * FROM Tbl_Leave WHERE id=$1 FOR UPDATE`
	err := tx.Get(&leave, query, leaveID)
	return leave, err
}

// Get All Leave Timming
// Get All Leave Timing
func (r *Repository) GetLeaveTiming() ([]models.LeaveTimingResponse, error) {
	var data []models.LeaveTimingResponse
	query := `
		SELECT id, type, timing, created_at, updated_at
		FROM Tbl_Half
		ORDER BY id
	`
	err := r.DB.Select(&data, query)

	return data, err
}

// Get Leave Timing By ID
func (r *Repository) GetLeaveTimingByID(id int) (*models.LeaveTimingResponse, error) {
	var data models.LeaveTimingResponse

	query := `
		SELECT *
		FROM Tbl_Half
		WHERE id = $1
	`

	err := r.DB.Get(&data, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &data, nil
}

func (r *Repository) UpdateLeaveTiming(tx *sqlx.Tx, id int, timing string) error {
	query := `
		UPDATE Tbl_Half
		SET timing = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	res, err := tx.Exec(query, timing, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
