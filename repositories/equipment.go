package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
)

//
// ======================
// CATEGORY REPOSITORIES
// ======================
//

// CreateCategory
func (r *Repository) CreateCategory(tx *sqlx.Tx, data models.EquipmentCategoryRequest) error {
	_, err := tx.Exec(`
		INSERT INTO tbl_equipment_category (name, description)
		VALUES ($1,$2)
	`, data.Name, data.Description)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// GetAllCategory
func (r *Repository) GetAllCategory() ([]models.EquipmentCategoryRes, error) {
	var res []models.EquipmentCategoryRes

	err := r.DB.Select(&res, `
		SELECT id, name, description, created_at, updated_at
		FROM tbl_equipment_category
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateCategory
func (r *Repository) UpdateCategory(tx *sqlx.Tx, id uuid.UUID, data models.EquipmentCategoryRequest) error {
	result, err := tx.Exec(`
		UPDATE tbl_equipment_category
		SET name=$1, description=$2, updated_at=now()
		WHERE id=$3
	`, data.Name, data.Description, id)

	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// DeleteCategory
func (r *Repository) DeleteCategory(tx *sqlx.Tx, id uuid.UUID) error {
	result, err := tx.Exec(`DELETE FROM tbl_equipment_category WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

//
// ======================
// EQUIPMENT REPOSITORIES
// ======================
//

func (r *Repository) CreateEquipment(tx *sqlx.Tx, data models.EquipmentRequest) error {
	_, err := tx.Exec(`
		INSERT INTO tbl_equipment
		(name, category_id, is_shared, price, total_quantity, remaining_quantity, purchase_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		data.Name,
		data.CategoryID,
		data.IsShared,
		data.Price,
		data.TotalQuantity,
		data.TotalQuantity, // RemainingQuantity initially equals TotalQuantity
		data.PurchaseDate,  // Optional purchase date
	)
	if err != nil {
		return fmt.Errorf("failed to create equipment: %w", err)
	}
	return nil
}

// GetAllEquipment
func (r *Repository) GetAllEquipment() ([]models.EquipmentRes, error) {
	res := []models.EquipmentRes{}

	err := r.DB.Select(&res, `
		SELECT id, name, category_id, is_shared, price,
		       total_quantity, remaining_quantity,
		       purchase_date,
		       created_at, updated_at
		FROM tbl_equipment
		ORDER BY name
	`)
	return res, err
}

// GetEquipmentByCategory
func (r *Repository) GetEquipmentByCategory(categoryID uuid.UUID) ([]models.EquipmentRes, error) {
	res := []models.EquipmentRes{}

	err := r.DB.Select(&res, `
		SELECT id, name, category_id, is_shared, price,
		       total_quantity, remaining_quantity,
		       purchase_date,
		       created_at, updated_at
		FROM tbl_equipment
		WHERE category_id = $1
		ORDER BY name
	`, categoryID)

	return res, err
}

// UpdateEquipment
func (r *Repository) UpdateEquipment(tx *sqlx.Tx, id uuid.UUID, data models.EquipmentRequest) error {
	// Get current remaining quantity and total quantity
	var currentRemaining, currentTotal int
	err := tx.QueryRow(`
		SELECT remaining_quantity, total_quantity 
		FROM tbl_equipment 
		WHERE id = $1
	`, id).Scan(&currentRemaining, &currentTotal)
	if err != nil {
		return fmt.Errorf("equipment not found")
	}

	// Calculate new remaining quantity based on the difference in total quantity
	// If total quantity increases, add the difference to remaining quantity
	// If total quantity decreases, subtract from remaining (but not below 0)
	newRemaining := currentRemaining + (data.TotalQuantity - currentTotal)
	if newRemaining < 0 {
		newRemaining = 0
	}

	result, err := tx.Exec(`
		UPDATE tbl_equipment
		SET name = $1,
		    category_id = $2,
		    is_shared = $3,
		    price = $4,
		    total_quantity = $5,
		    remaining_quantity = $6,
		    purchase_date = $7,
		    updated_at = now()
		WHERE id = $8
	`,
		data.Name,
		data.CategoryID,
		data.IsShared,
		data.Price,
		data.TotalQuantity,
		newRemaining,
		data.PurchaseDate,
		id,
	)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("equipment not found")
	}

	return nil
}

// DeleteEquipment
func (r *Repository) DeleteEquipment(tx *sqlx.Tx, id uuid.UUID) error {
	result, err := tx.Exec(`DELETE FROM tbl_equipment WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("equipment not found")
	}
	return nil
}

//
// ======================
// ASSIGNMENT REPOSITORIES
// ======================
//

// AssignEquipment
func (r *Repository) AssignEquipment(tx *sqlx.Tx, req models.AssignEquipmentRequest) error {
	var remaining int
	var isShared bool

	err := tx.QueryRow(`
		SELECT remaining_quantity, is_shared
		FROM tbl_equipment
		WHERE id=$1
		FOR UPDATE
	`, req.EquipmentID).Scan(&remaining, &isShared)
	if err != nil {
		return fmt.Errorf("equipment not found")
	}

	if (!isShared && remaining < 1) || (isShared && remaining < req.Quantity) {
		return fmt.Errorf("not enough equipment available")
	}

	_, err = tx.Exec(`
		INSERT INTO tbl_equipment_assignment
		(equipment_id, employee_id, assigned_by, quantity)
		VALUES ($1,$2,$3,$4)
	`, req.EquipmentID, req.EmployeeID, req.AssignedBy, req.Quantity)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE tbl_equipment
		SET remaining_quantity = remaining_quantity - $1
		WHERE id=$2
	`, req.Quantity, req.EquipmentID)

	return err
}

// GetAllAssignedEquipment
func (r *Repository) GetAllAssignedEquipment() ([]models.AssignEquipmentResponse, error) {
	res := []models.AssignEquipmentResponse{}

	err := r.DB.Select(&res, `
		SELECT 
		       e.full_name AS employee_name,
		       e.email AS employee_email,
		       eq.name AS equipment_name,
		       eq.purchase_date,
		       ea.quantity,
		       ab.full_name AS approved_by_name
		FROM tbl_equipment_assignment ea
		JOIN tbl_employee e  ON e.id = ea.employee_id
		JOIN tbl_equipment eq ON eq.id = ea.equipment_id
		JOIN tbl_employee ab ON ab.id = ea.assigned_by
		WHERE ea.returned_at IS NULL
		ORDER BY ea.assigned_at DESC
	`)
	return res, err
}

// GetAssignedEquipmentByEmployee
func (r *Repository) GetAssignedEquipmentByEmployee(employeeID uuid.UUID) ([]models.AssignEquipmentResponse, error) {
	res := []models.AssignEquipmentResponse{}

	err := r.DB.Select(&res, `
		SELECT 
		    e.full_name AS employee_name,
		    e.email AS employee_email,
		    eq.name AS equipment_name,
		    eq.purchase_date,
		    ea.quantity,
		    ab.full_name AS approved_by_name
		FROM tbl_equipment_assignment ea
		JOIN tbl_employee e ON e.id = ea.employee_id
		JOIN tbl_equipment eq ON eq.id = ea.equipment_id
		JOIN tbl_employee ab ON ab.id = ea.assigned_by
		WHERE ea.returned_at IS NULL AND e.id=$1
		ORDER BY ea.assigned_at DESC
	`, employeeID)

	return res, err
}

// RemoveEquipment (Return)
func (r *Repository) RemoveEquipment(tx *sqlx.Tx, req models.RemoveEquipmentRequest) error {
	var qty int

	err := tx.Get(&qty, `
		SELECT quantity
		FROM tbl_equipment_assignment
		WHERE equipment_id=$1 AND employee_id=$2 AND returned_at IS NULL
	`, req.EquipmentID, req.EmployeeID)
	if err != nil {
		return fmt.Errorf("active assignment not found")
	}

	_, err = tx.Exec(`
		UPDATE tbl_equipment_assignment
		SET returned_at=now()
		WHERE equipment_id=$1 AND employee_id=$2 AND returned_at IS NULL
	`, req.EquipmentID, req.EmployeeID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE tbl_equipment
		SET remaining_quantity = remaining_quantity + $1
		WHERE id=$2
	`, qty, req.EquipmentID)

	return err
}

// UpdateAssignment (Quantity change OR reassignment)
func (r *Repository) UpdateAssignment(tx *sqlx.Tx, req models.UpdateAssignmentRequest) error {
	var currentQty int

	// 1️ Get current assignment quantity
	err := tx.Get(&currentQty, `
		SELECT quantity
		FROM tbl_equipment_assignment
		WHERE equipment_id=$1 AND employee_id=$2 AND returned_at IS NULL
	`, req.EquipmentID, req.FromEmployeeID)
	if err != nil {
		return fmt.Errorf("no active assignment found")
	}

	// 2️ Reassignment to another employee
	if req.ToEmployeeID != nil {
		if req.Quantity > currentQty {
			return fmt.Errorf("quantity exceeds assigned amount")
		}

		// Reduce or remove from current employee
		if req.Quantity == currentQty {
			_, err = tx.Exec(`
				UPDATE tbl_equipment_assignment
				SET returned_at=now()
				WHERE equipment_id=$1 AND employee_id=$2 AND returned_at IS NULL
			`, req.EquipmentID, req.FromEmployeeID)
		} else {
			_, err = tx.Exec(`
				UPDATE tbl_equipment_assignment
				SET quantity = quantity - $1
				WHERE equipment_id=$2 AND employee_id=$3 AND returned_at IS NULL
			`, req.Quantity, req.EquipmentID, req.FromEmployeeID)
		}
		if err != nil {
			return err
		}

		// Assign to new employee
		_, err = tx.Exec(`
			INSERT INTO tbl_equipment_assignment
			(equipment_id, employee_id, assigned_by, quantity)
			VALUES ($1,$2,$3,$4)
		`, req.EquipmentID, *req.ToEmployeeID, req.AssignedBy, req.Quantity)
		return err
	}

	// 3️ Quantity update for same employee
	diff := req.Quantity - currentQty
	if diff > 0 {
		var remaining int
		err := tx.Get(&remaining, `
			SELECT remaining_quantity
			FROM tbl_equipment
			WHERE id=$1
			FOR UPDATE
		`, req.EquipmentID)
		if err != nil || remaining < diff {
			return fmt.Errorf("not enough quantity available")
		}

		_, _ = tx.Exec(`
			UPDATE tbl_equipment
			SET remaining_quantity = remaining_quantity - $1
			WHERE id=$2
		`, diff, req.EquipmentID)
	} else if diff < 0 {
		_, _ = tx.Exec(`
			UPDATE tbl_equipment
			SET remaining_quantity = remaining_quantity + $1
			WHERE id=$2
		`, -diff, req.EquipmentID)
	}

	// Update assignment quantity
	_, err = tx.Exec(`
		UPDATE tbl_equipment_assignment
		SET quantity=$1
		WHERE equipment_id=$2 AND employee_id=$3 AND returned_at IS NULL
	`, req.Quantity, req.EquipmentID, req.FromEmployeeID)

	return err
}
