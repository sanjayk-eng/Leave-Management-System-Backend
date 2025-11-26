# Email Notification System Setup

## Overview
The system uses Google Apps Script as an email service to send automated notifications for various events.

## Features Implemented

### 1. Employee Creation Notification
- **Trigger**: When admin creates a new employee
- **Recipient**: New employee
- **Content**: Welcome message with login credentials
- **File**: `controllers/employee.go` - `CreateEmployee()` function

### 2. Leave Application Notification
- **Trigger**: When employee applies for leave
- **Recipients**: 
  - Employee's direct manager
  - All ADMIN users
  - All SUPERADMIN users
- **Content**: Leave request details for approval
- **File**: `controllers/leave.go` - `ApplyLeave()` function

### 3. Leave Approval Notification
- **Trigger**: When manager/admin approves leave
- **Recipient**: Employee who applied
- **Content**: Approval confirmation with leave details
- **File**: `controllers/leave.go` - `ActionLeave()` function

### 4. Leave Rejection Notification
- **Trigger**: When manager/admin rejects leave
- **Recipient**: Employee who applied
- **Content**: Rejection notice with leave details
- **File**: `controllers/leave.go` - `ActionLeave()` function

## Technical Implementation

### Email Service
- **Provider**: Google Apps Script
- **URL**: `https://script.google.com/macros/s/AKfycbxsNl0-rGsVKoszXmURHXoFuxjJeKJTRcC_7AdAA61N56ghaMwdto6RmIBdno4Hz0vQDA/exec`
- **Method**: HTTP POST
- **Format**: JSON

### Request Format
```json
{
  "to": "recipient@zenithive.com",
  "subject": "Email Subject",
  "body": "Email body content"
}
```

### Code Structure
```
utils/notification.go
├── SendEmail()                      // Core email sending function
├── SendEmployeeCreationEmail()      // Welcome email for new employees
├── SendLeaveApplicationEmail()      // Notify managers/admins of leave request
├── SendLeaveApprovalEmail()         // Notify employee of approval
└── SendLeaveRejectionEmail()        // Notify employee of rejection
```

## Async Processing
All email notifications are sent asynchronously using goroutines to ensure:
- API responses are not blocked
- Fast response times
- Non-critical email failures don't affect core operations

Example:
```go
go func() {
    if err := utils.SendEmployeeCreationEmail(email, name, password); err != nil {
        fmt.Printf("Failed to send email: %v\n", err)
    }
}()
```

## Error Handling
- Email failures are logged but don't block the main operation
- Failed emails are printed to console for monitoring
- HTTP timeout set to 10 seconds to prevent hanging

## Testing Email Notifications

### 1. Test Employee Creation Email
```bash
curl -X POST http://localhost:8080/api/employee/ \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "email": "testuser@zenithive.com",
    "role": "EMPLOYEE",
    "password": "test123",
    "salary": 40000,
    "joining_date": "2024-11-25T00:00:00Z"
  }'
```
**Expected**: Welcome email sent to `testuser@zenithive.com`

### 2. Test Leave Application Email
```bash
# First login as employee
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@zenithive.com",
    "password": "password123"
  }'

# Apply for leave
curl -X POST http://localhost:8080/api/leaves/apply \
  -H "Authorization: Bearer <employee_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "leave_type_id": 1,
    "start_date": "2024-12-01T00:00:00Z",
    "end_date": "2024-12-03T00:00:00Z"
  }'
```
**Expected**: Notification emails sent to manager, all admins, and superadmins

### 3. Test Leave Approval Email
```bash
curl -X POST http://localhost:8080/api/leaves/<leave_id>/action \
  -H "Authorization: Bearer <manager_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "APPROVE"
  }'
```
**Expected**: Approval email sent to employee

### 4. Test Leave Rejection Email
```bash
curl -X POST http://localhost:8080/api/leaves/<leave_id>/action \
  -H "Authorization: Bearer <manager_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "REJECT"
  }'
```
**Expected**: Rejection email sent to employee

## Monitoring
Check console logs for email status:
```
Failed to send welcome email to user@zenithive.com: <error>
Failed to send leave application notifications: <error>
Failed to send approval email: <error>
Failed to send rejection email: <error>
```

## Customization

### Modify Email Templates
Edit the email body in `utils/notification.go`:

```go
body := fmt.Sprintf(`
Your custom email template here
Employee: %s
...
`, variables)
```

### Change Email Service URL
Update the constant in `utils/notification.go`:
```go
const GOOGLE_SCRIPT_URL = "your-new-url"
```

### Add New Notification Types
1. Create a new function in `utils/notification.go`
2. Call it from the appropriate controller
3. Use goroutine for async execution

Example:
```go
// In utils/notification.go
func SendPayrollNotification(email, name string, amount float64) error {
    subject := "Payroll Processed"
    body := fmt.Sprintf("Dear %s, Your salary of %.2f has been processed.", name, amount)
    return SendEmail(email, subject, body)
}

// In controller
go func() {
    utils.SendPayrollNotification(email, name, amount)
}()
```

## Troubleshooting

### Emails Not Being Sent
1. Check if Google Apps Script URL is accessible
2. Verify JSON payload format
3. Check network connectivity
4. Review console logs for errors

### Emails Delayed
- Emails are sent asynchronously, slight delays are normal
- Check Google Apps Script execution logs
- Verify email service quota limits

### Wrong Recipients
- Verify role assignments in database
- Check manager assignments
- Review SQL queries in notification code

## Security Considerations
- Passwords are sent in plain text via email (consider adding password reset flow)
- Email content is not encrypted in transit (relies on HTTPS)
- No rate limiting on email sending (consider adding if needed)

## Future Enhancements
- [ ] HTML email templates
- [ ] Email queue system for retry logic
- [ ] Email delivery status tracking
- [ ] Configurable email templates via admin panel
- [ ] Email preferences per user
- [ ] Digest emails (daily/weekly summaries)
- [ ] SMS notifications integration
