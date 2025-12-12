package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/models"
)

// GetLogs - only for super_admin to get logs filtered by days
func (h *HandlerFunc) GetLogs(c *gin.Context) {
	// Check if user is SUPERADMIN
	role, exists := c.Get("role")
	if !exists || role != "SUPERADMIN" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied. Only SUPERADMIN can view logs",
		})
		return
	}

	// Get days parameter from query (default to 7 days if not provided or empty)
	daysParam := c.Query("days")
	days := 7 // Default value

	// If parameter is provided and not empty, try to parse it
	if daysParam != "" {
		parsedDays, err := strconv.Atoi(daysParam)
		if err != nil || parsedDays < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid days parameter. Must be a positive integer",
			})
			return
		}
		days = parsedDays
	}

	// Calculate the date threshold
	dateThreshold := time.Now().AddDate(0, 0, -days)

	// Query to get logs with user names, filtered by days
	query := `
		SELECT 
			l.id,
			e.full_name as user_name,
			l.action,
			l.component,
			l.created_at
		FROM tbl_log l
		JOIN Tbl_Employee e ON l.from_user_id = e.id
		WHERE l.created_at >= $1
		ORDER BY l.created_at DESC
	`

	rows, err := h.Query.DB.Query(query, dateThreshold)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch logs",
		})
		return
	}
	defer rows.Close()

	var logs []models.LogResponse
	for rows.Next() {
		var log models.LogResponse
		err := rows.Scan(
			&log.ID,
			&log.UserName,
			&log.Action,
			&log.Component,
			&log.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to scan log data",
			})
			return
		}
		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error iterating through logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logs retrieved successfully",
		"data": gin.H{
			"logs":        logs,
			"total_count": len(logs),
			"days_filter": days,
			"date_from":   dateThreshold.Format("2006-01-02"),
		},
	})
}
