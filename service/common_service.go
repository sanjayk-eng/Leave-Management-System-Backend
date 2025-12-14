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

// CalculateAbsentDaysForMonth calculates the number of absent days for a specific month
// Now handles timing-based leaves (half days vs full days) correctly
// Uses the pre-calculated days from leave records which include timing considerations
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
}
