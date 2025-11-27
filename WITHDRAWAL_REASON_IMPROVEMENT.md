# Withdrawal Reason - Database Storage Improvement

## 🎯 What Was Improved

Added database storage for withdrawal reasons in the `Tbl_Leave` table, so the reason is permanently saved and can be retrieved later.

---

## ✨ New Features

### 1. Reason Saved to Database
- Withdrawal reason is now stored in the `reason` column of `Tbl_Leave`
- Reason is permanently saved with the leave record
- Can be retrieved later for audit purposes

### 2. Default Reason
- If no reason provided, defaults to "Withdrawn by [ROLE]"
- Example: "Withdrawn by ADMIN" or "Withdrawn by MANAGER"

### 3. Reason in Response
- API response now includes `withdrawal_reason` field
- Frontend can display the reason immediately

---

## 📊 Before vs After

### Before ❌
```json
// Reason only sent in email, not saved
POST /api/leaves/:id/withdraw
{
  "reason": "Emergency"
}

// Response - no reason field
{
  "message": "leave withdrawn successfully",
  "leave_id": "uuid",
  "days_restored": 5
}

// Database - reason not saved
```

### After ✅
```json
// Reason saved to database
POST /api/leaves/:id/withdraw
{
  "reason": "Emergency project requirement"
}

// Response - includes reason
{
  "message": "leave withdrawn successfully and balance restored",
  "leave_id": "uuid",
  "days_restored": 5,
  "withdrawal_by": "admin_uuid",
  "withdrawal_role": "ADMIN",
  "withdrawal_reason": "Emergency project requirement"
}

// Database - reason saved in Tbl_Leave.reason column
```

---

## 🔄 Implementation Details

### Database Update
```sql
UPDATE Tbl_Leave 
SET status='WITHDRAWN', 
    reason='Emergency project requirement',  -- ✅ NEW: Reason saved
    updated_at=NOW() 
WHERE id=$1
```

### Default Reason Logic
```go
withdrawalReason := input.Reason
if withdrawalReason == "" {
    withdrawalReason = fmt.Sprintf("Withdrawn by %s", role)
}
```

### Response Enhancement
```go
c.JSON(200, gin.H{
    "message":           "leave withdrawn successfully and balance restored",
    "leave_id":          leaveID,
    "days_restored":     leave.Days,
    "withdrawal_by":     currentUserID,
    "withdrawal_role":   role,
    "withdrawal_reason": withdrawalReason,  // ✅ NEW: Reason in response
})
```

---

## 🎯 Use Cases

### Use Case 1: Admin Withdraws with Reason
```bash
curl -X POST http://localhost:8080/api/leaves/LEAVE_ID/withdraw \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Critical production issue requires all hands on deck"
  }'

# Response
{
  "message": "leave withdrawn successfully and balance restored",
  "leave_id": "uuid",
  "days_restored": 5,
  "withdrawal_by": "admin_uuid",
  "withdrawal_role": "ADMIN",
  "withdrawal_reason": "Critical production issue requires all hands on deck"
}

# Database: reason saved ✅
```

---

### Use Case 2: Manager Withdraws Without Reason
```bash
curl -X POST http://localhost:8080/api/leaves/LEAVE_ID/withdraw \
  -H "Authorization: Bearer <manager_token>" \
  -H "Content-Type: application/json" \
  -d '{}'

# Response
{
  "message": "leave withdrawn successfully and balance restored",
  "leave_id": "uuid",
  "days_restored": 3,
  "withdrawal_by": "manager_uuid",
  "withdrawal_role": "MANAGER",
  "withdrawal_reason": "Withdrawn by MANAGER"  // ✅ Default reason
}

# Database: default reason saved ✅
```

---

### Use Case 3: View Leave History with Reason
```bash
# Get all leaves
GET /api/leaves/all

# Response includes withdrawn leaves with reasons
{
  "total": 10,
  "data": [
    {
      "id": "uuid",
      "employee": "John Doe",
      "leave_type": "Annual Leave",
      "start_date": "2024-12-01",
      "end_date": "2024-12-05",
      "days": 5,
      "status": "WITHDRAWN",
      "reason": "Emergency project requirement",  // ✅ Reason visible
      "applying_date": "2024-11-20"
    }
  ]
}
```

---

## 📧 Email Notification

Email already includes the reason (no changes needed):

**Subject**: Leave Request Withdrawn

**Body**:
```
Dear John Doe,

Your approved leave request has been withdrawn by Admin Smith (ADMIN).

Leave Type: Annual Leave
Start Date: 2024-12-01
End Date: 2024-12-05
Duration: 5.0 days
Status: WITHDRAWN
Reason: Emergency project requirement  ✅

Your leave balance has been restored. The 5.0 days have been credited back to your account.

Best regards,
Zenithive Leave Management System
```

---

## 🔍 Audit Trail Benefits

### Before ❌
- Reason only in email (can be deleted)
- No permanent record
- Can't retrieve reason later
- Difficult to audit

### After ✅
- Reason in database (permanent)
- Always available for retrieval
- Easy to audit
- Complete history

---

## 💻 Frontend Integration

### Display Withdrawal Reason
```javascript
const LeaveCard = ({ leave }) => {
  return (
    <div className="leave-card">
      <h3>{leave.leave_type}</h3>
      <p>Status: {leave.status}</p>
      
      {/* Show reason if withdrawn */}
      {leave.status === 'WITHDRAWN' && leave.reason && (
        <div className="withdrawal-reason">
          <strong>Withdrawal Reason:</strong>
          <p>{leave.reason}</p>
        </div>
      )}
      
      <p>Dates: {leave.start_date} to {leave.end_date}</p>
      <p>Days: {leave.days}</p>
    </div>
  );
};
```

### Withdraw Leave with Reason
```javascript
const withdrawLeave = async (leaveId) => {
  const reason = prompt('Enter reason for withdrawal:');
  
  const response = await fetch(`/api/leaves/${leaveId}/withdraw`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ 
      reason: reason || '' // Optional
    })
  });
  
  const data = await response.json();
  
  alert(`Leave withdrawn successfully!
Reason: ${data.withdrawal_reason}
Days restored: ${data.days_restored}`);
};
```

---

## 📊 Database Schema

### Tbl_Leave Table
```sql
CREATE TABLE Tbl_Leave (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    leave_type_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    days FLOAT NOT NULL,
    status VARCHAR(20) NOT NULL,
    reason TEXT DEFAULT '',  -- ✅ Stores withdrawal/cancellation reason
    applied_by UUID,
    approved_by UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

---

## ✅ Benefits

### Transparency
✅ Clear record of why leave was withdrawn  
✅ Helps employees understand decisions  
✅ Reduces confusion and disputes  

### Audit Trail
✅ Permanent record in database  
✅ Can be retrieved anytime  
✅ Useful for HR reviews  
✅ Compliance and reporting  

### Better Communication
✅ Reason visible in leave history  
✅ Included in email notifications  
✅ Available in API responses  

### Default Handling
✅ Never empty - always has a reason  
✅ Default shows who withdrew it  
✅ No null/empty values  

---

## 🧪 Testing

### Test 1: Withdraw with Reason
```bash
POST /api/leaves/:id/withdraw
Body: { "reason": "Emergency" }

Expected:
- Database: reason = "Emergency"
- Response: withdrawal_reason = "Emergency"
- Email: includes "Emergency"
```

### Test 2: Withdraw without Reason
```bash
POST /api/leaves/:id/withdraw
Body: {}

Expected:
- Database: reason = "Withdrawn by ADMIN"
- Response: withdrawal_reason = "Withdrawn by ADMIN"
- Email: includes "Withdrawn by ADMIN"
```

### Test 3: View Leave History
```bash
GET /api/leaves/all

Expected:
- Withdrawn leaves show reason field
- Reason is visible in response
```

---

## 📁 Files Modified

1. ✅ `controllers/leave.go` - Updated WithdrawLeave to save reason
2. ✅ `LEAVE_CANCEL_WITHDRAW.md` - Updated documentation
3. ✅ `WITHDRAWAL_REASON_IMPROVEMENT.md` - This document

---

## ✅ Summary

### What Changed
✅ Withdrawal reason now saved to database  
✅ Default reason if none provided  
✅ Reason included in API response  
✅ Reason visible in leave history  

### Benefits
✅ Permanent audit trail  
✅ Better transparency  
✅ Improved communication  
✅ Complete history  

### Database
✅ Uses existing `reason` column  
✅ No migration needed  
✅ Always has a value (never empty)  

---

**Updated**: November 27, 2024  
**Status**: ✅ COMPLETE & IMPROVED
