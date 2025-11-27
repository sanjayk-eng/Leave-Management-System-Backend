# Leave Management - Complete Summary

## 🎉 New Features Added

Two new leave management endpoints have been implemented:

1. ✅ **DELETE /api/leaves/:id/cancel** - Cancel pending leave
2. ✅ **POST /api/leaves/:id/withdraw** - Withdraw approved leave

---

## 📊 Quick Comparison

| Feature | Cancel | Withdraw |
|---------|--------|----------|
| **Endpoint** | DELETE /:id/cancel | POST /:id/withdraw |
| **Applies To** | Pending leaves | Approved leaves |
| **Who Can Use** | Employee (own) or Admin | Admin or Manager |
| **Balance Impact** | None | Restores balance |
| **Reason Required** | No | Optional |

---

## 🔐 Permissions

### Cancel Leave
| Role | Permission |
|------|------------|
| SUPERADMIN | ✅ Cancel any pending leave |
| ADMIN | ✅ Cancel any pending leave |
| MANAGER | ❌ No |
| EMPLOYEE | ✅ Cancel own pending leave |

### Withdraw Leave
| Role | Permission |
|------|------------|
| SUPERADMIN | ✅ Withdraw any approved leave |
| ADMIN | ✅ Withdraw any approved leave |
| MANAGER | ✅ Withdraw team members' approved leave |
| EMPLOYEE | ❌ No |

---

## 🎯 Quick Examples

### Cancel Pending Leave
```bash
# Employee cancels own leave
curl -X DELETE http://localhost:8080/api/leaves/LEAVE_ID/cancel \
  -H "Authorization: Bearer <employee_token>"

# Response
{
  "message": "leave request cancelled successfully",
  "leave_id": "uuid"
}
```

### Withdraw Approved Leave
```bash
# Admin withdraws approved leave with reason
curl -X POST http://localhost:8080/api/leaves/LEAVE_ID/withdraw \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Emergency project requirement"
  }'

# Response
{
  "message": "leave withdrawn successfully and balance restored",
  "leave_id": "uuid",
  "days_restored": 5,
  "withdrawal_by": "admin_uuid",
  "withdrawal_role": "ADMIN"
}
```

---

## 📧 Email Notifications

### Cancel Notification
**Subject**: Leave Request Cancelled

**Content**:
- Leave details (type, dates, duration)
- Status: CANCELLED
- Security note

### Withdraw Notification
**Subject**: Leave Request Withdrawn

**Content**:
- Leave details (type, dates, duration)
- Who withdrew it and their role
- Reason (if provided)
- Balance restoration confirmation
- Status: WITHDRAWN

---

## 🔄 Leave Status Flow

```
Apply Leave
    ↓
┌─────────┐
│ Pending │ ──cancel──> CANCELLED
└────┬────┘
     │
  approve
     │
     ▼
┌──────────┐
│ APPROVED │ ──withdraw──> WITHDRAWN (balance restored)
└──────────┘
```

---

## ✨ Key Features

### Cancel Leave
✅ Employee self-service  
✅ Admin override capability  
✅ Only for pending leaves  
✅ Email notification  
✅ No balance impact  

### Withdraw Leave
✅ Admin/Manager control  
✅ Automatic balance restoration  
✅ Optional reason field  
✅ Manager hierarchy validation  
✅ Email notification with details  

---

## 🧪 Testing

### Test Cancel
```bash
# 1. Apply leave (as employee)
POST /api/leaves/apply

# 2. Cancel it (as same employee)
DELETE /api/leaves/:id/cancel

# Expected: Success, status = CANCELLED
```

### Test Withdraw
```bash
# 1. Apply leave (as employee)
POST /api/leaves/apply

# 2. Approve it (as manager)
POST /api/leaves/:id/action
Body: { "action": "APPROVE" }

# 3. Withdraw it (as admin/manager)
POST /api/leaves/:id/withdraw
Body: { "reason": "Emergency" }

# Expected: Success, status = WITHDRAWN, balance restored
```

---

## 📁 Files Modified

1. ✅ `controllers/leave.go` - Added CancelLeave() and WithdrawLeave()
2. ✅ `routes/router.go` - Added cancel and withdraw routes
3. ✅ `utils/notification.go` - Added email functions
4. ✅ `LEAVE_CANCEL_WITHDRAW.md` - Detailed documentation
5. ✅ `LEAVE_MANAGEMENT_SUMMARY.md` - This file

---

## 💻 Frontend Integration

```javascript
// Cancel leave
const cancelLeave = async (leaveId) => {
  const response = await fetch(`/api/leaves/${leaveId}/cancel`, {
    method: 'DELETE',
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return await response.json();
};

// Withdraw leave
const withdrawLeave = async (leaveId, reason) => {
  const response = await fetch(`/api/leaves/${leaveId}/withdraw`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ reason })
  });
  return await response.json();
};
```

---

## ✅ Status

✅ **Implementation Complete**  
✅ **No Syntax Errors**  
✅ **Email Notifications Added**  
✅ **Documentation Complete**  
✅ **Production Ready**  

---

**Created**: November 27, 2024  
**Status**: Ready for Testing 🚀
