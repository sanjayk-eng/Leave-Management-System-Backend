# Leave Cancel & Withdraw Features

## Overview

Two new endpoints have been added to manage leave requests:

1. **Cancel Leave** - For pending leave requests
2. **Withdraw Leave** - For approved leave requests

---

## 1. Cancel Leave (Pending Requests)

### Endpoint
```
DELETE /api/leaves/:id/cancel
```

### Description
Allows employees to cancel their own pending leave requests, or admins to cancel any pending leave.

### Authentication
Required: JWT Bearer Token

### URL Parameters
- `id` (UUID, required): Leave request ID

### Permissions

| Role | Can Cancel |
|------|------------|
| SUPERADMIN | ✅ Any pending leave |
| ADMIN | ✅ Any pending leave |
| HR | ❌ No |
| MANAGER | ❌ No |
| EMPLOYEE | ✅ Own pending leave only |

### Business Rules
- ✅ Only leaves with status **"Pending"** can be cancelled
- ✅ Employee can only cancel their own leave
- ✅ Admin can cancel any pending leave
- ❌ Cannot cancel approved, rejected, or withdrawn leaves

### Success Response (200)
```json
{
  "message": "leave request cancelled successfully",
  "leave_id": "770e8400-e29b-41d4-a716-446655440002"
}
```

### Error Responses

**400 Bad Request - Invalid Leave ID**
```json
{
  "error": {
    "code": 400,
    "message": "invalid leave ID"
  }
}
```

**400 Bad Request - Cannot Cancel**
```json
{
  "error": {
    "code": 400,
    "message": "cannot cancel leave with status: APPROVED. Only pending leaves can be cancelled"
  }
}
```

**403 Forbidden**
```json
{
  "error": {
    "code": 403,
    "message": "you can only cancel your own leave requests"
  }
}
```

**404 Not Found**
```json
{
  "error": {
    "code": 404,
    "message": "leave request not found"
  }
}
```

### cURL Examples

**Employee cancels own pending leave:**
```bash
curl -X DELETE http://localhost:8080/api/leaves/770e8400-e29b-41d4-a716-446655440002/cancel \
  -H "Authorization: Bearer <employee_token>"
```

**Admin cancels any pending leave:**
```bash
curl -X DELETE http://localhost:8080/api/leaves/770e8400-e29b-41d4-a716-446655440002/cancel \
  -H "Authorization: Bearer <admin_token>"
```

### Email Notification

When a leave is cancelled, the employee receives:

**Subject:** Leave Request Cancelled

**Body:**
```
Dear [Employee Name],

Your leave request has been cancelled.

Leave Type: Annual Leave
Start Date: 2024-12-01
End Date: 2024-12-05
Duration: 5.0 days
Status: CANCELLED

If you did not cancel this leave request, please contact your manager or HR department immediately.

Best regards,
Zenithive Leave Management System
```

---

## 2. Withdraw Leave (Approved Requests)

### Endpoint
```
POST /api/leaves/:id/withdraw
```

### Description
Allows Admin/Manager to withdraw an approved leave and restore the employee's leave balance.

### Authentication
Required: JWT Bearer Token

### URL Parameters
- `id` (UUID, required): Leave request ID

### Request Body (Optional)
```json
{
  "reason": "Emergency project requirement"
}
```

### Permissions

| Role | Can Withdraw |
|------|--------------|
| SUPERADMIN | ✅ Any approved leave |
| ADMIN | ✅ Any approved leave |
| HR | ❌ No |
| MANAGER | ✅ Team members' approved leave only |
| EMPLOYEE | ❌ No |

### Business Rules
- ✅ Only leaves with status **"APPROVED"** can be withdrawn
- ✅ Manager can only withdraw leaves of their team members
- ✅ Admin can withdraw any approved leave
- ✅ Leave balance is automatically restored
- ✅ Optional reason can be provided
- ❌ Cannot withdraw pending, rejected, or cancelled leaves

### Success Response (200)
```json
{
  "message": "leave withdrawn successfully and balance restored",
  "leave_id": "770e8400-e29b-41d4-a716-446655440002",
  "days_restored": 5,
  "withdrawal_by": "660e8400-e29b-41d4-a716-446655440001",
  "withdrawal_role": "ADMIN",
  "withdrawal_reason": "Emergency project requirement"
}
```

**Note**: If no reason is provided, it defaults to "Withdrawn by [ROLE]"

### Error Responses

**400 Bad Request - Invalid Leave ID**
```json
{
  "error": {
    "code": 400,
    "message": "invalid leave ID"
  }
}
```

**400 Bad Request - Cannot Withdraw**
```json
{
  "error": {
    "code": 400,
    "message": "cannot withdraw leave with status: Pending. Only approved leaves can be withdrawn"
  }
}
```

**403 Forbidden - Not Authorized**
```json
{
  "error": {
    "code": 403,
    "message": "only SUPERADMIN, ADMIN, and MANAGER can withdraw approved leaves"
  }
}
```

**403 Forbidden - Manager Restriction**
```json
{
  "error": {
    "code": 403,
    "message": "managers can only withdraw leaves of their team members"
  }
}
```

**404 Not Found**
```json
{
  "error": {
    "code": 404,
    "message": "leave request not found"
  }
}
```

### cURL Examples

**Admin withdraws approved leave:**
```bash
curl -X POST http://localhost:8080/api/leaves/770e8400-e29b-41d4-a716-446655440002/withdraw \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Emergency project requirement"
  }'
```

**Manager withdraws team member's leave:**
```bash
curl -X POST http://localhost:8080/api/leaves/770e8400-e29b-41d4-a716-446655440002/withdraw \
  -H "Authorization: Bearer <manager_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Urgent client meeting"
  }'
```

**Withdraw without reason:**
```bash
curl -X POST http://localhost:8080/api/leaves/770e8400-e29b-41d4-a716-446655440002/withdraw \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Email Notification

When a leave is withdrawn, the employee receives:

**Subject:** Leave Request Withdrawn

**Body:**
```
Dear [Employee Name],

Your approved leave request has been withdrawn by [Admin Name] (ADMIN).

Leave Type: Annual Leave
Start Date: 2024-12-01
End Date: 2024-12-05
Duration: 5.0 days
Status: WITHDRAWN
Reason: Emergency project requirement

Your leave balance has been restored. The 5.0 days have been credited back to your account.

If you have any questions about this withdrawal, please contact your manager or HR department.

Best regards,
Zenithive Leave Management System
```

---

## Comparison: Cancel vs Withdraw

| Aspect | Cancel | Withdraw |
|--------|--------|----------|
| **Applies To** | Pending leaves | Approved leaves |
| **Who Can Do** | Employee (own) or Admin | Admin or Manager |
| **Balance Impact** | No impact (not yet deducted) | Restores balance |
| **Status Change** | Pending → CANCELLED | APPROVED → WITHDRAWN |
| **Reason Required** | No | Optional |
| **Notification** | Yes | Yes |

---

## Leave Status Flow

```
┌─────────┐
│ Pending │ ──cancel──> CANCELLED
└────┬────┘
     │
  approve
     │
     ▼
┌──────────┐
│ APPROVED │ ──withdraw──> WITHDRAWN
└──────────┘
```

---

## Use Cases

### Use Case 1: Employee Changes Mind
**Scenario**: Employee applied for leave but changed their mind

**Action**: Employee cancels pending leave
```bash
DELETE /api/leaves/:id/cancel
```

**Result**: 
- Leave status → CANCELLED
- No balance impact (wasn't deducted yet)
- Email notification sent

---

### Use Case 2: Emergency at Work
**Scenario**: Employee's approved leave needs to be cancelled due to emergency

**Action**: Manager withdraws approved leave
```bash
POST /api/leaves/:id/withdraw
Body: { "reason": "Critical production issue" }
```

**Result**:
- Leave status → WITHDRAWN
- Leave balance restored
- Email notification with reason sent

---

### Use Case 3: Admin Corrects Mistake
**Scenario**: Leave was approved by mistake

**Action**: Admin withdraws the leave
```bash
POST /api/leaves/:id/withdraw
Body: { "reason": "Approved in error" }
```

**Result**:
- Leave status → WITHDRAWN
- Balance restored
- Employee notified

---

## Database Changes

### Leave Status Values
- `Pending` - Initial status when applied
- `APPROVED` - Approved by manager/admin
- `REJECTED` - Rejected by manager/admin
- `CANCELLED` - Cancelled by employee/admin (pending only)
- `WITHDRAWN` - Withdrawn by admin/manager (approved only)

### Leave Balance Impact

**Cancel (Pending Leave)**:
- No database update needed
- Balance was never deducted

**Withdraw (Approved Leave)**:
```sql
UPDATE Tbl_Leave_balance 
SET used = used - [days], 
    closing = closing + [days], 
    updated_at = NOW()
WHERE employee_id = [id] 
  AND leave_type_id = [type_id] 
  AND year = CURRENT_YEAR
```

---

## Testing Checklist

### Cancel Leave Tests
- [ ] Employee can cancel own pending leave ✅
- [ ] Employee cannot cancel others' pending leave ✅
- [ ] Admin can cancel any pending leave ✅
- [ ] Cannot cancel approved leave ✅
- [ ] Cannot cancel rejected leave ✅
- [ ] Cannot cancel already cancelled leave ✅
- [ ] Email notification sent ✅

### Withdraw Leave Tests
- [ ] Admin can withdraw any approved leave ✅
- [ ] Manager can withdraw team member's leave ✅
- [ ] Manager cannot withdraw other team's leave ✅
- [ ] Employee cannot withdraw leave ✅
- [ ] Cannot withdraw pending leave ✅
- [ ] Cannot withdraw rejected leave ✅
- [ ] Balance is restored correctly ✅
- [ ] Email notification sent with reason ✅

---

## Frontend Integration

### React - Cancel Leave
```javascript
const cancelLeave = async (leaveId) => {
  try {
    const response = await fetch(`/api/leaves/${leaveId}/cancel`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error.message);
    }

    const data = await response.json();
    alert(data.message);
    return data;
  } catch (error) {
    console.error('Error cancelling leave:', error);
    throw error;
  }
};
```

### React - Withdraw Leave
```javascript
const withdrawLeave = async (leaveId, reason = '') => {
  try {
    const response = await fetch(`/api/leaves/${leaveId}/withdraw`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ reason })
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error.message);
    }

    const data = await response.json();
    alert(`${data.message}. ${data.days_restored} days restored.`);
    return data;
  } catch (error) {
    console.error('Error withdrawing leave:', error);
    throw error;
  }
};
```

### Component Example
```javascript
const LeaveActions = ({ leave, currentUser }) => {
  const canCancel = leave.status === 'Pending' && 
    (leave.employee_id === currentUser.id || currentUser.role === 'ADMIN');
  
  const canWithdraw = leave.status === 'APPROVED' && 
    (currentUser.role === 'ADMIN' || currentUser.role === 'SUPERADMIN' || 
     (currentUser.role === 'MANAGER' && leave.manager_id === currentUser.id));

  return (
    <div>
      {canCancel && (
        <button onClick={() => cancelLeave(leave.id)}>
          Cancel Leave
        </button>
      )}
      
      {canWithdraw && (
        <button onClick={() => {
          const reason = prompt('Enter reason for withdrawal (optional):');
          withdrawLeave(leave.id, reason);
        }}>
          Withdraw Leave
        </button>
      )}
    </div>
  );
};
```

---

## Files Modified

1. ✅ `controllers/leave.go` - Added CancelLeave and WithdrawLeave functions
2. ✅ `routes/router.go` - Added cancel and withdraw routes
3. ✅ `utils/notification.go` - Added email notification functions
4. ✅ `LEAVE_CANCEL_WITHDRAW.md` - This documentation

---

## Summary

✅ **Cancel Leave** - For pending requests  
✅ **Withdraw Leave** - For approved requests  
✅ **Balance Restoration** - Automatic for withdrawals  
✅ **Email Notifications** - Sent for both actions  
✅ **Role-Based Permissions** - Proper access control  
✅ **Comprehensive Error Handling** - Clear error messages  
✅ **Documentation Complete** - Full API docs  

---

**Created**: November 27, 2024  
**Status**: ✅ Production Ready
