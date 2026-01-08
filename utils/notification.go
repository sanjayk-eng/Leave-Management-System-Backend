package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// GetSMTPConfig reads SMTP configuration from environment variables
func GetSMTPConfig() (*SMTPConfig, error) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || portStr == "" || username == "" || password == "" || from == "" {
		return nil, fmt.Errorf("missing SMTP configuration: ensure SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM are set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	return &SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
	}, nil
}

// SendEmail sends an email using SMTP
func SendEmail(to, subject, body string) error {
	config, err := GetSMTPConfig()
	if err != nil {
		return fmt.Errorf("SMTP configuration error: %v", err)
	}

	fmt.Printf("Attempting to send email to: %s with subject: %s\n", to, subject)

	// Create message
	message := fmt.Sprintf("From: %s\r\n", config.From)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// Setup authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Use smtp.SendMail for simpler, more reliable SMTP handling
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	fmt.Printf("SMTP Config - Host: %s, Port: %d, From: %s\n", config.Host, config.Port, config.From)
	
	err = smtp.SendMail(addr, auth, config.From, []string{to}, []byte(message))
	if err != nil {
		detailedErr := fmt.Errorf("SMTP send failed - Host: %s, Port: %d, To: %s, Error: %v", 
			config.Host, config.Port, to, err)
		fmt.Printf("SMTP ERROR: %v\n", detailedErr)
		fmt.Printf("Troubleshooting: Check SMTP credentials, network connectivity, and firewall settings\n")
		return detailedErr
	}

	fmt.Printf("Email sent successfully to: %s\n", to)
	return nil
}

// SendEmailToMultiple sends email to multiple recipients
func SendEmailToMultiple(recipients []string, subject, body string) error {
	config, err := GetSMTPConfig()
	if err != nil {
		return fmt.Errorf("SMTP configuration error: %v", err)
	}

	if len(recipients) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	fmt.Printf("Attempting to send email to %d recipients with subject: %s\n", len(recipients), subject)

	// Create message with multiple recipients
	message := fmt.Sprintf("From: %s\r\n", config.From)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ", "))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// Setup authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Use smtp.SendMail for simpler, more reliable SMTP handling
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	err = smtp.SendMail(addr, auth, config.From, recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email to multiple recipients: %v", err)
	}

	fmt.Printf("Email sent successfully to %d recipients\n", len(recipients))
	return nil
}

// SendEmployeeCreationEmail sends notification to newly created employee
func SendEmployeeCreationEmail(employeeEmail, employeeName, password string) error {
	subject := "Welcome to Zenithive - Your Account Has Been Created"
	body := fmt.Sprintf(`
Dear %s,

Welcome to Zenithive!

Your employee account has been successfully created. Below are your login credentials:

Email: %s
Password: %s

Please login to the system and change your password at your earliest convenience.

Login URL: [https://zenithiveapp.netlify.app]

If you have any questions, please contact your HR department.

Best regards,
Zenithive HR Team
`, employeeName, employeeEmail, password)

	return SendEmail(employeeEmail, subject, body)
}

// SendLeaveApplicationEmail sends notification to manager, admin, and superadmin
func SendLeaveApplicationEmail(recipients []string, employeeName, leaveType, startDate, endDate string, days float64, reason string) error {
	subject := fmt.Sprintf("Leave Application - %s", employeeName)
	body := fmt.Sprintf(`
Dear Manager/Admin,

A new leave application has been submitted and requires your review.

Employee: %s
Leave Type: %s
Start Date: %s
End Date: %s
Duration: %.1f days
Reason: %s
Status: Pending Approval

Please login to the system to approve or reject this leave request.

Best regards,
Zenithive Leave Management System
`, employeeName, leaveType, startDate, endDate, days, reason)

	return SendEmailToMultiple(recipients, subject, body)
}

func SendLeaveManagerRejectionEmail(
	AdminEmail []string,
	empEmail string,
	employeeName, leaveType, startDate, endDate string,
	days float64, rejectedBy string,
) error {

	subject := "Leave Request - Manager Rejection (Pending Final Decision)"

	// ------------------------
	// EMPLOYEE EMAIL (Step 1)
	// ------------------------
	empBody := fmt.Sprintf(`
Dear %s,

Your leave request has been REJECTED by your manager %s.

This is a first-level rejection. The request has now been forwarded to Admin/SuperAdmin for final review.

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: MANAGER_REJECTED

For more information, please contact your manager.

Best regards,
Zenithive Leave Management System
`, employeeName, rejectedBy, leaveType, startDate, endDate, days)

	if err := SendEmail(empEmail, subject, empBody); err != nil {
		return err
	}

	// ------------------------
	// ADMIN EMAIL (Step 1)
	// ------------------------
	adminBody := fmt.Sprintf(`
Dear Admin,

A leave request has been REJECTED at manager level by %s.

This leave now requires final rejection approval from Admin/SuperAdmin.

Leave Details:
- Employee: %s
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: MANAGER_REJECTED

Please log in to the admin panel to complete the final review.

Best regards,
Zenithive Leave Management System
`, rejectedBy, employeeName, leaveType, startDate, endDate, days)

	return SendEmailToMultiple(AdminEmail, subject, adminBody)
}

// SendLeaveManagerApprovalEmail sends notification for manager-level approval (first step)
func SendLeaveManagerApprovalEmail(
	AdminEmail []string,
	employeeEmail, employeeName, leaveType, startDate, endDate string,
	days float64, approvedBy string,
) error {

	subject := "Leave Approved by Manager"

	// ------------------------
	// 1) EMPLOYEE EMAIL
	// ------------------------
	empBody := fmt.Sprintf(`
Dear %s,

Your leave application has been APPROVED by your manager %s.

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: MANAGER APPROVED

Note: Your leave is now pending final approval from Admin/SuperAdmin.

Best regards,
Zenithive Leave Management System
`, employeeName, approvedBy, leaveType, startDate, endDate, days)

	if err := SendEmail(employeeEmail, subject, empBody); err != nil {
		return err
	}

	// ------------------------
	// 2) ADMIN EMAIL
	// ------------------------
	adminBody := fmt.Sprintf(`
Dear Admin,

A leave request has been APPROVED by the manager %s.

Leave Details:
- Employee: %s
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: MANAGER APPROVED

Please review and take final action.

Best regards,
Zenithive Leave Management System
`, approvedBy, employeeName, leaveType, startDate, endDate, days)

	return SendEmailToMultiple(AdminEmail, subject, adminBody)
}

// SendLeaveApprovalEmail sends notification to employee when leave is approved
func SendLeaveFinalApprovalEmail(
	AdminEmail []string,
	employeeEmail, employeeName, leaveType, startDate, endDate string,
	days float64, approvedBy string,
) error {

	subject := "Leave Approved"

	// ------------------------
	// 1) EMPLOYEE EMAIL
	// ------------------------
	empBody := fmt.Sprintf(`
Dear %s,

Your leave application has been APPROVED by %s.

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: APPROVED

Enjoy your time off!

Best regards,
Zenithive Leave Management System
`, employeeName, approvedBy, leaveType, startDate, endDate, days)

	if err := SendEmail(employeeEmail, subject, empBody); err != nil {
		return err
	}

	// ------------------------
	// 2) ADMIN EMAIL
	// ------------------------
	adminBody := fmt.Sprintf(`
Dear Admin,

The leave request for employee %s has been APPROVED by %s.

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: APPROVED

Best regards,
Zenithive Leave Management System
`, employeeName, approvedBy, leaveType, startDate, endDate, days)

	return SendEmailToMultiple(AdminEmail, subject, adminBody)
}

// SendLeaveRejectionEmail sends notification to employee when leave is rejected
func SendLeaveRejectionEmail(
	AdminEmail []string,
	empEmail string,
	employeeName, leaveType, startDate, endDate string,
	days float64, rejectedBy string,
) error {

	subject := "Leave Request Rejected"

	// ------------------------
	// 1) EMPLOYEE EMAIL
	// ------------------------
	empBody := fmt.Sprintf(`
Dear %s,

We regret to inform you that your leave application has been REJECTED by %s.

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: REJECTED

Please contact your manager if you require more information.

Best regards,
Zenithive Leave Management System
`, employeeName, rejectedBy, leaveType, startDate, endDate, days)

	if err := SendEmail(empEmail, subject, empBody); err != nil {
		return err
	}

	// ------------------------
	// 2) ADMIN EMAIL
	// ------------------------
	adminBody := fmt.Sprintf(`
Dear Admin,

A leave request has been REJECTED by %s.

Leave Details:
- Employee: %s
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: REJECTED

Best regards,
Zenithive Leave Management System
`, rejectedBy, employeeName, leaveType, startDate, endDate, days)

	return SendEmailToMultiple(AdminEmail, subject, adminBody)
}

// SendLeaveAddedByAdminEmail sends notification to employee when admin/manager adds leave on their behalf
func SendLeaveAddedByAdminEmail(employeeEmail, employeeName, leaveType, startDate, endDate string, days float64, addedBy, addedByRole string) error {
	subject := fmt.Sprintf("Leave Added to Your Account - %s", leaveType)
	body := fmt.Sprintf(`
Dear %s,

A leave has been added to your account by %s (%s).

Leave Type: %s
Start Date: %s
End Date: %s
Duration: %.1f days
Status: APPROVED

This leave has been automatically approved and your leave balance has been updated accordingly.

If you have any questions about this leave entry, please contact your manager or HR department.

Best regards,
Zenithive Leave Management System
`, employeeName, addedBy, addedByRole, leaveType, startDate, endDate, days)

	return SendEmail(employeeEmail, subject, body)
}

// SendPasswordUpdateEmail sends notification to employee when their password is updated by admin
func SendPasswordUpdateEmail(employeeEmail, employeeName, newPassword, updatedByEmail, updatedByRole string) error {
	subject := "Your Password Has Been Updated"
	body := fmt.Sprintf(`
Dear %s,

Your account password has been updated by %s (%s).

Your new login credentials are:
Email: %s
Password: %s

If you did not request this change, please contact your HR department immediately.

For security reasons, we recommend:
1. Login with your new password
2. Change your password to something memorable
3. Keep your password secure and do not share it with anyone

Login URL: [https://zenithiveapp.netlify.app]

Best regards,
Zenithive HR Team
`, employeeName, updatedByEmail, updatedByRole, employeeEmail, newPassword)

	return SendEmail(employeeEmail, subject, body)
}

// SendLeaveCancellationEmail sends notification when leave is cancelled
func SendLeaveCancellationEmail(employeeEmail, employeeName, leaveType, startDate, endDate string, days float64) error {
	subject := "Leave Request Cancelled"
	body := fmt.Sprintf(`
Dear %s,

Your leave request has been cancelled.

Leave Type: %s
Start Date: %s
End Date: %s
Duration: %.1f days
Status: CANCELLED

If you did not cancel this leave request, please contact your manager or HR department immediately.

Best regards,
Zenithive Leave Management System
`, employeeName, leaveType, startDate, endDate, days)

	return SendEmail(employeeEmail, subject, body)
}

// SendLeaveWithdrawalPendingEmail sends notification to admins when manager requests withdrawal
func SendLeaveWithdrawalPendingEmail(recipients []string, employeeName, leaveType, startDate, endDate string, days float64, requestedBy, reason string) error {
	subject := fmt.Sprintf("Leave Withdrawal Request - %s", employeeName)

	reasonText := ""
	if reason != "" {
		reasonText = fmt.Sprintf("\nReason: %s", reason)
	}

	body := fmt.Sprintf(`
Dear Admin,

A leave withdrawal request has been submitted and requires your approval.

Employee: %s
Leave Type: %s
Start Date: %s
End Date: %s
Duration: %.1f days
Requested By: %s (MANAGER)
Status: Pending Withdrawal Approval%s

Please login to the system to approve or reject this withdrawal request.

Best regards,
Zenithive Leave Management System
`, employeeName, leaveType, startDate, endDate, days, requestedBy, reasonText)

	return SendEmailToMultiple(recipients, subject, body)
}

// SendLeaveWithdrawalEmail sends notification when approved leave is withdrawn
// SendLeaveWithdrawalEmail sends notification when a leave is withdrawn
func SendLeaveWithdrawalEmail(
	adminEmails []string,
	employeeEmail, employeeName, leaveType, startDate, endDate string,
	days float64, withdrawnBy, withdrawnByRole, reason string,
) error {

	subject := "Leave Request Withdrawn"

	// Optional reason text
	reasonText := ""
	if reason != "" {
		reasonText = fmt.Sprintf("\nReason: %s", reason)
	}

	// ------------------------
	// 1) EMPLOYEE EMAIL
	// ------------------------
	empBody := fmt.Sprintf(`
Dear %s,

Your approved leave request has been WITHDRAWN by %s (%s).

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: WITHDRAWN%s

Your leave balance has been restored. The %.1f days have been credited back to your account.

If you have any questions regarding this withdrawal, please contact your manager or HR department.

Best regards,
Zenithive Leave Management System
`, employeeName, withdrawnBy, withdrawnByRole, leaveType, startDate, endDate, days, reasonText, days)

	if err := SendEmail(employeeEmail, subject, empBody); err != nil {
		return err
	}

	// ------------------------
	// 2) ADMIN EMAIL
	// ------------------------
	adminBody := fmt.Sprintf(`
Dear Admin,

The leave request of employee %s has been WITHDRAWN by %s (%s).

Leave Details:
- Leave Type: %s
- Start Date: %s
- End Date: %s
- Duration: %.1f days
- Status: WITHDRAWN%s

The employee's leave balance has been restored.

Best regards,
Zenithive Leave Management System
`, employeeName, withdrawnBy, withdrawnByRole, leaveType, startDate, endDate, days, reasonText)

	return SendEmailToMultiple(adminEmails, subject, adminBody)
}

// SendPayslipWithdrawalEmail sends notification when payslip is withdrawn
func SendPayslipWithdrawalEmail(employeeEmail, employeeName string, month, year int, netSalary float64, withdrawnBy, withdrawnByRole, reason string) error {
	monthNames := []string{"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December"}

	subject := fmt.Sprintf("Payslip Withdrawn - %s %d", monthNames[month], year)

	reasonText := ""
	if reason != "" {
		reasonText = fmt.Sprintf("\nReason: %s", reason)
	}

	body := fmt.Sprintf(`
Dear %s,

Your payslip for %s %d has been withdrawn by %s (%s).

Pay Period: %s %d
Net Salary: ₹%.2f
Status: WITHDRAWN%s

This payslip has been marked as withdrawn and may require reprocessing. Please contact your HR department or payroll administrator for more information.

If you have any questions about this withdrawal, please reach out to your manager or HR department.

Best regards,
Zenithive Payroll Management System
`, employeeName, monthNames[month], year, withdrawnBy, withdrawnByRole, monthNames[month], year, netSalary, reasonText)

	return SendEmail(employeeEmail, subject, body)
}
