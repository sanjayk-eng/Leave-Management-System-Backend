# Email Notification Troubleshooting Guide

## 🔍 Issue: "unsupported protocol scheme" Error

### Error Message
```
Failed to send withdrawal email: failed to send email: Post "": unsupported protocol scheme ""
```

### Root Cause
The `GOOGLE_SCRIPT_URL` environment variable was being read at package initialization time (before the .env file was loaded), resulting in an empty string.

---

## ✅ Fix Applied

### Changed in `utils/notification.go`

**Before (❌ Wrong):**
```go
var GOOGLE_SCRIPT_URL = os.Getenv("GOOGLE_SCRIPT_URL") // Read at package init

func SendEmail(to, subject, body string) error {
    // ... code ...
    resp, err := client.Post(GOOGLE_SCRIPT_URL, ...) // Uses empty string
}
```

**After (✅ Correct):**
```go
func SendEmail(to, subject, body string) error {
    // Read at runtime (after .env is loaded)
    googleScriptURL := os.Getenv("GOOGLE_SCRIPT_URL")
    
    // Check if URL is set
    if googleScriptURL == "" {
        return fmt.Errorf("GOOGLE_SCRIPT_URL environment variable is not set")
    }
    
    // ... code ...
    resp, err := client.Post(googleScriptURL, ...) // Uses actual URL
}
```

---

## 🧪 Testing Email Notifications

### 1. Verify Environment Variable

**Check .env file:**
```bash
cat .env | grep GOOGLE_SCRIPT_URL
```

**Expected output:**
```
GOOGLE_SCRIPT_URL=https://script.google.com/macros/s/YOUR_SCRIPT_ID/exec
```

### 2. Test Email Sending

**Restart your server:**
```bash
# Stop the server (Ctrl+C)
# Start again
go run main.go
```

### 3. Trigger a Notification

**Apply for leave:**
```bash
curl -X POST http://localhost:8082/api/leaves/apply \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "leave_type_id": 1,
    "start_date": "2024-12-01T00:00:00Z",
    "end_date": "2024-12-03T00:00:00Z",
    "reason": "Personal work that needs attention"
  }'
```

**Check server logs:**
```
✅ Success: No error messages
❌ Error: "Failed to send email: ..." in logs
```

---

## 📧 Email Notification Events

### 1. Leave Application
**Trigger**: Employee applies for leave  
**Recipients**: Manager + All Admins  
**Function**: `SendLeaveApplicationEmail()`  

### 2. Leave Approval
**Trigger**: Manager/Admin approves leave  
**Recipients**: Employee  
**Function**: `SendLeaveApprovalEmail()`  

### 3. Leave Rejection
**Trigger**: Manager/Admin rejects leave  
**Recipients**: Employee  
**Function**: `SendLeaveRejectionEmail()`  

### 4. Leave Cancellation
**Trigger**: Employee cancels pending leave  
**Recipients**: Employee  
**Function**: `SendLeaveCancellationEmail()`  

### 5. Leave Withdrawal
**Trigger**: Admin withdraws approved leave  
**Recipients**: Employee  
**Function**: `SendLeaveWithdrawalEmail()`  

### 6. Employee Creation
**Trigger**: Admin creates new employee  
**Recipients**: New employee  
**Function**: `SendEmployeeCreationEmail()`  

### 7. Password Update
**Trigger**: Admin/HR updates employee password  
**Recipients**: Employee  
**Function**: `SendPasswordUpdateEmail()`  

### 8. Leave Added by Admin
**Trigger**: Admin/Manager adds leave for employee  
**Recipients**: Employee  
**Function**: `SendLeaveAddedByAdminEmail()`  

---

## 🔧 Common Issues & Solutions

### Issue 1: Empty GOOGLE_SCRIPT_URL

**Symptom:**
```
Failed to send email: GOOGLE_SCRIPT_URL environment variable is not set
```

**Solution:**
1. Check `.env` file exists in project root
2. Verify `GOOGLE_SCRIPT_URL` is set in `.env`
3. Restart the server

---

### Issue 2: Invalid Google Script URL

**Symptom:**
```
Failed to send email: Post "https://...": dial tcp: lookup script.google.com: no such host
```

**Solution:**
1. Verify the Google Script URL is correct
2. Check internet connection
3. Ensure Google Apps Script is deployed

---

### Issue 3: Google Script Returns Error

**Symptom:**
```
email service returned status: 400
```

**Solution:**
1. Check Google Apps Script logs
2. Verify script accepts POST requests
3. Ensure JSON format is correct:
```json
{
  "to": "email@example.com",
  "subject": "Subject",
  "body": "Body"
}
```

---

### Issue 4: Email Not Received

**Symptom:** No error in logs, but email not received

**Possible Causes:**
1. Email in spam folder
2. Google Apps Script quota exceeded
3. Invalid recipient email address
4. Gmail blocking the sender

**Solution:**
1. Check spam folder
2. Check Google Apps Script execution logs
3. Verify recipient email is correct
4. Check Gmail settings

---

## 🔐 Google Apps Script Setup

### Required Script

```javascript
function doPost(e) {
  try {
    // Parse the JSON payload
    var data = JSON.parse(e.postData.contents);
    
    // Send email
    GmailApp.sendEmail(
      data.to,
      data.subject,
      data.body
    );
    
    // Return success
    return ContentService.createTextOutput(
      JSON.stringify({ success: true })
    ).setMimeType(ContentService.MimeType.JSON);
    
  } catch (error) {
    // Return error
    return ContentService.createTextOutput(
      JSON.stringify({ 
        success: false, 
        error: error.toString() 
      })
    ).setMimeType(ContentService.MimeType.JSON);
  }
}
```

### Deployment Steps

1. Go to https://script.google.com
2. Create new project
3. Paste the script above
4. Click "Deploy" → "New deployment"
5. Select type: "Web app"
6. Execute as: "Me"
7. Who has access: "Anyone"
8. Click "Deploy"
9. Copy the web app URL
10. Add to `.env` file as `GOOGLE_SCRIPT_URL`

---

## 📊 Notification Flow Diagram

```
Employee Action
    ↓
Backend Controller
    ↓
Commit Transaction
    ↓
Async Goroutine (go func())
    ↓
Fetch Email Details from DB
    ↓
Call utils.SendEmail()
    ↓
Read GOOGLE_SCRIPT_URL (runtime)
    ↓
POST to Google Apps Script
    ↓
Google Apps Script sends email
    ↓
Email delivered to recipient
```

---

## 🧪 Testing Checklist

### Environment Setup
- [ ] `.env` file exists in project root
- [ ] `GOOGLE_SCRIPT_URL` is set in `.env`
- [ ] Google Apps Script is deployed
- [ ] Server restarted after `.env` changes

### Email Notifications
- [ ] Leave application sends email to manager/admin
- [ ] Leave approval sends email to employee
- [ ] Leave rejection sends email to employee
- [ ] Leave cancellation sends email to employee
- [ ] Leave withdrawal sends email to employee
- [ ] Employee creation sends email to new employee
- [ ] Password update sends email to employee

### Error Handling
- [ ] Empty URL returns clear error message
- [ ] Invalid URL returns clear error message
- [ ] Network errors are logged
- [ ] Email failures don't crash the application

---

## 📝 Debugging Tips

### 1. Enable Detailed Logging

Add this to your notification functions:
```go
fmt.Printf("📧 Sending email to: %s\n", to)
fmt.Printf("📧 Subject: %s\n", subject)
fmt.Printf("📧 Google Script URL: %s\n", googleScriptURL)
```

### 2. Test Email Function Directly

Create a test endpoint:
```go
func (h *HandlerFunc) TestEmail(c *gin.Context) {
    err := utils.SendEmail(
        "test@example.com",
        "Test Email",
        "This is a test email from the system",
    )
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"message": "Email sent successfully"})
}
```

### 3. Check Google Apps Script Logs

1. Go to https://script.google.com
2. Open your project
3. Click "Executions" in left sidebar
4. Check for errors or successful executions

---

## ✅ Verification

After applying the fix, you should see:

**✅ Success:**
```
Leave applied successfully
(No error messages in logs)
```

**❌ Before Fix:**
```
Failed to send email: Post "": unsupported protocol scheme ""
```

**✅ After Fix:**
```
(Email sent successfully, no errors)
```

---

## 📁 Files Modified

1. ✅ `utils/notification.go` - Fixed GOOGLE_SCRIPT_URL loading
2. ✅ `EMAIL_NOTIFICATION_TROUBLESHOOTING.md` - This guide

---

## 🎯 Summary

**Problem**: Environment variable read at wrong time  
**Solution**: Read environment variable at runtime  
**Result**: Emails now send successfully  

---

**Updated**: November 27, 2024  
**Status**: ✅ FIXED
