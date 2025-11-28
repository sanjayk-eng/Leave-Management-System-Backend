# Quick Reference Guide - Zenithive API

## 🚀 Quick Start

### Base URL
```
http://localhost:8080/api
```

### Authentication
```bash
# Login
POST /api/auth/login
{
  "email": "user@zenithive.com",
  "password": "password"
}

# Use token in all requests
Authorization: Bearer <your_token>
```

---

## 📋 Endpoint Summary

### Authentication (1 endpoint)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/auth/login` | Public | Login and get JWT token |

### Employee Management (10 endpoints)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/employee/` | SUPERADMIN, ADMIN | List all employees |
| GET | `/employee/my-team` | MANAGER | Get team members |
| GET | `/employee/:id` | All | Get employee details |
| POST | `/employee/` | SUPERADMIN, ADMIN | Create employee |
| PATCH | `/employee/:id` | All (limited) | Update employee info |
| PATCH | `/employee/:id/password` | SUPERADMIN, ADMIN, HR | Update password |
| PATCH | `/employee/:id/role` | SUPERADMIN, ADMIN, HR | Update role |
| PATCH | `/employee/:id/manager` | SUPERADMIN, ADMIN, HR | Update manager |
| PUT | `/employee/deactivate/:id` | SUPERADMIN, ADMIN, HR | Deactivate/Activate |
| GET | `/employee/:id/reports` | Self, Manager, Admin | Get direct reports |

### Leave Management (9 endpoints)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/leaves/apply` | EMPLOYEE | Apply for leave |
| POST | `/leaves/admin-add` | SUPERADMIN, ADMIN, MANAGER* | Add leave for employee |
| POST | `/leaves/admin-add/policy` | SUPERADMIN | Create leave policy |
| GET | `/leaves/Get-All-Leave-Policy` | All | Get leave policies |
| GET | `/leaves/manager/history` | MANAGER | Get team leave history |
| POST | `/leaves/:id/action` | MANAGER, ADMIN, SUPERADMIN | Approve/Reject leave |
| DELETE | `/leaves/:id/cancel` | Employee (own), ADMIN, SUPERADMIN | Cancel pending leave |
| POST | `/leaves/:id/withdraw` | ADMIN, SUPERADMIN, MANAGER | Withdraw approved leave |
| GET | `/leaves/all` | All (filtered) | Get all leaves |

### Leave Balances (2 endpoints)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/leave-balances/employee/:id` | All | Get leave balances |
| POST | `/leave-balances/:id/adjust` | ADMIN, SUPERADMIN | Adjust leave balance |

### Payroll (4 endpoints)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/payroll/run` | SUPERADMIN, ADMIN | Run payroll preview |
| POST | `/payroll/:id/finalize` | **SUPERADMIN ONLY** | Finalize payroll |
| GET | `/payroll/payslip` | All (filtered) | Get finalized payslips |
| GET | `/payroll/payslips/:id/pdf` | All | Download payslip PDF |

### Settings (3 endpoints)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/settings/` | ADMIN, SUPERADMIN | Get company settings |
| PUT | `/settings/` | ADMIN, SUPERADMIN | Update settings |

### Holidays (3 endpoints)
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/settings/holidays/` | SUPERADMIN | Add holiday |
| GET | `/settings/holidays/` | SUPERADMIN | List holidays |
| DELETE | `/settings/holidays/:id` | SUPERADMIN | Delete holiday |

**Total: 32 endpoints**

---

## 🔐 Role Permissions Quick View

| Feature | SUPERADMIN | ADMIN | HR | MANAGER | EMPLOYEE |
|---------|:----------:|:-----:|:--:|:-------:|:--------:|
| Manage Employees | ✅ | ✅ | ⚠️ | ❌ | ❌ |
| Manage Leaves | ✅ | ✅ | ✅ | ⚠️ | ⚠️ |
| Manage Payroll | ✅ | ⚠️ | ❌ | ❌ | ❌ |
| Manage Settings | ✅ | ✅ | ❌ | ❌ | ❌ |
| Finalize Payroll | ✅ | ❌ | ❌ | ❌ | ❌ |

✅ Full Access | ⚠️ Limited Access | ❌ No Access

---

## 💡 Common Use Cases

### 1. Employee Onboarding
```bash
# 1. Create employee
POST /api/employee/
{
  "full_name": "John Doe",
  "email": "john@zenithive.com",
  "role": "EMPLOYEE",
  "password": "temp123",
  "salary": 50000,
  "joining_date": "2024-12-01T00:00:00Z"
}

# 2. Assign manager
PATCH /api/employee/:id/manager
{
  "manager_id": "manager_uuid"
}
```

### 2. Leave Application Flow
```bash
# 1. Employee applies
POST /api/leaves/apply
{
  "leave_type_id": 1,
  "start_date": "2024-12-10T00:00:00Z",
  "end_date": "2024-12-12T00:00:00Z",
  "reason": "Family vacation"
}

# 2. Manager approves
POST /api/leaves/:id/action
{
  "action": "APPROVE"
}

# 3. Check balance
GET /api/leave-balances/employee/:id
```

### 3. Monthly Payroll
```bash
# 1. Run payroll (ADMIN)
POST /api/payroll/run
{
  "month": 11,
  "year": 2024
}

# 2. Review preview
# Check calculations

# 3. Finalize (SUPERADMIN only)
POST /api/payroll/:id/finalize

# 4. Employee downloads
GET /api/payroll/payslips/:id/pdf
```

---

## 🔑 Key Business Rules

### Leave Management
- ✅ Working days = Mon-Fri minus holidays
- ✅ Cannot apply for overlapping leaves
- ✅ Must have sufficient balance
- ✅ Only pending leaves can be cancelled
- ✅ Only approved leaves can be withdrawn

### Payroll
- ✅ Deduction = (Salary / Working Days) × Absent Days
- ✅ Only SUPERADMIN can finalize
- ✅ Cannot finalize twice
- ✅ Cannot run for future months

### Employee Management
- ✅ Email must end with @zenithive.com
- ✅ ADMIN/HR cannot change own role
- ✅ ADMIN/HR cannot modify SUPERADMIN
- ✅ Cannot assign self as manager

---

## 📧 Email Notifications

Automatic emails sent for:
- ✅ Employee created (welcome email)
- ✅ Leave applied (to manager)
- ✅ Leave approved/rejected (to employee)
- ✅ Leave cancelled (to employee)
- ✅ Leave withdrawn (to employee with reason)
- ✅ Leave added by admin (to employee)
- ✅ Password updated (to employee)

---

## ⚠️ Important Restrictions

### SUPERADMIN Only
- Finalize payroll
- Create leave policies
- Manage holidays
- Promote to SUPERADMIN

### ADMIN/HR Cannot
- Change their own role
- Modify SUPERADMIN users
- Finalize payroll
- Promote to SUPERADMIN

### MANAGER Can (if enabled)
- Add leave for team members
- Approve team leaves
- View team leave history

### EMPLOYEE Can
- Update own name
- Apply for leave
- Cancel own pending leaves
- View own data

---

## 🚨 Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| 401 Unauthorized | Invalid/expired token | Re-login |
| 403 Forbidden | Insufficient permissions | Check role |
| 400 Bad Request | Invalid input | Check request body |
| 404 Not Found | Resource doesn't exist | Verify ID |
| 500 Server Error | Server issue | Check logs |

---

## 📊 Data Formats

### Dates
```
ISO 8601: "2024-12-01T00:00:00Z"
```

### UUIDs
```
"550e8400-e29b-41d4-a716-446655440000"
```

### Currency
```
Float: 50000.00 (no currency symbol)
```

### Status Values
```
Employee: "active" | "deactive"
Leave: "Pending" | "APPROVED" | "REJECTED" | "CANCELLED" | "WITHDRAWN"
Payroll: "PREVIEW" | "FINALIZED"
```

---

## 🔧 Testing Commands

### Get Token
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@zenithive.com","password":"admin123"}' \
  | jq -r '.token')
```

### Test Endpoint
```bash
curl -X GET http://localhost:8080/api/employee/ \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📱 Response Formats

### Success Response
```json
{
  "message": "Operation successful",
  "data": {...}
}
```

### Error Response
```json
{
  "error": {
    "code": 400,
    "message": "Detailed error message"
  }
}
```

---

## 🎯 Quick Tips

1. **Always include Authorization header** (except login)
2. **Email must end with @zenithive.com**
3. **Passwords minimum 6 characters**
4. **UUIDs must be valid format**
5. **Dates in ISO 8601 format**
6. **Only SUPERADMIN can finalize payroll**
7. **Check role permissions before calling endpoints**
8. **Use transactions for multi-step operations**

---

**Last Updated**: November 27, 2024  
**Version**: 1.0
