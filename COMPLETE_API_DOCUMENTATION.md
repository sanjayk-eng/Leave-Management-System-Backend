# Complete API Documentation - Zenithive User Management System

## 📋 Table of Contents

1. [System Overview](#system-overview)
2. [Authentication](#authentication)
3. [Employee Management](#employee-management)
4. [Leave Management](#leave-management)
5. [Leave Balances](#leave-balances)
6. [Payroll Management](#payroll-management)
7. [Company Settings](#company-settings)
8. [Holidays](#holidays)
9. [Role-Based Permissions](#role-based-permissions)
10. [Business Logic](#business-logic)
11. [Error Handling](#error-handling)

---

## System Overview

### Base URL
```
http://localhost:{APP_PORT}/api
```

### Technology Stack
- **Backend**: Go (Gin Framework)
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **PDF Generation**: gofpdf

### Roles
- **SUPERADMIN**: Full system access, can finalize payroll
- **ADMIN**: Administrative access, cannot finalize payroll
- **HR**: Human resources operations
- **MANAGER**: Team management and leave approvals
- **EMPLOYEE**: Basic employee operations

---

## Authentication

### 1. Login
**POST** `/api/auth/login`

Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "email": "user@zenithive.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@zenithive.com",
    "role": "EMPLOYEE"
  }
}
```

**Error Responses:**
- `400`: Invalid request payload
- `401`: Login failed — email not found or wrong password

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@zenithive.com",
    "password": "admin123"
  }'
```

---

## Employee Management

### 1. List All Employees
**GET** `/api/employee/`

**Access**: SUPERADMIN, ADMIN

**Response:**
```json
{
  "message": "Employees fetched",
  "employees": [
    {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@zenithive.com",
      "role": "EMPLOYEE",
      "status": "active",
      "manager_id": "uuid",
      "manager_name": "Jane Smith",
      "salary": 50000.00,
      "joining_date": "2024-01-15T00:00:00Z"
    }
  ]
}
```

### 2. Get My Team
**GET** `/api/employee/my-team`

**Access**: MANAGER

Get all employees reporting to the current manager.

**Response:**
```json
{
  "message": "Team members fetched successfully",
  "team_members": [
    {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@zenithive.com",
      "role": "EMPLOYEE",
      "status": "active"
    }
  ]
}
```

### 3. Get Employee By ID
**GET** `/api/employee/:id`

**Access**: All authenticated users

**Response:**
```json
{
  "message": "employee details fetched successfully",
  "employee": {
    "id": "uuid",
    "full_name": "John Doe",
    "email": "john@zenithive.com",
    "role": "EMPLOYEE",
    "status": "active",
    "manager_name": "Jane Smith",
    "salary": 50000.00,
    "joining_date": "2024-01-15T00:00:00Z"
  }
}
```

### 4. Create Employee
**POST** `/api/employee/`

**Access**: SUPERADMIN, ADMIN

**Request Body:**
```json
{
  "full_name": "Jane Smith",
  "email": "jane@zenithive.com",
  "role": "EMPLOYEE",
  "password": "password123",
  "salary": 45000.00,
  "joining_date": "2024-11-25T00:00:00Z"
}
```

**Success Response (201):**
```json
{
  "message": "employee created"
}
```

**Validations:**
- Email must end with `@zenithive.com`
- Email must be unique
- Role must exist in system

### 5. Update Employee Information
**PATCH** `/api/employee/:id`

**Access**: 
- Anyone can update their own **name**
- Only SUPERADMIN and ADMIN can update **email, salary, joining_date**

**Request Body (all fields optional):**
```json
{
  "full_name": "John Updated Doe",
  "email": "john.updated@zenithive.com",
  "salary": 55000.00,
  "joining_date": "2024-01-01T00:00:00Z"
}
```

**Success Response (200):**
```json
{
  "message": "employee information updated successfully",
  "employee_id": "uuid"
}
```

### 6. Update Employee Password
**PATCH** `/api/employee/:id/password`

**Access**: SUPERADMIN, ADMIN, HR

**Request Body:**
```json
{
  "new_password": "newSecurePassword123"
}
```

**Success Response (200):**
```json
{
  "message": "password updated successfully",
  "employee_id": "uuid"
}
```

**Features:**
- Password hashed with bcrypt
- Minimum 6 characters
- Email notification sent to employee

### 7. Update Employee Role
**PATCH** `/api/employee/:id/role`

**Access**: SUPERADMIN, ADMIN, HR

**Request Body:**
```json
{
  "role": "MANAGER"
}
```

**Success Response (200):**
```json
{
  "message": "role updated successfully",
  "employee_id": "uuid",
  "old_role": "EMPLOYEE",
  "new_role": "MANAGER"
}
```

**Restrictions:**
- ADMIN and HR cannot change their own role
- ADMIN and HR cannot modify SUPERADMIN users
- ADMIN and HR cannot promote to SUPERADMIN
- Cannot change role if employee is a manager with subordinates

### 8. Update Employee Manager
**PATCH** `/api/employee/:id/manager`

**Access**: SUPERADMIN, ADMIN, HR

**Request Body:**
```json
{
  "manager_id": "uuid"
}
```

**Success Response (200):**
```json
{
  "message": "manager updated successfully",
  "employee_id": "uuid",
  "manager_id": "uuid"
}
```

**Restrictions:**
- Cannot assign employee as their own manager
- Non-SUPERADMIN cannot assign themselves as manager to others
- Manager must have MANAGER role
- Manager must be active

### 9. Deactivate/Activate Employee
**PUT** `/api/employee/deactivate/:id`

**Access**: SUPERADMIN, ADMIN, HR

**Success Response (200):**
```json
{
  "message": "Employee status updated successfully",
  "new_status": "deactive"
}
```

**Restrictions:**
- ADMIN and HR cannot deactivate SUPERADMIN users

---

## Leave Management

### 1. Apply for Leave
**POST** `/api/leaves/apply`

**Access**: EMPLOYEE

**Request Body:**
```json
{
  "leave_type_id": 1,
  "start_date": "2024-12-01T00:00:00Z",
  "end_date": "2024-12-05T00:00:00Z",
  "reason": "Family vacation planned for the holidays"
}
```

**Success Response (200):**
```json
{
  "message": "Leave applied successfully",
  "leave_id": "uuid",
  "days": 5
}
```

**Business Logic:**
- Calculates working days (excludes weekends and holidays)
- Checks leave balance
- Checks for overlapping leaves
- Status set to "Pending"
- Manager notified via email

### 2. Admin Add Leave
**POST** `/api/leaves/admin-add`

**Access**: SUPERADMIN, ADMIN, MANAGER (if setting enabled)

**Request Body:**
```json
{
  "employee_id": "uuid",
  "leave_type_id": 1,
  "start_date": "2024-12-10T00:00:00Z",
  "end_date": "2024-12-12T00:00:00Z",
  "reason": "Medical emergency"
}
```

**Success Response (200):**
```json
{
  "message": "Leave added successfully",
  "leave_id": "uuid",
  "days": 3
}
```

**Features:**
- Leave status set to "APPROVED" immediately
- Balance updated automatically
- Manager can only add for their team members (if setting enabled)

### 3. Create Leave Policy
**POST** `/api/leaves/admin-add/policy`

**Access**: SUPERADMIN

**Request Body:**
```json
{
  "name": "Sick Leave",
  "is_paid": true,
  "default_entitlement": 10
}
```

**Success Response (200):**
```json
{
  "id": 3,
  "name": "Sick Leave",
  "is_paid": true,
  "default_entitlement": 10,
  "created_at": "2024-11-25T10:00:00Z"
}
```

### 4. Get All Leave Policies
**GET** `/api/leaves/Get-All-Leave-Policy`

**Access**: All authenticated users

**Response:**
```json
[
  {
    "id": 1,
    "name": "Annual Leave",
    "is_paid": true,
    "default_entitlement": 20,
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

### 5. Approve/Reject Leave
**POST** `/api/leaves/:id/action`

**Access**: MANAGER, ADMIN, SUPERADMIN

**Request Body:**
```json
{
  "action": "APPROVE"
}
```
or
```json
{
  "action": "REJECT"
}
```

**Success Response (200):**
```json
{
  "message": "Leave approved successfully"
}
```

**Business Logic:**
- On APPROVE: Balance updated, employee notified
- On REJECT: Balance unchanged, employee notified
- Only pending leaves can be actioned

### 6. Cancel Leave
**DELETE** `/api/leaves/:id/cancel`

**Access**: Employee (own leaves), ADMIN, SUPERADMIN

**Success Response (200):**
```json
{
  "message": "leave request cancelled successfully",
  "leave_id": "uuid"
}
```

**Restrictions:**
- Only pending leaves can be cancelled
- Employee can only cancel their own leaves
- Email notification sent

### 7. Withdraw Leave
**POST** `/api/leaves/:id/withdraw`

**Access**: ADMIN, SUPERADMIN, MANAGER (for team members)

**Request Body:**
```json
{
  "reason": "Project deadline requires employee presence"
}
```

**Success Response (200):**
```json
{
  "message": "Leave withdrawn successfully",
  "leave_id": "uuid"
}
```

**Business Logic:**
- Only approved leaves can be withdrawn
- Leave balance restored
- Employee notified with reason

### 8. Get All Leaves
**GET** `/api/leaves/all`

**Access**: All authenticated users (filtered by role)

**Response:**
```json
{
  "total": 2,
  "data": [
    {
      "id": "uuid",
      "employee": "John Doe",
      "leave_type": "Annual Leave",
      "start_date": "2024-12-01T00:00:00Z",
      "end_date": "2024-12-05T00:00:00Z",
      "days": 5,
      "reason": "Family vacation",
      "status": "Pending",
      "applied_at": "2024-11-25T10:00:00Z"
    }
  ]
}
```

**Filtering:**
- EMPLOYEE: Only their own leaves
- MANAGER: Team members' leaves
- ADMIN/SUPERADMIN: All leaves

### 9. Get Manager Leave History
**GET** `/api/leaves/manager/history`

**Access**: MANAGER

Get leave history for all team members.

**Response:**
```json
{
  "message": "Manager leave history fetched successfully",
  "total": 10,
  "leaves": [...]
}
```

---

## Leave Balances

### 1. Get Leave Balances
**GET** `/api/leave-balances/employee/:id`

**Access**: All authenticated users

**Response:**
```json
{
  "employee_id": "uuid",
  "balances": [
    {
      "leave_type": "Annual Leave",
      "used": 5,
      "total": 20,
      "available": 15
    }
  ]
}
```

**Permission:**
- EMPLOYEE: Can only view own balances
- Others: Can view any employee's balances

### 2. Adjust Leave Balance
**POST** `/api/leave-balances/:id/adjust`

**Access**: ADMIN, SUPERADMIN

**Request Body:**
```json
{
  "leave_type_id": 1,
  "quantity": 5,
  "reason": "Bonus leave for exceptional performance"
}
```

**Success Response (200):**
```json
{
  "message": "Leave balance adjusted successfully",
  "new_adjusted": 5,
  "new_closing": 25,
  "year": 2024
}
```

**Features:**
- Use positive quantity to add leaves
- Use negative quantity to deduct leaves
- Adjustment logged in audit table

---

## Payroll Management

### 1. Run Payroll
**POST** `/api/payroll/run`

**Access**: SUPERADMIN, ADMIN

**Request Body:**
```json
{
  "month": 11,
  "year": 2024
}
```

**Success Response (200):**
```json
{
  "payroll_run_id": "uuid",
  "month": 11,
  "year": 2024,
  "total_payroll": 285000.00,
  "total_deductions": 15000.00,
  "employees_count": 6,
  "payroll_preview": [
    {
      "employee_id": "uuid",
      "employee": "John Doe",
      "basic_salary": 50000.00,
      "working_days": 22,
      "absent_days": 2,
      "deductions": 4545.45,
      "net_salary": 45454.55
    }
  ]
}
```

**Business Logic:**
- Calculates deductions based on approved leaves
- Formula: `Deduction = (Salary / Working Days) × Absent Days`
- Creates payroll run with status "PREVIEW"
- Cannot run for future months

### 2. Finalize Payroll
**POST** `/api/payroll/:id/finalize`

**Access**: **SUPERADMIN ONLY** ⚠️

**Success Response (200):**
```json
{
  "message": "Payroll finalized successfully",
  "payroll_run_id": "uuid",
  "payslips": ["uuid1", "uuid2"]
}
```

**Business Logic:**
- Generates payslips for all employees
- Updates status to "FINALIZED"
- Cannot finalize twice
- Only SUPERADMIN can finalize (security measure)

### 3. Get Finalized Payslips
**GET** `/api/payroll/payslip`

**Access**: 
- SUPERADMIN, ADMIN: All payslips
- EMPLOYEE: Only own payslips

**Response:**
```json
{
  "message": "Finalized payslips fetched successfully",
  "total_payslips": 10,
  "data": [
    {
      "payslip_id": "uuid",
      "employee_id": "uuid",
      "full_name": "John Doe",
      "month": 11,
      "year": 2024,
      "basic_salary": 50000.00,
      "net_salary": 45454.55,
      "pdf_path": "./tmp/payslip_uuid.pdf"
    }
  ]
}
```

### 4. Download Payslip PDF
**GET** `/api/payroll/payslips/:id/pdf`

**Access**: All authenticated users

**Response**: PDF file download

**PDF Features:**
- Professional design with company branding
- Employee information
- Earnings breakdown
- Deductions details
- Attendance summary
- Calculation breakdown
- Generated timestamp

---

## Company Settings

### 1. Get Company Settings
**GET** `/api/settings/`

**Access**: ADMIN, SUPERADMIN

**Response:**
```json
{
  "settings": {
    "id": "uuid",
    "working_days_per_month": 22,
    "allow_manager_add_leave": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-11-25T10:00:00Z"
  }
}
```

### 2. Update Company Settings
**PUT** `/api/settings/`

**Access**: ADMIN, SUPERADMIN

**Request Body:**
```json
{
  "working_days_per_month": 22,
  "allow_manager_add_leave": true
}
```

**Success Response (200):**
```json
{
  "message": "Company settings updated successfully"
}
```

---

## Holidays

### 1. Add Holiday
**POST** `/api/settings/holidays/`

**Access**: SUPERADMIN

**Request Body:**
```json
{
  "name": "Christmas",
  "date": "2024-12-25T00:00:00Z",
  "type": "HOLIDAY"
}
```

**Success Response (200):**
```json
{
  "message": "Holiday added successfully",
  "id": "uuid"
}
```

### 2. Get All Holidays
**GET** `/api/settings/holidays/`

**Access**: SUPERADMIN

**Response:**
```json
[
  {
    "id": 1,
    "name": "Christmas",
    "date": "2024-12-25T00:00:00Z",
    "day": "Wednesday",
    "type": "HOLIDAY",
    "created_at": "2024-11-25T10:00:00Z"
  }
]
```

### 3. Delete Holiday
**DELETE** `/api/settings/holidays/:id`

**Access**: SUPERADMIN

**Success Response (200):**
```json
{
  "message": "Holiday deleted successfully"
}
```

---

## Role-Based Permissions

### Complete Permission Matrix

| Feature | SUPERADMIN | ADMIN | HR | MANAGER | EMPLOYEE |
|---------|------------|-------|-----|---------|----------|
| **Authentication** |
| Login | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Employee Management** |
| List All Employees | ✅ | ✅ | ❌ | ❌ | ❌ |
| Get My Team | ❌ | ❌ | ❌ | ✅ | ❌ |
| Get Employee By ID | ✅ | ✅ | ✅ | ✅ | ✅ |
| Create Employee | ✅ | ✅ | ❌ | ❌ | ❌ |
| Update Name | ✅ All | ✅ All | ✅ Own | ✅ Own | ✅ Own |
| Update Email/Salary/Joining Date | ✅ | ✅ | ❌ | ❌ | ❌ |
| Update Password | ✅ | ✅ | ✅ | ❌ | ❌ |
| Update Role | ✅ | ✅ | ✅ | ❌ | ❌ |
| Update Manager | ✅ | ✅ | ✅ | ❌ | ❌ |
| Deactivate Employee | ✅ | ✅ | ✅ | ❌ | ❌ |
| **Leave Management** |
| Apply Leave | ✅ | ✅ | ✅ | ✅ | ✅ |
| Admin Add Leave | ✅ | ✅ | ❌ | ✅* | ❌ |
| Create Leave Policy | ✅ | ❌ | ❌ | ❌ | ❌ |
| Get Leave Policies | ✅ | ✅ | ✅ | ✅ | ✅ |
| Approve/Reject Leave | ✅ | ✅ | ✅ | ✅ | ❌ |
| Cancel Leave | ✅ All | ✅ All | ❌ | ❌ | ✅ Own |
| Withdraw Leave | ✅ | ✅ | ❌ | ✅ Team | ❌ |
| Get All Leaves | ✅ All | ✅ All | ✅ All | ✅ Team | ✅ Own |
| Get Manager History | ❌ | ❌ | ❌ | ✅ | ❌ |
| **Leave Balances** |
| Get Leave Balances | ✅ All | ✅ All | ✅ All | ✅ All | ✅ Own |
| Adjust Leave Balance | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Payroll** |
| Run Payroll | ✅ | ✅ | ❌ | ❌ | ❌ |
| Finalize Payroll | ✅ | ❌ | ❌ | ❌ | ❌ |
| Get Payslips | ✅ All | ✅ All | ❌ | ❌ | ✅ Own |
| Download Payslip PDF | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Settings** |
| Get Settings | ✅ | ✅ | ❌ | ❌ | ❌ |
| Update Settings | ✅ | ✅ | ❌ | ❌ | ❌ |
| Manage Holidays | ✅ | ❌ | ❌ | ❌ | ❌ |

*Manager can add leave only if `allow_manager_add_leave` setting is enabled

---

## Business Logic

### Leave Calculation
```
Working Days = Count(Monday-Friday) - Holidays
Excludes: Saturdays, Sundays, Company Holidays
```

### Payroll Calculation
```
Per Day Salary = Basic Salary / Working Days Per Month
Deduction = Per Day Salary × Absent Days
Net Salary = Basic Salary - Deduction
```

### Leave Balance Formula
```
Closing Balance = Opening + Accrued - Used + Adjusted
```

### Leave Approval Workflow
```
1. Employee applies → Status: Pending
2. Manager/Admin reviews
3a. Approved → Balance updated, Status: APPROVED
3b. Rejected → Balance unchanged, Status: REJECTED
```

### Payroll Workflow
```
1. ADMIN runs payroll → Status: PREVIEW
2. Review calculations
3. SUPERADMIN finalizes → Status: FINALIZED
4. Payslips generated
5. Employees download PDFs
```

---

## Error Handling

### Standard Error Response Format
```json
{
  "error": {
    "code": 400,
    "message": "Detailed error message"
  }
}
```

### Common HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful request |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid input or validation error |
| 401 | Unauthorized | Missing or invalid authentication token |
| 403 | Forbidden | User doesn't have permission |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server-side error |

### Common Error Messages

**Authentication:**
- "Missing Authorization header"
- "Invalid or expired token"
- "Login failed — email not found"
- "Login failed — wrong password"

**Permissions:**
- "not permitted"
- "Only SUPERADMIN can finalize payroll"
- "ADMIN and HR cannot change their own role"
- "you can only update your own name"

**Validation:**
- "invalid employee ID"
- "email must end with @zenithive.com"
- "email already exists"
- "password must be at least 6 characters long"

**Business Logic:**
- "Insufficient leave balance"
- "Overlapping leave exists"
- "cannot cancel leave with status: APPROVED"
- "Payroll already finalized"

---

## Security Features

### Authentication
- ✅ JWT token-based authentication
- ✅ Token expiration
- ✅ Password hashing with bcrypt
- ✅ Minimum password length (6 characters)

### Authorization
- ✅ Role-based access control (RBAC)
- ✅ Route-level middleware protection
- ✅ Function-level permission checks
- ✅ Self-modification restrictions

### Data Protection
- ✅ Password hash never exposed in responses
- ✅ Email domain validation (@zenithive.com)
- ✅ Email uniqueness enforcement
- ✅ SQL injection prevention (parameterized queries)

### Audit Trail
- ✅ Updated_at timestamps on all modifications
- ✅ Leave adjustment logging
- ✅ Payroll finalization tracking
- ✅ Email notifications for critical actions

---

## Email Notifications

### Notification Events

1. **Employee Created**
   - Recipient: New employee
   - Content: Welcome message with credentials

2. **Leave Applied**
   - Recipient: Manager, Admin, SUPERADMIN
   - Content: Leave request details

3. **Leave Approved**
   - Recipient: Employee
   - Content: Approval confirmation

4. **Leave Rejected**
   - Recipient: Employee
   - Content: Rejection notification

5. **Leave Cancelled**
   - Recipient: Employee
   - Content: Cancellation confirmation

6. **Leave Withdrawn**
   - Recipient: Employee
   - Content: Withdrawal notification with reason

7. **Leave Added by Admin**
   - Recipient: Employee
   - Content: Leave addition notification

8. **Password Updated**
   - Recipient: Employee
   - Content: Password change notification

---

## Environment Variables

Required in `.env` file:

```env
APP_PORT=8080
FRONTEND_SERVER=http://localhost:3000
SERACT_KEY=your_jwt_secret_key
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=user_management_db
GOOGLE_SCRIPT_URL=your_email_service_url
```

---

## Quick Start Guide

### 1. Setup
```bash
# Clone repository
git clone <repository-url>

# Install dependencies
go mod download

# Setup database
# Run migrations in pkg/migration/

# Configure environment
cp .env.example .env
# Edit .env with your settings
```

### 2. Run Server
```bash
go run main.go
```

### 3. Test Authentication
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@zenithive.com",
    "password": "admin123"
  }'
```

### 4. Use Token
```bash
# Save token from login response
TOKEN="your_jwt_token_here"

# Make authenticated request
curl -X GET http://localhost:8080/api/employee/ \
  -H "Authorization: Bearer $TOKEN"
```

---

## Best Practices

### For Developers
- ✅ Always use parameterized queries
- ✅ Validate input data
- ✅ Handle errors gracefully
- ✅ Use transactions for multi-step operations
- ✅ Log important events
- ✅ Send notifications asynchronously

### For API Consumers
- ✅ Store JWT token securely
- ✅ Handle token expiration
- ✅ Validate responses
- ✅ Handle errors appropriately
- ✅ Use HTTPS in production
- ✅ Don't expose sensitive data in logs

### For Administrators
- ✅ Review payroll before finalizing
- ✅ Regularly backup database
- ✅ Monitor system logs
- ✅ Keep employee data updated
- ✅ Configure holidays annually
- ✅ Review leave policies periodically

---

## Support & Maintenance

### Common Issues

**Issue**: Token expired
**Solution**: Re-login to get new token

**Issue**: Email not sent
**Solution**: Check GOOGLE_SCRIPT_URL configuration

**Issue**: Cannot finalize payroll
**Solution**: Ensure user is SUPERADMIN

**Issue**: Leave balance incorrect
**Solution**: Use adjust balance endpoint to correct

---

## Version History

**Version 1.0** (November 2024)
- Initial release
- Complete employee management
- Leave management system
- Payroll processing
- Role-based access control
- Email notifications

---

## Contact & Support

For issues, questions, or feature requests, please contact the development team.

---

**Last Updated**: November 27, 2024  
**API Version**: 1.0  
**Status**: Production Ready ✅
