package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
)

type EmployeeAuthData struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type Repository struct {
	DB *sqlx.DB
}

func InitializeRepo(db *sqlx.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetEmployeeByEmail(email string) (EmployeeAuthData, error) {
	var emp EmployeeAuthData
	query := `
		SELECT 
			e.id,
			e.email,
			e.password,
			r.type AS role
		FROM Tbl_Employee e
		JOIN Tbl_Role r ON e.role_id = r.id
		WHERE e.email = $1
		LIMIT 1;
	`
	err := r.DB.Get(&emp, query, email)
	return emp, err
}

func (r *Repository) GetAllEmployees() (*sql.Rows, error) {
	query := `
        SELECT 
            e.id, e.full_name, e.email, e.status,
            r.type AS role, e.password, e.manager_id,
            e.salary, e.joining_date,
            e.created_at, e.updated_at, e.deleted_at
        FROM Tbl_Employee e
        JOIN Tbl_Role r ON e.role_id = r.id
        ORDER BY e.full_name
    `
	return r.DB.Query(query)
}

// ------------------ CHECK EMAIL EXISTS ------------------
func (r *Repository) CheckEmailExists(email string) (bool, error) {
	var existing string
	err := r.DB.QueryRow(
		`SELECT email FROM Tbl_Employee WHERE email=$1`, email,
	).Scan(&existing)

	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// ------------------ GET ROLE ID ------------------
func (r *Repository) GetRoleID(role string) (string, error) {
	var id string
	err := r.DB.QueryRow(`SELECT id FROM Tbl_Role WHERE type=$1`, role).Scan(&id)
	return id, err
}

// ------------------ CREATE EMPLOYEE ------------------
func (r *Repository) InsertEmployee(fullName, email, roleID, password string, salary *float64, joining *time.Time) error {
	_, err := r.DB.Exec(`
		INSERT INTO Tbl_Employee (full_name, email, role_id, password, salary, joining_date)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, fullName, email, roleID, password, salary, joining)
	return err
}

// ------------------ GET CURRENT ROLE NAME ------------------
func (r *Repository) GetEmployeeCurrentRole(empID string) (string, error) {
	var role string
	err := r.DB.QueryRow(`
        SELECT R.TYPE
        FROM TBL_EMPLOYEE E
        JOIN TBL_ROLE R ON E.ROLE_ID = R.ID
        WHERE E.ID = $1
    `, empID).Scan(&role)
	return role, err
}

// ------------------ UPDATE ROLE ------------------
func (r *Repository) UpdateEmployeeRole(empID string, newRole string) (string, error) {
	var id string
	query := `
        UPDATE TBL_EMPLOYEE
        SET ROLE_ID = (SELECT ID FROM TBL_ROLE WHERE TYPE=$1),
            UPDATED_AT = NOW()
        WHERE ID = $2
        RETURNING ID;
    `
	err := r.DB.QueryRow(query, newRole, empID).Scan(&id)
	return id, err
}

// ------------------ CHECK MANAGER EXISTS ------------------
func (r *Repository) ManagerExists(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM TBL_EMPLOYEE WHERE ID=$1)`,
		id,
	).Scan(&exists)
	return exists, err
}

// ------------------ UPDATE MANAGER ------------------
func (r *Repository) UpdateManager(empID, managerID uuid.UUID) error {
	_, err := r.DB.Exec(`
        UPDATE TBL_EMPLOYEE
        SET MANAGER_ID=$1, UPDATED_AT=NOW()
        WHERE ID=$2
    `, managerID, empID)
	return err
}

// AddHoliday inserts a holiday into the database
func (r *Repository) AddHoliday(name string, date time.Time, typ string) (string, error) {
	if typ == "" {
		typ = "HOLIDAY"
	}
	day := date.Weekday().String()
	var id string
	err := r.DB.QueryRow(`
		INSERT INTO Tbl_Holiday (name, date, day, type, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`, name, date, day, typ).Scan(&id)
	return id, err
}

// GetAllHolidays fetches all holidays
func (r *Repository) GetAllHolidays() ([]models.Holiday, error) {
	rows, err := r.DB.Queryx(`SELECT id, name, date, day, type, created_at, updated_at FROM Tbl_Holiday ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holidays []models.Holiday
	for rows.Next() {
		var h models.Holiday
		if err := rows.StructScan(&h); err != nil {
			return nil, err
		}
		holidays = append(holidays, h)
	}
	return holidays, nil
}

// DeleteHoliday deletes a holiday by ID
func (r *Repository) DeleteHoliday(id string) error {
	_, err := r.DB.Exec(`DELETE FROM Tbl_Holiday WHERE id=$1`, id)
	return err
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
