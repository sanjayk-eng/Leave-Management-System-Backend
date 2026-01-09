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
		WHERE e.id = $1
		ORDER BY ea.assigned_at DESC
	`, employeeID)

	return res, err
}

// RemoveEquipment (Return)
// HARD DELETE: Permanently removes the most recent assignment from the database
func (r *Repository) RemoveEquipment(tx *sqlx.Tx, req models.RemoveEquipmentRequest) error {
	var assignmentID uuid.UUID
	var qty int

	// Get the most recent assignment (by assigned_at DESC)
	// This handles cases where the same equipment was assigned multiple times to the same employee
	err := tx.QueryRow(`
		SELECT id, quantity
		FROM tbl_equipment_assignment
		WHERE equipment_id = $1 AND employee_id = $2
		ORDER BY assigned_at DESC
		LIMIT 1
	`, req.EquipmentID, req.EmployeeID).Scan(&assignmentID, &qty)
	
	if err != nil {
		return fmt.Errorf("assignment not found for equipment_id=%s, employee_id=%s: %w", 
			req.EquipmentID, req.EmployeeID, err)
	}

	// HARD DELETE: Permanently delete the assignment row from database
	result, err := tx.Exec(`
		DELETE FROM tbl_equipment_assignment
		WHERE id = $1
	`, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to delete assignment id=%s: %w", assignmentID, err)
	}

	// Check if assignment was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("assignment id=%s not found", assignmentID)
	}

	// Update equipment remaining quantity (add back the quantity that was assigned)
	result2, err := tx.Exec(`
		UPDATE tbl_equipment
		SET remaining_quantity = remaining_quantity + $1
		WHERE id = $2
	`, qty, req.EquipmentID)
	if err != nil {
		return fmt.Errorf("failed to update equipment quantity: %w", err)
	}

	// Verify equipment was updated
	rowsAffected2, err := result2.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check equipment update: %w", err)
	}
	if rowsAffected2 == 0 {
		return fmt.Errorf("equipment id=%s not found", req.EquipmentID)
	}

	return nil
}

// UpdateAssignment (Quantity change OR reassignment)
func (r *Repository) UpdateAssignment(tx *sqlx.Tx, req models.UpdateAssignmentRequest) error {
	var assignmentID uuid.UUID
	var currentQty int

	// 1️ Get the most recent assignment ID and quantity
	// CRITICAL FIX: Get assignment ID first to ensure we update only ONE assignment
	// This prevents updating all assignments when multiple exist
	err := tx.QueryRow(`
		SELECT id, quantity
		FROM tbl_equipment_assignment
		WHERE equipment_id = $1 AND employee_id = $2
		ORDER BY assigned_at DESC
		LIMIT 1
	`, req.EquipmentID, req.FromEmployeeID).Scan(&assignmentID, &currentQty)
	if err != nil {
		return fmt.Errorf("assignment not found for equipment_id=%s, employee_id=%s: %w", 
			req.EquipmentID, req.FromEmployeeID, err)
	}

	// 2️ Reassignment to another employee
	if req.ToEmployeeID != nil {
		if req.Quantity > currentQty {
			return fmt.Errorf("quantity %d exceeds assigned amount %d", req.Quantity, currentQty)
		}

		// ========================================
		// STEP 1: REMOVE FROM CURRENT EMPLOYEE FIRST
		// ========================================
		// This ensures current employee's assignment is removed before assigning to new employee
		if req.Quantity == currentQty {
			// Full reassignment - HARD DELETE the current employee's assignment
			result, err := tx.Exec(`
				DELETE FROM tbl_equipment_assignment
				WHERE id = $1
			`, assignmentID)
			if err != nil {
				return fmt.Errorf("failed to delete current employee assignment: %w", err)
			}
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return fmt.Errorf("current employee assignment id=%s not found", assignmentID)
			}
		} else {
			// Partial reassignment - reduce quantity from current employee
			result, err := tx.Exec(`
				UPDATE tbl_equipment_assignment
				SET quantity = quantity - $1
				WHERE id = $2
			`, req.Quantity, assignmentID)
			if err != nil {
				return fmt.Errorf("failed to reduce current employee assignment quantity: %w", err)
			}
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return fmt.Errorf("current employee assignment id=%s not found", assignmentID)
			}
		}

		// ========================================
		// STEP 2: Add back quantity to equipment remaining_quantity
		// ========================================
		// This frees up the quantity that was assigned to current employee
		_, err = tx.Exec(`
			UPDATE tbl_equipment
			SET remaining_quantity = remaining_quantity + $1
			WHERE id = $2
		`, req.Quantity, req.EquipmentID)
		if err != nil {
			return fmt.Errorf("failed to update equipment quantity: %w", err)
		}

		// ========================================
		// STEP 3: Assign to new employee
		// ========================================
		// Check if new employee already has an assignment for this equipment
		// Get the most recent assignment (if exists) to update only that one
		var newEmployeeAssignmentID uuid.UUID
		var existingQty int
		err = tx.QueryRow(`
			SELECT id, quantity
			FROM tbl_equipment_assignment
			WHERE equipment_id = $1 AND employee_id = $2
			ORDER BY assigned_at DESC
			LIMIT 1
		`, req.EquipmentID, *req.ToEmployeeID).Scan(&newEmployeeAssignmentID, &existingQty)

		if err == nil {
			// New employee already has assignment - update the most recent one
			// No need to reduce remaining_quantity because we already added it back in Step 2
			_, err = tx.Exec(`
				UPDATE tbl_equipment_assignment
				SET quantity = quantity + $1, assigned_by = $2
				WHERE id = $3
			`, req.Quantity, req.AssignedBy, newEmployeeAssignmentID)
			if err != nil {
				return fmt.Errorf("failed to update new employee assignment: %w", err)
			}
			// No need to reduce remaining_quantity - it's a transfer, not new assignment
		} else {
			// New employee doesn't have assignment - create new one
			// Reduce remaining_quantity because we're assigning to new employee
			// (We added it back in Step 2, now we reduce it again for the new assignment)
			_, err = tx.Exec(`
				INSERT INTO tbl_equipment_assignment
				(equipment_id, employee_id, assigned_by, quantity)
				VALUES ($1, $2, $3, $4)
			`, req.EquipmentID, *req.ToEmployeeID, req.AssignedBy, req.Quantity)
			if err != nil {
				return fmt.Errorf("failed to create new employee assignment: %w", err)
			}
			// Reduce remaining_quantity because we're assigning to new employee
			_, err = tx.Exec(`
				UPDATE tbl_equipment
				SET remaining_quantity = remaining_quantity - $1
				WHERE id = $2
			`, req.Quantity, req.EquipmentID)
			if err != nil {
				return fmt.Errorf("failed to update equipment quantity: %w", err)
			}
		}

		return nil
	}

	// 3️ Quantity update for same employee
	diff := req.Quantity - currentQty
	if diff > 0 {
		var remaining int
		err := tx.Get(&remaining, `
			SELECT remaining_quantity
			FROM tbl_equipment
			WHERE id = $1
			FOR UPDATE
		`, req.EquipmentID)
		if err != nil {
			return fmt.Errorf("equipment not found: %w", err)
		}
		if remaining < diff {
			return fmt.Errorf("not enough quantity available: need %d, have %d", diff, remaining)
		}

		_, err = tx.Exec(`
			UPDATE tbl_equipment
			SET remaining_quantity = remaining_quantity - $1
			WHERE id = $2
		`, diff, req.EquipmentID)
		if err != nil {
			return fmt.Errorf("failed to reduce equipment quantity: %w", err)
		}
	} else if diff < 0 {
		_, err = tx.Exec(`
			UPDATE tbl_equipment
			SET remaining_quantity = remaining_quantity + $1
			WHERE id = $2
		`, -diff, req.EquipmentID)
		if err != nil {
			return fmt.Errorf("failed to increase equipment quantity: %w", err)
		}
	}

	// CRITICAL FIX: Update assignment quantity by ID, not by equipment_id + employee_id
	// This ensures we only update the specific assignment we found
	result, err := tx.Exec(`
		UPDATE tbl_equipment_assignment
		SET quantity = $1
		WHERE id = $2
	`, req.Quantity, assignmentID)
	if err != nil {
		return fmt.Errorf("failed to update assignment quantity: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("assignment id=%s not found", assignmentID)
	}

	return nil
}
