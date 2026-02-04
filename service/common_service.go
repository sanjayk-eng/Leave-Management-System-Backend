package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/repositories"
)

func CalculateWorkingDays(Query *repositories.Repository, tx *sqlx.Tx, start, end time.Time) (float64, error) {
	// 1️ Validate date range
	if end.Before(start) {
		return 0, fmt.Errorf("end date cannot be before start date")
	}

	// Normalize dates to midnight UTC to avoid timezone issues
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	// 2️ Fetch holidays within range
	holidays, err := Query.GetByFilterHolidayBetwweenTwoDates(tx, start, end)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch holidays: %v", err)
	}

	// Convert slice to a map for O(1) lookup
	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h.Format("2006-01-02")] = true
	}

	// 3️ Count working days
	workingDays := 0
	var workingDaysList []string // For debugging

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dayStr := d.Format("2006-01-02")
		weekday := d.Weekday()

		// Skip Saturday and Sunday
		if weekday == time.Saturday || weekday == time.Sunday {
			fmt.Printf("DEBUG: Skipping weekend: %s (%s)\n", dayStr, weekday)
			continue
		}

		// Skip holidays
		if holidayMap[dayStr] {
			fmt.Printf("DEBUG: Skipping holiday: %s\n", dayStr)
			continue
		}

		workingDays++
		workingDaysList = append(workingDaysList, fmt.Sprintf("%s (%s)", dayStr, weekday))
	}

	fmt.Printf("DEBUG: Working days calculated: %d - Days: %v\n", workingDays, workingDaysList)
	return float64(workingDays), nil
}

// CalculateWorkingDaysWithTiming calculates working days based on timing type
// timingID: 1 = First Half (0.5 days), 2 = Second Half (0.5 days), 3 = Full Day (1.0 days)
func CalculateWorkingDaysWithTiming(Query *repositories.Repository, tx *sqlx.Tx, start, end time.Time, timingID int) (float64, error) {
	// First calculate the base working days
	baseDays, err := CalculateWorkingDays(Query, tx, start, end)
	if err != nil {
		return 0, err
	}

	// Apply timing multiplier
	switch timingID {
	case 1, 2: // First Half or Second Half
		return baseDays * 0.5, nil
	case 3: // Full Day
		return baseDays, nil
	default:
		return 0, fmt.Errorf("invalid timing ID: %d. Must be 1 (First Half), 2 (Second Half), or 3 (Full Day)", timingID)
	}
}

type LeaveSummary struct {
	PaidDays   float64
	UnpaidDays float64
}

// CalculateAbsentDaysForMonth calculates the number of absent days for a specific month
// Now handles timing-based leaves (half days vs full days) correctly
// Uses the pre-calculated days from leave records which include timing considerations

/*
func CalculateAbsentDaysForMonth(db *sqlx.DB, employeeID uuid.UUID, month, year int) float64 {
	// Get first and last day of the payroll month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1) // Last day of the month

	// Fetch all approved UNPAID leaves that overlap with this month
	// Include timing information to understand the leave type
	type LeaveRecord struct {
		StartDate  time.Time `db:"start_date"`
		EndDate    time.Time `db:"end_date"`
		Days       float64   `db:"days"`        // Pre-calculated days (includes timing: 0.5 for half days, 1.0+ for full days)
		TimingID   *int      `db:"half_id"`     // Timing ID (1=First Half, 2=Second Half, 3=Full Day)
		TimingType *string   `db:"timing_type"` // Timing type for debugging
	}

	var leaves []LeaveRecord
	err := db.Select(&leaves, `
		SELECT l.start_date, l.end_date, l.days, l.half_id, h.type as timing_type
		FROM Tbl_Leave l
		JOIN Tbl_Leave_type lt ON l.leave_type_id = lt.id
		LEFT JOIN Tbl_Half h ON l.half_id = h.id
		WHERE l.employee_id=$1
		AND l.status='APPROVED'
		AND lt.is_paid = false
		AND l.start_date <= $2
		AND l.end_date >= $3
	`, employeeID, lastDay, firstDay)

	if err != nil {
		fmt.Printf("Error fetching leaves for payroll: %v\n", err)
		return -1
	}

	// Calculate total absent days within this month
	totalAbsentDays := 0.0

	for _, leave := range leaves {
		// For leaves that span across months, we need to calculate the proportion
		// that falls within the payroll month

		// Determine the overlap period
		overlapStart := leave.StartDate
		if overlapStart.Before(firstDay) {
			overlapStart = firstDay
		}

		overlapEnd := leave.EndDate
		if overlapEnd.After(lastDay) {
			overlapEnd = lastDay
		}

		// Calculate working days in the overlap period (excluding weekends and holidays)
		workingDaysInOverlap := 0
		for d := overlapStart; !d.After(overlapEnd); d = d.AddDate(0, 0, 1) {
			// Skip weekends
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				continue
			}

			// Check if it's a holiday
			var isHoliday bool
			err := db.Get(&isHoliday, `
				SELECT EXISTS(SELECT 1 FROM Tbl_Holiday WHERE date=$1)
			`, d)
			if err == nil && !isHoliday {
				workingDaysInOverlap++
			}
		}

		// Calculate total working days in the entire leave period
		totalWorkingDaysInLeave := 0
		for d := leave.StartDate; !d.After(leave.EndDate); d = d.AddDate(0, 0, 1) {
			// Skip weekends
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				continue
			}

			// Check if it's a holiday
			var isHoliday bool
			err := db.Get(&isHoliday, `
				SELECT EXISTS(SELECT 1 FROM Tbl_Holiday WHERE date=$1)
			`, d)
			if err == nil && !isHoliday {
				totalWorkingDaysInLeave++
			}
		}

		// Calculate the proportional absent days for this month
		// This preserves the timing-based calculation (half days vs full days)
		if totalWorkingDaysInLeave > 0 {
			proportionalDays := leave.Days * (float64(workingDaysInOverlap) / float64(totalWorkingDaysInLeave))
			totalAbsentDays += proportionalDays

			// Debug logging
			timingType := "FULL"
			if leave.TimingType != nil {
				timingType = *leave.TimingType
			}
			fmt.Printf("DEBUG Payroll: Leave %s to %s, Type: %s, Total Days: %.1f, Overlap Days: %d/%d, Proportional: %.2f\n",
				leave.StartDate.Format("2006-01-02"), leave.EndDate.Format("2006-01-02"),
				timingType, leave.Days, workingDaysInOverlap, totalWorkingDaysInLeave, proportionalDays)
		}
	}

	fmt.Printf("DEBUG Payroll: Employee %s, Month %d/%d, Total Absent Days: %.2f\n",
		employeeID, month, year, totalAbsentDays)

	return totalAbsentDays
}*/

func CalculateAbsentDaysForMonth(db *sqlx.DB, employeeID uuid.UUID, month, year int) LeaveSummary {
	// 1. Define time boundaries
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)

	// 2. Optimization: Fetch holidays once for the period
	var holidays []time.Time
	_ = db.Select(&holidays, `SELECT date FROM Tbl_Holiday WHERE date >= $1 AND date <= $2`, firstDay, lastDay)

	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		holidayMap[h.Format("2006-01-02")] = true
	}

	// 3. Updated SQL: Removed "is_paid = false" to get ALL approved leaves
	type LeaveRecord struct {
		StartDate  time.Time `db:"start_date"`
		EndDate    time.Time `db:"end_date"`
		Days       float64   `db:"days"`
		IsPaid     bool      `db:"is_paid"` // Now fetching this field
		TimingType *string   `db:"timing_type"`
	}

	var leaves []LeaveRecord
	err := db.Select(&leaves, `
        SELECT l.start_date, l.end_date, l.days, lt.is_paid, h.type as timing_type
        FROM Tbl_Leave l
        JOIN Tbl_Leave_type lt ON l.leave_type_id = lt.id
        LEFT JOIN Tbl_Half h ON l.half_id = h.id
        WHERE l.employee_id=$1 
        AND l.status='APPROVED'
        AND l.start_date <= $2
        AND l.end_date >= $3
    `, employeeID, lastDay, firstDay)

	if err != nil {
		fmt.Printf("Error fetching leaves: %v\n", err)
		return LeaveSummary{}
	}

	summary := LeaveSummary{}

	// 4. Calculate days
	for _, leave := range leaves {
		overlapStart := leave.StartDate
		if overlapStart.Before(firstDay) {
			overlapStart = firstDay
		}

		overlapEnd := leave.EndDate
		if overlapEnd.After(lastDay) {
			overlapEnd = lastDay
		}

		actualDaysInMonth := 0.0
		for d := overlapStart; !d.After(overlapEnd); d = d.AddDate(0, 0, 1) {
			// Skip weekends and holidays
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday || holidayMap[d.Format("2006-01-02")] {
				continue
			}

			if leave.StartDate.Equal(leave.EndDate) && leave.Days < 1.0 {
				actualDaysInMonth += leave.Days
			} else {
				actualDaysInMonth += 1.0
			}
		}

		// 5. Categorize based on IsPaid flag
		if leave.IsPaid {
			summary.PaidDays += actualDaysInMonth
		} else {
			summary.UnpaidDays += actualDaysInMonth
		}
	}

	return summary
}

// LeaveBalanceData represents raw balance data from database
type LeaveBalanceData struct {
	LeaveTypeID int
	Opening     float64
	Accrued     float64
	Used        float64
	Adjusted    float64
	Closing     float64
}

// LeaveTypeData represents leave type information
type LeaveTypeData struct {
	LeaveTypeID        int
	LeaveTypeName      string
	DefaultEntitlement float64
}

// CalculatedBalance represents the calculated leave balance result
type CalculatedBalance struct {
	LeaveTypeID int     `json:"leave_type_id"`
	LeaveType   string  `json:"leave_type"`
	Opening     float64 `json:"opening"`
	Accrued     float64 `json:"accrued"`
	Used        float64 `json:"used"`
	Adjusted    float64 `json:"adjusted"`
	Total       float64 `json:"total"`
	Available   float64 `json:"available"`
}

// CalculateLeaveBalances calculates leave balances using map-based approach
// This function takes leave types and balance records, then calculates the final balances
func CalculateLeaveBalances(leaveTypes []LeaveTypeData, balanceRecords []LeaveBalanceData) []CalculatedBalance {
	// Create a map of leave_type_id -> balance for O(1) lookup
	balanceMap := make(map[int]LeaveBalanceData)
	for _, balance := range balanceRecords {
		balanceMap[balance.LeaveTypeID] = balance
	}

	var calculatedBalances []CalculatedBalance

	// Calculate balances for each leave type
	for _, lt := range leaveTypes {
		balance, exists := balanceMap[lt.LeaveTypeID]

		var opening, accrued, used, adjusted, total, available float64

		if exists {
			// Balance record exists - use actual values from database
			opening = balance.Opening
			accrued = balance.Accrued
			used = balance.Used
			adjusted = balance.Adjusted
			// Total = Opening + Accrued
			total = opening + accrued
			// Available = Closing (which is calculated as: opening + accrued - used + adjusted)
			available = balance.Closing
		} else {
			// No balance record exists - use default entitlement
			opening = lt.DefaultEntitlement
			accrued = 0
			used = 0
			adjusted = 0
			// Total = Default Entitlement (treated as opening)
			total = lt.DefaultEntitlement
			// Available = Default Entitlement (since nothing used yet)
			available = lt.DefaultEntitlement
		}

		calculatedBalances = append(calculatedBalances, CalculatedBalance{
			LeaveTypeID: lt.LeaveTypeID,
			LeaveType:   lt.LeaveTypeName,
			Opening:     opening,
			Accrued:     accrued,
			Used:        used,
			Adjusted:    adjusted,
			Total:       total,
			Available:   available,
		})
	}

	return calculatedBalances
}
