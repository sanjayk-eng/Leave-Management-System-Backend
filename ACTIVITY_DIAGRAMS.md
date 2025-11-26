# Complete Activity Diagrams for All Routes

## Table of Contents
1. [Authentication Routes](#authentication-routes)
2. [Employee Routes](#employee-routes)
3. [Leave Routes](#leave-routes)
4. [Leave Balance Routes](#leave-balance-routes)
5. [Payroll Routes](#payroll-routes)
6. [Settings Routes](#settings-routes)
7. [Holiday Routes](#holiday-routes)

---

## Authentication Routes

### 1. POST /api/auth/login - User Login

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Client Sends Login      │
                    │  Request                 │
                    │  {email, password}       │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Parse JSON Request Body │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                          ┌──────────────┐
                          │  Valid JSON? │
                          └──────┬───────┘
                                 │
                     ┌───────────┴───────────┐
                     │ NO                    │ YES
                     ▼                       ▼
            ┌─────────────────┐    ┌─────────────────────┐
            │  Return 400     │    │  Query Database:    │
            │  "Invalid       │    │  SELECT employee    │
            │  request        │    │  WHERE email = ?    │
            │  payload"       │    └──────────┬──────────┘
            └─────────────────┘               │
                                              ▼
                                       ┌──────────────┐
                                       │  User Found? │
                                       └──────┬───────┘
                                              │
                                  ┌───────────┴───────────┐
                                  │ NO                    │ YES
                                  ▼                       ▼
                        ┌─────────────────┐    ┌─────────────────────┐
                        │  Return 401     │    │  Verify Password    │
                        │  "Login failed  │    │  using bcrypt       │
                        │  - email not    │    └──────────┬──────────┘
                        │  found"         │               │
                        └─────────────────┘               ▼
                                               ┌─────────────────────┐
                                               │  Password Correct?  │
                                               └──────────┬──────────┘
                                                          │
                                              ┌───────────┴───────────┐
                                              │ NO                    │ YES
                                              ▼                       ▼
                                    ┌─────────────────┐    ┌─────────────────────┐
                                    │  Return 401     │    │  Generate JWT Token │
                                    │  "Login failed  │    │  with:              │
                                    │  - wrong        │    │  - user_id          │
                                    │  password"      │    │  - role             │
                                    └─────────────────┘    │  - secret_key       │
                                                           └──────────┬──────────┘
                                                                      │
                                                                      ▼
                                                           ┌─────────────────────┐
                                                           │  Token Generated?   │
                                                           └──────────┬──────────┘
                                                                      │
                                                          ┌───────────┴───────────┐
                                                          │ NO                    │ YES
                                                          ▼                       ▼
                                                ┌─────────────────┐    ┌─────────────────────┐
                                                │  Return 500     │    │  Return 200         │
                                                │  "Failed to     │    │  {                  │
                                                │  generate       │    │    success: true,   │
                                                │  token"         │    │    message: "Login  │
                                                └─────────────────┘    │    successful",     │
                                                                       │    token: "jwt...", │
                                                                       │    user: {          │
                                                                       │      id, email,     │
                                                                       │      role           │
                                                                       │    }                │
                                                                       │  }                  │
                                                                       └──────────┬──────────┘
                                                                                  │
                                                                                  ▼
                                                                                 END
```

---

## Employee Routes

### 2. GET /api/employee/ - Get All Employees

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Client Sends Request    │
                    │  with JWT Token          │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Auth Middleware         │
                    │  Validates JWT Token     │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                          ┌──────────────┐
                          │  Valid Token?│
                          └──────┬───────┘
                                 │
                     ┌───────────┴───────────┐
                     │ NO                    │ YES
                     ▼                       ▼
            ┌─────────────────┐    ┌─────────────────────┐
            │  Return 401     │    │  Extract user_id    │
            │  "Unauthorized" │    │  and role from      │
            └─────────────────┘    │  token claims       │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Check Role:        │
                                   │  SUPERADMIN or      │
                                   │  Admin?             │
                                   └──────────┬──────────┘
                                              │
                                  ┌───────────┴───────────┐
                                  │ NO                    │ YES
                                  ▼                       ▼
                        ┌─────────────────┐    ┌─────────────────────┐
                        │  Return 401     │    │  Query Database:    │
                        │  "not permitted"│    │  SELECT * FROM      │
                        └─────────────────┘    │  Tbl_Employee       │
                                               │  JOIN Tbl_Role      │
                                               └──────────┬──────────┘
                                                          │
                                                          ▼
                                               ┌─────────────────────┐
                                               │  Query Successful?  │
                                               └──────────┬──────────┘
                                                          │
                                              ┌───────────┴───────────┐
                                              │ NO                    │ YES
                                              ▼                       ▼
                                    ┌─────────────────┐    ┌─────────────────────┐
                                    │  Return 500     │    │  Iterate Rows       │
                                    │  with error     │    │  Build Employee     │
                                    └─────────────────┘    │  Array              │
                                                           └──────────┬──────────┘
                                                                      │
                                                                      ▼
                                                           ┌─────────────────────┐
                                                           │  Return 200         │
                                                           │  {                  │
                                                           │    message:         │
                                                           │    "Employees       │
                                                           │    fetched",        │
                                                           │    employees: [...]  │
                                                           │  }                  │
                                                           └──────────┬──────────┘
                                                                      │
                                                                      ▼
                                                                     END
```

---

### 3. POST /api/employee/ - Create Employee

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Auth Middleware         │
                    │  Validates JWT           │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Check Role:             │
                    │  SUPERADMIN or Admin?    │
                    └────────────┬─────────────┘
                                 │
                     ┌───────────┴───────────┐
                     │ NO                    │ YES
                     ▼                       ▼
            ┌─────────────────┐    ┌─────────────────────┐
            │  Return 401     │    │  Parse JSON Body    │
            │  "not permitted"│    │  {full_name, email, │
            └─────────────────┘    │  role, password,    │
                                   │  salary, joining}   │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Validate Input     │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                       ┌──────────────┐
                                       │  Valid?      │
                                       └──────┬───────┘
                                              │
                                  ┌───────────┴───────────┐
                                  │ NO                    │ YES
                                  ▼                       ▼
                        ┌─────────────────┐    ┌─────────────────────┐
                        │  Return 400     │    │  Check Email Suffix │
                        │  with error     │    │  ends with          │
                        └─────────────────┘    │  @zenithive.com?    │
                                               └──────────┬──────────┘
                                                          │
                                              ┌───────────┴───────────┐
                                              │ NO                    │ YES
                                              ▼                       ▼
                                    ┌─────────────────┐    ┌─────────────────────┐
                                    │  Return 400     │    │  Query: Check Email │
                                    │  "email must    │    │  EXISTS in DB       │
                                    │  end with       │    └──────────┬──────────┘
                                    │  @zenithive.com"│               │
                                    └─────────────────┘               ▼
                                                           ┌─────────────────────┐
                                                           │  Email Exists?      │
                                                           └──────────┬──────────┘
                                                                      │
                                                          ┌───────────┴───────────┐
                                                          │ YES                   │ NO
                                                          ▼                       ▼
                                                ┌─────────────────┐    ┌─────────────────────┐
                                                │  Return 400     │    │  Query: Get Role ID │
                                                │  "email already │    │  WHERE type = ?     │
                                                │  exists"        │    └──────────┬──────────┘
                                                └─────────────────┘               │
                                                                                  ▼
                                                                       ┌─────────────────────┐
                                                                       │  Role Found?        │
                                                                       └──────────┬──────────┘
                                                                                  │
                                                                      ┌───────────┴───────────┐
                                                                      │ NO                    │ YES
                                                                      ▼                       ▼
                                                            ┌─────────────────┐    ┌─────────────────────┐
                                                            │  Return 400     │    │  Hash Password      │
                                                            │  "role not      │    │  using bcrypt       │
                                                            │  found"         │    └──────────┬──────────┘
                                                            └─────────────────┘               │
                                                                                              ▼
                                                                                   ┌─────────────────────┐
                                                                                   │  INSERT INTO        │
                                                                                   │  Tbl_Employee       │
                                                                                   │  (full_name, email, │
                                                                                   │  role_id, password, │
                                                                                   │  salary, joining)   │
                                                                                   └──────────┬──────────┘
                                                                                              │
                                                                                              ▼
                                                                                   ┌─────────────────────┐
                                                                                   │  Insert Successful? │
                                                                                   └──────────┬──────────┘
                                                                                              │
                                                                                  ┌───────────┴───────────┐
                                                                                  │ NO                    │ YES
                                                                                  ▼                       ▼
                                                                        ┌─────────────────┐    ┌─────────────────────┐
                                                                        │  Return 500     │    │  Spawn Goroutine    │
                                                                        │  "failed to     │    │  (Async)            │
                                                                        │  create         │    └──────────┬──────────┘
                                                                        │  employee"      │               │
                                                                        └─────────────────┘               ▼
                                                                                             ┌─────────────────────┐
                                                                                             │  Send Welcome Email │
                                                                                             │  to Employee with:  │
                                                                                             │  - Email            │
                                                                                             │  - Password         │
                                                                                             │  - Login URL        │
                                                                                             └──────────┬──────────┘
                                                                                                        │
                                                                                                        ▼
                                                                                             ┌─────────────────────┐
                                                                                             │  Return 201         │
                                                                                             │  {message:          │
                                                                                             │  "employee created"}│
                                                                                             └──────────┬──────────┘
                                                                                                        │
                                                                                                        ▼
                                                                                                       END
```

---

### 4. PATCH /api/employee/:id/role - Update Employee Role

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Auth Middleware         │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Check Role:             │
                    │  SUPERADMIN or ADMIN?    │
                    └────────────┬─────────────┘
                                 │
                     ┌───────────┴───────────┐
                     │ NO                    │ YES
                     ▼                       ▼
            ┌─────────────────┐    ┌─────────────────────┐
            │  Return 401     │    │  Extract Employee   │
            │  "not permitted"│    │  ID from URL Param  │
            └─────────────────┘    └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Parse JSON Body    │
                                   │  {role: "MANAGER"}  │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Query: Get Current │
                                   │  Role of Employee   │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                       ┌──────────────┐
                                       │  Found?      │
                                       └──────┬───────┘
                                              │
                                  ┌───────────┴───────────┐
                                  │ NO                    │ YES
                                  ▼                       ▼
                        ┌─────────────────┐    ┌─────────────────────┐
                        │  Return 500     │    │  Current Role ==    │
                        │  with error     │    │  New Role?          │
                        └─────────────────┘    └──────────┬──────────┘
                                                          │
                                              ┌───────────┴───────────┐
                                              │ YES                   │ NO
                                              ▼                       ▼
                                    ┌─────────────────┐    ┌─────────────────────┐
                                    │  Return 200     │    │  Query: Get Role ID │
                                    │  "already same  │    │  for New Role       │
                                    │  role"          │    └──────────┬──────────┘
                                    └─────────────────┘               │
                                                                      ▼
                                                           ┌─────────────────────┐
                                                           │  UPDATE             │
                                                           │  Tbl_Employee       │
                                                           │  SET role_id = ?    │
                                                           │  WHERE id = ?       │
                                                           └──────────┬──────────┘
                                                                      │
                                                                      ▼
                                                           ┌─────────────────────┐
                                                           │  Update Successful? │
                                                           └──────────┬──────────┘
                                                                      │
                                                          ┌───────────┴───────────┐
                                                          │ NO                    │ YES
                                                          ▼                       ▼
                                                ┌─────────────────┐    ┌─────────────────────┐
                                                │  Return 500     │    │  Return 200         │
                                                │  with error     │    │  {                  │
                                                └─────────────────┘    │    message: "role   │
                                                                       │    updated",        │
                                                                       │    employee_id: id  │
                                                                       │  }                  │
                                                                       └──────────┬──────────┘
                                                                                  │
                                                                                  ▼
                                                                                 END
```

---

### 5. PATCH /api/employee/:id/manager - Update Employee Manager

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Auth Middleware         │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Check Role:             │
                    │  SUPERADMIN/ADMIN/HR?    │
                    └────────────┬─────────────┘
                                 │
                     ┌───────────┴───────────┐
                     │ NO                    │ YES
                     ▼                       ▼
            ┌─────────────────┐    ┌─────────────────────┐
            │  Return 401     │    │  Parse Employee ID  │
            │  "not permitted"│    │  from URL           │
            └─────────────────┘    └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Parse JSON Body    │
                                   │  {manager_id: UUID} │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Employee ID ==     │
                                   │  Manager ID?        │
                                   └──────────┬──────────┘
                                              │
                                  ┌───────────┴───────────┐
                                  │ YES                   │ NO
                                  ▼                       ▼
                        ┌─────────────────┐    ┌─────────────────────┐
                        │  Return 400     │    │  Query: Check       │
                        │  "cannot assign │    │  Manager EXISTS     │
                        │  self"          │    └──────────┬──────────┘
                        └─────────────────┘               │
                                                          ▼
                                               ┌─────────────────────┐
                                               │  Manager Exists?    │
                                               └──────────┬──────────┘
                                                          │
                                              ┌───────────┴───────────┐
                                              │ NO                    │ YES
                                              ▼                       ▼
                                    ┌─────────────────┐    ┌─────────────────────┐
                                    │  Return 404     │    │  UPDATE             │
                                    │  "manager not   │    │  Tbl_Employee       │
                                    │  found"         │    │  SET manager_id = ? │
                                    └─────────────────┘    │  WHERE id = ?       │
                                                           └──────────┬──────────┘
                                                                      │
                                                                      ▼
                                                           ┌─────────────────────┐
                                                           │  Update Successful? │
                                                           └──────────┬──────────┘
                                                                      │
                                                          ┌───────────┴───────────┐
                                                          │ NO                    │ YES
                                                          ▼                       ▼
                                                ┌─────────────────┐    ┌─────────────────────┐
                                                │  Return 500     │    │  Return 200         │
                                                │  "failed"       │    │  {                  │
                                                └─────────────────┘    │    message:         │
                                                                       │    "manager         │
                                                                       │    updated",        │
                                                                       │    employee_id,     │
                                                                       │    manager_id       │
                                                                       │  }                  │
                                                                       └──────────┬──────────┘
                                                                                  │
                                                                                  ▼
                                                                                 END
```

---

### 6. GET /api/employee/:id/reports - Get Employee Reports

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Auth Middleware         │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Extract Employee ID     │
                    │  from URL Param          │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Return 200              │
                    │  {message: "Get employee │
                    │  reports"}               │
                    │  (Placeholder)           │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                                END

Note: This endpoint is currently a placeholder and needs full implementation.
```

---

## Leave Routes

### 7. POST /api/leaves/apply - Apply for Leave

```
                                START
                                  │
                                  ▼
                    ┌──────────────────────────┐
                    │  Auth Middleware         │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  Extract user_id & role  │
                    │  from JWT Claims         │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                          ┌──────────────┐
                          │  Role ==     │
                          │  EMPLOYEE?   │
                          └──────┬───────┘
                                 │
                     ┌───────────┴───────────┐
                     │ NO                    │ YES
                     ▼                       ▼
            ┌─────────────────┐    ┌─────────────────────┐
            │  Return 403     │    │  Parse JSON Body    │
            │  "Only employees│    │  {leave_type_id,    │
            │  can apply      │    │  start_date,        │
            │  leave"         │    │  end_date}          │
            └─────────────────┘    └──────────┬──────────┘
                                              │
                                              ▼
                                   ┌─────────────────────┐
                                   │  Validate Input     │
                                   └──────────┬──────────┘
                                              │
                                              ▼
                                       ┌──────────────┐
                                       │  Valid?      │
                                       └──────┬───────┘
                                              │
                                  ┌───────────┴───────────┐
                                  │ NO                    │ YES
                                  ▼                       ▼
                        ┌─────────────────┐    ┌─────────────────────┐
                        │  Return 400     │    │  BEGIN TRANSACTION  │
                        │  "Invalid input"│    └──────────┬──────────┘
                        └─────────────────┘               │
                                                          ▼
                                               ┌─────────────────────┐
                                               │  Query: Fetch       │
                                               │  Manager ID         │
                                               │  WHERE emp_id = ?   │
                                               └──────────┬──────────┘
                                                          │
                                                          ▼
                                                   ┌──────────────┐
                                                   │  Has Manager?│
                                                   └──────┬───────┘
                                                          │
                                              ┌───────────┴───────────┐
                                              │ NO                    │ YES
                                              ▼                       ▼
                                    ┌─────────────────┐    ┌─────────────────────┐
                                    │  ROLLBACK       │    │  Calculate Working  │
                                    │  Return 400     │    │  Days:              │
                                    │  "Manager not   │    │  - Skip Weekends    │
                                    │  assigned"      │    │  - Skip Holidays    │
                                    └─────────────────┘    │  - Count Mon-Fri    │
                                                           └──────────┬──────────┘
                                                                      │
                                                                      ▼
                                                           ┌─────────────────────┐
                                                           │  Days > 0?          │
                                                           └──────────┬──────────┘
                                                                      │
                                                          ┌───────────┴───────────┐
                                                          │ NO                    │ YES
                                                          ▼                       ▼
                                                ┌─────────────────┐    ┌─────────────────────┐
                                                │  ROLLBACK       │    │  Query: Validate    │
                                                │  Return 400     │    │  Leave Type EXISTS  │
                                                │  "Days must be  │    └──────────┬──────────┘
                                                │  > 0"           │               │
                                                └─────────────────┘               ▼
                                                                       ┌─────────────────────┐
                                                                       │  Leave Type Valid?  │
                                                                       └──────────┬──────────┘
                                                                                  │
                                                                      ┌───────────┴───────────┐
                                                                      │ NO                    │ YES
                                                                      ▼                       ▼
                                                            ┌─────────────────┐    ┌─────────────────────┐
                                                            │  ROLLBACK       │    │  Query: Get/Create  │
                                                            │  Return 400     │    │  Leave Balance for  │
                                                            │  "Invalid leave │    │  Current Year       │
                                                            │  type"          │    └──────────┬──────────┘
                                                            └─────────────────┘               │
                                                                                              ▼
                                                                                   ┌─────────────────────┐
                                                                                   │  Balance Exists?    │
                                                                                   └──────────┬──────────┘
                                                                                              │
                                                                                  ┌───────────┴───────────┐
                                                                                  │ NO                    │ YES
                                                                                  ▼                       ▼
                                                                        ┌─────────────────┐    ┌─────────────────────┐
                                                                        │  INSERT Balance │    │  Check Sufficient   │
                                                                        │  with Default   │    │  Balance?           │
                                                                        │  Entitlement    │    │  (closing >= days)  │
                                                                        └────────┬────────┘    └──────────┬──────────┘
                                                                                 │                        │
                                                                                 └────────────┬───────────┘
                                                                                              │
                                                                                  ┌───────────┴───────────┐
                                                                                  │ NO                    │ YES
                                                                                  ▼                       ▼
                                                                        ┌─────────────────┐    ┌─────────────────────┐
                                                                        │  ROLLBACK       │    │  Query: Check       │
                                                                        │  Return 400     │    │  Overlapping Leaves │
                                                                        │  "Insufficient  │    │  (Pending/Approved) │
                                                                        │  balance"       │    └──────────┬──────────┘
                                                                        └─────────────────┘               │
                                                                                                          ▼
                                                                                               ┌─────────────────────┐
                                                                                               │  Overlap Found?     │
                                                                                               └──────────┬──────────┘
                                                                                                          │
                                                                                              ┌───────────┴───────────┐
                                                                                              │ YES                   │ NO
                                                                                              ▼                       ▼
                                                                                    ┌─────────────────┐    ┌─────────────────────┐
                                                                                    │  ROLLBACK       │    │  INSERT INTO        │
                                                                                    │  Return 400     │    │  Tbl_Leave          │
                                                                                    │  "Overlapping   │    │  (status: Pending)  │
                                                                                    │  leave exists"  │    └──────────┬──────────┘
                                                                                    └─────────────────┘               │
                                                                                                                      ▼
                                                                                                           ┌─────────────────────┐
                                                                                                           │  COMMIT TRANSACTION │
                                                                                                           └──────────┬──────────┘
                                                                                                                      │
                                                                                                                      ▼
                                                                                                           ┌─────────────────────┐
                                                                                                           │  Spawn Goroutine    │
                                                                                                           │  (Async Email)      │
                                                                                                           └──────────┬──────────┘
                                                                                                                      │
                                                                                                                      ▼
                                                                                                           ┌─────────────────────┐
                                                                                                           │  Fetch:             │
                                                                                                           │  - Employee Name    │
                                                                                                           │  - Leave Type Name  │
                                                                                                           │  - Manager Email    │
                                                                                                           │  - Admin Emails     │
                                                                                                           └──────────┬──────────┘
                                                                                                                      │
                                                                                                                      ▼
                                                                                                           ┌─────────────────────┐
                                                                                                           │  Send Notification  │
                                                                                                           │  Emails to:         │
                                                                                                           │  - Manager          │
                                                                                                           │  - All Admins       │
                                                                                                           │  - All SuperAdmins  │
                                                                                                           └──────────┬──────────┘
                                                                                                                      │
                                                                                                                      ▼
                                                                                                           ┌─────────────────────┐
                                                                                                           │  Return 200         │
                                                                                                           │  {                  │
                                                                                                           │    message: "Leave  │
                                                                                                           │    applied          │
                                                                                                           │    successfully",   │
                                                                                                           │    leave_id: UUID,  │
                                                                                                           │    days: 5          │
                                                                                                           │  }                  │
                                                                                                           └──────────┬──────────┘
                                                                                                                      │
                                                                                                                      ▼
                                                                                                                     END
```
