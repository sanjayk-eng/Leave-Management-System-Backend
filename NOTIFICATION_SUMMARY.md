# Email Notification Integration - Summary

## What Was Added

### 1. New File: `utils/notification.go`
Complete email notification utility with the following functions:
- `SendEmail()` - Core function to send emails via Google Apps Script
- `SendEmployeeCreationEmail()` - Welcome email with credentials
- `SendLeaveApplicationEmail()` - Notify managers/admins of new leave requests
- `SendLeaveApprovalEmail()` - Notify employee of approved leave
- `SendLeaveRejectionEmail()` - Notify employee of rejected leave

### 2. Updated: `controllers/employee.go`
**Modified Function**: `CreateEmployee()`
- Added async email notification after successful employee creation
- Sends welcome email with login credentials to new employee
- Non-blocking implementation using goroutines

### 3. Updated: `controllers/leave.go`
**Modified Functions**:

#### `ApplyLeave()`
- Added notification to manager, admins, and superadmins
- Fetches all relevant recipients from database
- Sends leave application details for approval

#### `ActionLeave()`
- Added approval notification to employee
- Added rejection notification to employee
- Includes leave details in both scenarios

## Email Notification Flow

```
Employee Creation
    ↓
[Admin creates employee] → [Employee receives welcome email with password]

Leave Application
    ↓
[Employee applies] → [Manager + Admins + SuperAdmins receive notification]

Leave Approval/Rejection
    ↓
[Manager/Admin takes action] → [Employee receives approval/rejection email]
```

## Key Features

✅ **Asynchronous Processing**: All emails sent in background (non-blocking)
✅ **Error Handling**: Failed emails logged but don't break main flow
✅ **Multiple Recipients**: Supports sending to multiple admins/managers
✅ **Professional Templates**: Well-formatted email content
✅ **Google Apps Script Integration**: Uses provided email service URL

## Testing Commands

### Test Employee Creation Email
```bash
curl -X POST http://localhost:8080/api/employee/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john@zenithive.com",
    "role": "EMPLOYEE",
    "password": "welcome123",
    "salary": 50000,
    "joining_date": "2024-11-25T00:00:00Z"
  }'
```

### Test Leave Application Email
```bash
curl -X POST http://localhost:8080/api/leaves/apply \
  -H "Authorization: Bearer <employee_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "leave_type_id": 1,
    "start_date": "2024-12-01T00:00:00Z",
    "end_date": "2024-12-05T00:00:00Z"
  }'
```

### Test Approval Email
```bash
curl -X POST http://localhost:8080/api/leaves/<leave_id>/action \
  -H "Authorization: Bearer <manager_token>" \
  -H "Content-Type: application/json" \
  -d '{"action": "APPROVE"}'
```

### Test Rejection Email
```bash
curl -X POST http://localhost:8080/api/leaves/<leave_id>/action \
  -H "Authorization: Bearer <manager_token>" \
  -H "Content-Type: application/json" \
  -d '{"action": "REJECT"}'
```

## Email Service Configuration

**Google Apps Script URL**: 
```
https://script.google.com/macros/s/AKfycbxsNl0-rGsVKoszXmURHXoFuxjJeKJTRcC_7AdAA61N56ghaMwdto6RmIBdno4Hz0vQDA/exec
```

**Request Format**:
```json
{
  "to": "recipient@zenithive.com",
  "subject": "Email Subject",
  "body": "Email body content"
}
```

## Files Modified/Created

```
✨ NEW FILES:
├── utils/notification.go           (Email utility functions)
├── NOTIFICATION_SETUP.md           (Detailed setup guide)
└── NOTIFICATION_SUMMARY.md         (This file)

📝 MODIFIED FILES:
├── controllers/employee.go         (Added email on employee creation)
├── controllers/leave.go            (Added emails for leave actions)
└── API_DOCUMENTATION.md            (Added email notification section)
```

## Next Steps

1. **Test the integration**: Use the curl commands above to test each notification type
2. **Monitor logs**: Check console output for any email sending errors
3. **Verify Google Script**: Ensure the Apps Script URL is active and accessible
4. **Customize templates**: Modify email content in `utils/notification.go` if needed

## Notes

- All emails are sent asynchronously to avoid blocking API responses
- Email failures are logged but don't affect the main operation
- The system requires valid email addresses ending with `@zenithive.com`
- Passwords are sent in plain text (consider implementing password reset flow for production)
