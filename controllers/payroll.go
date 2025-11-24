package controllers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/utils"
)

// PayrollPreview represents preview data for a payroll run
type PayrollPreview struct {
	EmployeeID  uuid.UUID `json:"employee_id"`
	Employee    string    `json:"employee"`
	BasicSalary float64   `json:"basic_salary"`
	WorkingDays int       `json:"working_days"`
	AbsentDays  int       `json:"absent_days"`
	Deductions  float64   `json:"deductions"`
	NetSalary   float64   `json:"net_salary"`
}

// RunPayroll handles payroll preview
func (h *HandlerFunc) RunPayroll(c *gin.Context) {
	roleRaw, _ := c.Get("role")
	role := roleRaw.(string)
	if role != "SUPERADMIN" && role != "ADMIN" {
		utils.RespondWithError(c, 403, "Not authorized to run payroll")
		return
	}

	var input struct {
		Month int `json:"month" validate:"required"`
		Year  int `json:"year" validate:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, 400, "Invalid input: "+err.Error())
		return
	}

	// Fetch active employees
	var employees []struct {
		ID       uuid.UUID `db:"id"`
		FullName string    `db:"full_name"`
		Salary   float64   `db:"salary"`
		Status   string    `db:"status"`
	}
	err := h.Query.DB.Select(&employees, `SELECT id, full_name, salary, status FROM Tbl_Employee WHERE status='active'`)
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to fetch employees: "+err.Error())
		return
	}

	// Fetch working days from company settings
	var workingDays int
	err = h.Query.DB.Get(&workingDays, `
		SELECT working_days_per_month 
		FROM Tbl_Company_Settings 
		ORDER BY created_at DESC 
		LIMIT 1
	`)
	if err != nil {
		workingDays = 22 // fallback default
	}

	totalPayroll := 0.0
	totalDeductions := 0.0
	var previews []PayrollPreview

	for _, emp := range employees {
		// Calculate absent days from approved leaves
		var absentDays float64
		err := h.Query.DB.Get(&absentDays, `
			SELECT COALESCE(SUM(days),0) 
			FROM Tbl_Leave
			WHERE employee_id=$1 AND status='APPROVED'
			AND EXTRACT(MONTH FROM start_date)=$2
			AND EXTRACT(YEAR FROM start_date)=$3
		`, emp.ID, input.Month, input.Year)
		if err != nil {
			utils.RespondWithError(c, 500, "Failed to calculate absent days: "+err.Error())
			return
		}

		deduction := emp.Salary / float64(workingDays) * absentDays
		net := emp.Salary - deduction

		previews = append(previews, PayrollPreview{
			EmployeeID:  emp.ID,
			Employee:    emp.FullName,
			BasicSalary: emp.Salary,
			WorkingDays: workingDays,
			AbsentDays:  int(absentDays),
			Deductions:  deduction,
			NetSalary:   net,
		})

		totalPayroll += net
		totalDeductions += deduction
	}

	// Create a payroll run record
	runID := uuid.New()
	_, err = h.Query.DB.Exec(
		`INSERT INTO Tbl_Payroll_run (id, month, year, status) VALUES ($1,$2,$3,$4)`,
		runID, input.Month, input.Year, "PREVIEW",
	)
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to create payroll run: "+err.Error())
		return
	}

	c.JSON(200, gin.H{
		"payroll_run_id":   runID,
		"month":            input.Month,
		"year":             input.Year,
		"total_payroll":    totalPayroll,
		"total_deductions": totalDeductions,
		"employees_count":  len(employees),
		"payroll_preview":  previews,
	})
}

// FinalizePayroll - generates payslips
func (h *HandlerFunc) FinalizePayroll(c *gin.Context) {
	roleRaw, _ := c.Get("role")
	role := roleRaw.(string)
	if role != "SUPERADMIN" && role != "ADMIN" {
		utils.RespondWithError(c, 403, "Not authorized to finalize payroll")
		return
	}

	runID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondWithError(c, 400, "Invalid payroll run ID")
		return
	}

	// Fetch payroll run month/year
	var run struct {
		Month  int    `db:"month"`
		Year   int    `db:"year"`
		Status string `db:"status"`
	}
	err = h.Query.DB.Get(&run, `SELECT month, year, status FROM Tbl_Payroll_run WHERE id=$1`, runID)
	if err != nil {
		utils.RespondWithError(c, 404, "Payroll run not found")
		return
	}
	if run.Status == "FINALIZED" {
		utils.RespondWithError(c, 400, "Payroll already finalized")
		return
	}

	// Fetch working days from company settings
	var workingDaysSetting int
	err = h.Query.DB.Get(&workingDaysSetting, `SELECT working_days_per_month FROM Tbl_Company_Settings LIMIT 1`)
	if err != nil || workingDaysSetting <= 0 {
		workingDaysSetting = 22 // fallback default
	}

	tx, err := h.Query.DB.Beginx()
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to start transaction")
		return
	}
	defer tx.Rollback()

	// Fetch active employees
	var employees []struct {
		ID       uuid.UUID `db:"id"`
		FullName string    `db:"full_name"`
		Salary   float64   `db:"salary"`
	}
	err = tx.Select(&employees, `SELECT id, full_name, salary FROM Tbl_Employee WHERE status='active'`)
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to fetch employees: "+err.Error())
		return
	}

	var payslipIDs []uuid.UUID

	for _, emp := range employees {
		var absentDays float64
		err := tx.Get(&absentDays, `
			SELECT COALESCE(SUM(days),0) FROM Tbl_Leave
			WHERE employee_id=$1 AND status='APPROVED'
			AND EXTRACT(MONTH FROM start_date)=$2
			AND EXTRACT(YEAR FROM start_date)=$3
		`, emp.ID, run.Month, run.Year)
		if err != nil {
			utils.RespondWithError(c, 500, "Failed to calculate absent days: "+err.Error())
			return
		}

		deduction := emp.Salary / float64(workingDaysSetting) * absentDays
		net := emp.Salary - deduction

		payslipID := uuid.New()
		_, err = tx.Exec(`
			INSERT INTO Tbl_Payslip
			(id, payroll_run_id, employee_id, basic_salary, working_days, absent_days, deduction_amount, net_salary)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`, payslipID, runID, emp.ID, emp.Salary, workingDaysSetting, int(absentDays), deduction, net)
		if err != nil {
			utils.RespondWithError(c, 500, "Failed to insert payslip: "+err.Error())
			return
		}
		payslipIDs = append(payslipIDs, payslipID)
	}

	_, err = tx.Exec(`UPDATE Tbl_Payroll_run SET status='FINALIZED', updated_at=NOW() WHERE id=$1`, runID)
	if err != nil {
		utils.RespondWithError(c, 500, "Failed to update payroll run: "+err.Error())
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message":           "Payroll finalized",
		"payslip_ids":       payslipIDs,
		"working_days_used": workingDaysSetting,
	})
}

// GetPayslipPDF - GET /api/payroll/payslips/:id/pdf
func (h *HandlerFunc) GetPayslipPDF(c *gin.Context) {
	payslipID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid payslip ID"})
		return
	}

	var payslip struct {
		EmployeeName string  `db:"full_name"`
		Month        int     `db:"month"`
		Year         int     `db:"year"`
		BasicSalary  float64 `db:"basic_salary"`
		WorkingDays  int     `db:"working_days"`
		AbsentDays   int     `db:"absent_days"`
		Deductions   float64 `db:"deduction_amount"`
		NetSalary    float64 `db:"net_salary"`
	}

	err = h.Query.DB.Get(&payslip, `
		SELECT e.full_name, p.basic_salary, p.working_days, p.absent_days, p.deduction_amount, p.net_salary,
		       pr.month, pr.year
		FROM Tbl_Payslip p
		JOIN Tbl_Employee e ON e.id = p.employee_id
		JOIN Tbl_Payroll_run pr ON pr.id = p.payroll_run_id
		WHERE p.id = $1
	`, payslipID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Payslip not found: " + err.Error()})
		return
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// --- Header ---
	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(30, 30, 30) // Dark gray
	pdf.CellFormat(0, 12, fmt.Sprintf("Salary Payslip - %02d/%d", payslip.Month, payslip.Year),
		"", 1, "C", false, 0, "")
	pdf.Ln(6)

	// --- Employee Info ---
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(50, 10, "Employee Name:")
	pdf.SetFont("Arial", "", 16)
	pdf.Cell(0, 10, payslip.EmployeeName)
	pdf.Ln(12)

	// --- Table Header ---
	pdf.SetFont("Arial", "B", 14)
	pdf.SetFillColor(200, 200, 200) // Light gray
	pdf.CellFormat(90, 10, "Earnings / Info", "1", 0, "C", true, 0, "")
	pdf.CellFormat(90, 10, "Amount (INR)", "1", 1, "C", true, 0, "")

	// --- Table Content ---
	pdf.SetFont("Arial", "", 14)
	addRow := func(label string, value string) {
		pdf.CellFormat(90, 10, label, "1", 0, "L", false, 0, "")
		pdf.CellFormat(90, 10, value, "1", 1, "R", false, 0, "")
	}

	addRow("Basic Salary", fmt.Sprintf("%.2f", payslip.BasicSalary))
	addRow("Working Days", fmt.Sprintf("%d", payslip.WorkingDays))
	addRow("Absent Days", fmt.Sprintf("%d", payslip.AbsentDays))
	addRow("Deductions", fmt.Sprintf("-%.2f", payslip.Deductions))
	addRow("Net Salary", fmt.Sprintf("%.2f", payslip.NetSalary))

	pdf.Ln(6)

	// --- Calculation Section ---
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Calculation:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 14)
	pdf.MultiCell(0, 8, fmt.Sprintf("%.2f - %.2f = %.2f",
		payslip.BasicSalary, payslip.Deductions, payslip.NetSalary), "", "L", false)
	pdf.Ln(4)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "Deduction Formula:")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 14)
	pdf.MultiCell(0, 8, fmt.Sprintf("Basic Salary / Working Days * Absent Days = %.2f / %d * %d = %.2f",
		payslip.BasicSalary, payslip.WorkingDays, payslip.AbsentDays, payslip.Deductions), "", "L", false)

	// --- Save PDF ---
	os.MkdirAll("./tmp", os.ModePerm)
	pdfPath := fmt.Sprintf("./tmp/payslip_%s.pdf", payslipID)
	err = pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	// Update PDF path in DB
	_, _ = h.Query.DB.Exec("UPDATE Tbl_Payslip SET pdf_path=$1, updated_at=NOW() WHERE id=$2", pdfPath, payslipID)

	// Serve PDF
	c.File(pdfPath)
}
