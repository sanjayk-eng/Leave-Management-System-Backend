# Quick Start - Email Notifications

## ⚡ What's New?

Your backend now automatically sends email notifications for:
1. 📧 **Employee Creation** → Welcome email with credentials
2. 📧 **Leave Application** → Notify manager + admins
3. 📧 **Leave Approval** → Notify employee
4. 📧 **Leave Rejection** → Notify employee

---

## 🚀 Quick Test (5 Minutes)

### Step 1: Start Your Server
```bash
go run main.go
```

### Step 2: Create an Employee (Test Welcome Email)
```bash
curl -X POST http://localhost:8080/api/employee/ \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test Employee",
    "email": "test@zenithive.com",
    "role": "EMPLOYEE",
    "password": "welcome123",
    "salary": 50000,
    "joining_date": "2024-11-26T00:00:00Z"
  }'
```
✅ Check `test@zenithive.com` inbox for welcome email

### Step 3: Apply for Leave (Test Leave Application Email)
```bash
# First, login as the employee
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@zenithive.com", "password": "welcome123"}'

# Use the token from response to apply for leave
curl -X POST http://localhost:8080/api/leaves/apply \
  -H "Authorization: Bearer EMPLOYEE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "leave_type_id": 1,
    "start_date": "2024-12-01T00:00:00Z",
    "end_date": "2024-12-03T00:00:00Z"
  }'
```
✅ Check manager and admin inboxes for leave application notification

### Step 4: Approve Leave (Test Approval Email)
```bash
curl -X POST http://localhost:8080/api/leaves/LEAVE_ID/action \
  -H "Authorization: Bearer MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"action": "APPROVE"}'
```
✅ Check `test@zenithive.com` inbox for approval email

---

## 📋 What Was Changed?

### New File
- `utils/notification.go` - Email sending functions

### Modified Files
- `controllers/employee.go` - Added welcome email on creation
- `controllers/leave.go` - Added notifications for leave actions

### Documentation
- `API_DOCUMENTATION.md` - Updated with email notification section
- `NOTIFICATION_SETUP.md` - Detailed setup guide
- `NOTIFICATION_FLOW.md` - Visual flow diagrams
- `CHANGES_SUMMARY.md` - Complete implementation summary

---

## 🔍 How to Monitor

Watch your console for email logs:
```bash
# Successful sends = no output
# Failed sends will show:
Failed to send welcome email to user@zenithive.com: <error>
```

---

## ⚙️ Configuration

Email service is pre-configured with Google Apps Script:
```
URL: https://script.google.com/macros/s/AKfycbxsNl0-rGsVKoszXmURHXoFuxjJeKJTRcC_7AdAA61N56ghaMwdto6RmIBdno4Hz0vQDA/exec
```

No additional setup required! 🎉

---

## 🎯 Key Features

✅ **Non-Blocking** - Emails sent asynchronously  
✅ **Fault Tolerant** - Email failures don't break API  
✅ **Multi-Recipient** - Sends to multiple users  
✅ **Auto-Retry** - Google Apps Script handles retries  
✅ **Logged Errors** - Failed sends logged to console  

---

## 📚 Need More Info?

- **Setup Details**: See `NOTIFICATION_SETUP.md`
- **Flow Diagrams**: See `NOTIFICATION_FLOW.md`
- **Complete Changes**: See `CHANGES_SUMMARY.md`
- **API Docs**: See `API_DOCUMENTATION.md`

---

## ✅ Verification

Build and verify everything works:
```bash
go build -o test.exe
# Should compile without errors ✅
del test.exe
```

---

**Status**: ✅ Ready to Use  
**Version**: 1.0  
**Date**: November 26, 2024
