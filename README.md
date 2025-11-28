# Zenithive - User Management System Backend

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-316192?style=for-the-badge&logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success?style=for-the-badge)

**A comprehensive employee management system with leave tracking, payroll processing, and role-based access control**

[Features](#-features) • [Quick Start](#-quick-start) • [API Documentation](#-api-documentation) • [Architecture](#-architecture) • [Contributing](#-contributing)

</div>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [System Architecture](#-system-architecture)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Security](#-security)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🌟 Overview

Zenithive is a modern, enterprise-grade user management system designed for organizations with 25-100 employees. It provides comprehensive employee lifecycle management, automated leave tracking, intelligent payroll processing, and robust role-based access control.

### Key Highlights

- 🔐 **Secure Authentication** - JWT-based authentication with bcrypt password hashing
- 👥 **Employee Management** - Complete CRUD operations with hierarchy management
- 🏖️ **Leave Management** - Automated leave tracking with approval workflows
- 💰 **Payroll Processing** - Intelligent salary calculation with deduction management
- 📊 **Role-Based Access** - Granular permissions across 5 user roles
- 📧 **Email Notifications** - Automated notifications for all critical events
- 📄 **PDF Generation** - Professional payslip generation
- 🔄 **RESTful API** - Clean, well-documented REST API

---

## ✨ Features

### Employee Management
- ✅ Create, read, update, and deactivate employees
- ✅ Manager hierarchy with team management
- ✅ Role assignment and management
- ✅ Password management with secure hashing
- ✅ Employee profile with joining date, salary, and status
- ✅ Email domain validation (@zenithive.com)

### Leave Management
- ✅ Leave application with reason validation
- ✅ Multi-level approval workflow
- ✅ Working days calculation (excludes weekends & holidays)
- ✅ Leave balance tracking per employee
- ✅ Leave cancellation (pending leaves)
- ✅ Leave withdrawal (approved leaves)
- ✅ Admin can add leave on behalf of employees
- ✅ Multiple leave types (Annual, Sick, etc.)
- ✅ Leave policy management

### Payroll Management
- ✅ Monthly payroll processing
- ✅ Automatic deduction calculation based on leaves
- ✅ Payroll preview before finalization
- ✅ Professional PDF payslip generation
- ✅ Payslip download for employees
- ✅ SUPERADMIN-only finalization for security

### Access Control
- ✅ 5 distinct roles: SUPERADMIN, ADMIN, HR, MANAGER, EMPLOYEE
- ✅ Granular permissions per endpoint
- ✅ Self-modification restrictions
- ✅ Manager hierarchy validation
- ✅ JWT token-based authentication

### Notifications
- ✅ Email notifications for 8+ events
- ✅ Welcome emails for new employees
- ✅ Leave status notifications
- ✅ Password change alerts
- ✅ Async email processing

### Company Settings
- ✅ Configurable working days per month
- ✅ Manager leave addition toggle
- ✅ Holiday management
- ✅ Company-wide settings

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL 14+
- **ORM**: sqlx (SQL extensions)
- **Authentication**: JWT (golang-jwt)
- **Password Hashing**: bcrypt
- **PDF Generation**: gofpdf
- **Validation**: go-playground/validator

### Tools & Libraries
- **CORS**: gin-contrib/cors
- **UUID**: google/uuid
- **Environment**: godotenv

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                          │
│                    (React Frontend)                          │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP/REST API
                         │ JWT Authentication
┌────────────────────────▼────────────────────────────────────┐
│                     API Gateway Layer                        │
│                    (Gin Router + CORS)                       │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Middleware Layer                           │
│              (Auth, Validation, Logging)                     │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Controller Layer                           │
│   ┌──────────┬──────────┬──────────┬──────────┬─────────┐  │
│   │ Employee │  Leave   │ Payroll  │ Settings │  Auth   │  │
│   └──────────┴──────────┴──────────┴──────────┴─────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                  Repository Layer                            │
│              (Database Queries & Logic)                      │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Database Layer                             │
│                   (PostgreSQL)                               │
│   ┌──────────────────────────────────────────────────┐     │
│   │  Tables: Employee, Leave, Payroll, Settings      │     │
│   └──────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### Project Structure
```
UserMenagmentSystem_Backend/
├── controllers/          # Request handlers
│   ├── auth.go          # Authentication
│   ├── employee.go      # Employee management
│   ├── leave.go         # Leave management
│   ├── leave_balance.go # Leave balance operations
│   ├── payroll.go       # Payroll processing
│   └── settings.go      # Company settings
├── middlewere/          # Middleware functions
│   └── middlewere.go    # Auth middleware
├── models/              # Data models
│   └── models.go        # All data structures
├── repositories/        # Database layer
│   └── repo.go          # Database queries
├── routes/              # Route definitions
│   └── router.go        # API routes
├── utils/               # Utility functions
│   ├── auth.go          # JWT & password utils
│   ├── errors.go        # Error handling
│   └── notification.go  # Email notifications
├── pkg/                 # Packages
│   ├── config/          # Configuration
│   ├── database/        # Database connection
│   └── migration/       # SQL migrations
├── tmp/                 # Temporary files (PDFs)
├── .env                 # Environment variables
├── main.go              # Application entry point
├── go.mod               # Go dependencies
└── README.md            # This file
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/sanjayk-eng/UserMenagmentSystem_Backend.git
cd UserMenagmentSystem_Backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup PostgreSQL database**
```bash
# Create database
createdb user_management_db

# Run migrations
psql -d user_management_db -f pkg/migration/20251120110206_allschima.sql
psql -d user_management_db -f pkg/migration/20251120134716_tblrole.sql
psql -d user_management_db -f pkg/migration/20251123045525_tbl_holiday.sql
psql -d user_management_db -f pkg/migration/20251124053315_tbl_leave_adj_add_col.sql
psql -d user_management_db -f pkg/migration/20251124103513_tbl_setting_info.sql
psql -d user_management_db -f pkg/migration/20251124104539_tbl_setting_add.sql
```

4. **Configure environment variables**
```bash
# Create .env file
cp .env.example .env

# Edit .env with your settings
nano .env
```

**.env Configuration:**
```env
# Server Configuration
APP_PORT=8080
FRONTEND_SERVER=http://localhost:3000

# Security
SERACT_KEY=your_jwt_secret_key_here

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=user_management_db

# Email Service
GOOGLE_SCRIPT_URL=your_email_service_url
```

5. **Run the application**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Quick Test

```bash
# Test the API
curl http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@zenithive.com",
    "password": "admin123"
  }'
```

---

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication
All endpoints (except `/auth/login`) require JWT authentication:
```
Authorization: Bearer <your_jwt_token>
```

### Endpoint Summary

| Module | Endpoints | Description |
|--------|-----------|-------------|
| **Authentication** | 1 | Login and token generation |
| **Employee Management** | 10 | CRUD operations, role management |
| **Leave Management** | 9 | Apply, approve, cancel, withdraw leaves |
| **Leave Balances** | 2 | View and adjust leave balances |
| **Payroll** | 4 | Run, finalize, and download payslips |
| **Settings** | 3 | Company settings and holidays |
| **Total** | **32** | Complete API coverage |

### Quick Examples

**Login:**
```bash
POST /api/auth/login
{
  "email": "user@zenithive.com",
  "password": "password123"
}
```

**Create Employee:**
```bash
POST /api/employee/
Authorization: Bearer <token>
{
  "full_name": "John Doe",
  "email": "john@zenithive.com",
  "role": "EMPLOYEE",
  "password": "temp123",
  "salary": 50000,
  "joining_date": "2024-12-01T00:00:00Z"
}
```

**Apply Leave:**
```bash
POST /api/leaves/apply
Authorization: Bearer <token>
{
  "leave_type_id": 1,
  "start_date": "2024-12-10T00:00:00Z",
  "end_date": "2024-12-12T00:00:00Z",
  "reason": "Family vacation"
}
```

**Run Payroll:**
```bash
POST /api/payroll/run
Authorization: Bearer <token>
{
  "month": 11,
  "year": 2024
}
```

### Complete Documentation

For detailed API documentation, see:
- 📖 [Complete API Documentation](./COMPLETE_API_DOCUMENTATION.md)
- 🚀 [Quick Reference Guide](./QUICK_REFERENCE_GUIDE.md)

---

## 🗄️ Database Schema

### Core Tables

- **Tbl_Employee** - Employee information and hierarchy
- **Tbl_Role** - User roles (SUPERADMIN, ADMIN, HR, MANAGER, EMPLOYEE)
- **Tbl_Leave** - Leave requests and status
- **Tbl_Leave_Type** - Leave policies (Annual, Sick, etc.)
- **Tbl_Leave_Balance** - Employee leave balances
- **Tbl_Leave_Adjustment** - Manual balance adjustments
- **Tbl_Payroll_Run** - Payroll processing records
- **Tbl_Payslip** - Generated payslips
- **Tbl_Holiday** - Company holidays
- **Tbl_Company_Settings** - System configuration

### Entity Relationships

```
Employee ──┬── manages ──> Employee (Manager)
           ├── has ──> Leave
           ├── has ──> Leave_Balance
           └── has ──> Payslip

Leave ──┬── belongs to ──> Employee
        └── has type ──> Leave_Type

Payslip ──┬── belongs to ──> Employee
          └── part of ──> Payroll_Run
```

For detailed schema, see [SCHEMA.md](./SCHEMA.md)

---

## 🔒 Security

### Authentication & Authorization
- ✅ JWT token-based authentication
- ✅ Token expiration and refresh
- ✅ Password hashing with bcrypt (cost factor: 10)
- ✅ Role-based access control (RBAC)
- ✅ Route-level middleware protection

### Data Protection
- ✅ SQL injection prevention (parameterized queries)
- ✅ Email domain validation
- ✅ Password strength requirements (min 6 chars)
- ✅ Sensitive data never exposed in responses
- ✅ CORS configuration for frontend

### Access Control Rules
- ✅ ADMIN/HR cannot change their own role
- ✅ ADMIN/HR cannot modify SUPERADMIN users
- ✅ Only SUPERADMIN can finalize payroll
- ✅ Employees can only view/modify their own data
- ✅ Managers can only manage their team members

### Audit Trail
- ✅ All modifications tracked with timestamps
- ✅ Leave adjustments logged
- ✅ Payroll finalization tracked
- ✅ Email notifications for critical actions

---

## 🧪 Testing

### Manual Testing

```bash
# Set token variable
TOKEN="your_jwt_token_here"

# Test employee endpoints
curl -X GET http://localhost:8080/api/employee/ \
  -H "Authorization: Bearer $TOKEN"

# Test leave endpoints
curl -X GET http://localhost:8080/api/leaves/all \
  -H "Authorization: Bearer $TOKEN"

# Test payroll endpoints
curl -X GET http://localhost:8080/api/payroll/payslip \
  -H "Authorization: Bearer $TOKEN"
```

### Test Data

Default users (after migration):
- **SUPERADMIN**: superadmin@zenithive.com / superadmin123
- **ADMIN**: admin@zenithive.com / admin123
- **MANAGER**: manager@zenithive.com / manager123
- **EMPLOYEE**: employee@zenithive.com / employee123

---

## 🚢 Deployment

### Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
```

```bash
# Build and run
docker build -t zenithive-backend .
docker run -p 8080:8080 zenithive-backend
```

### Docker Compose

```yaml
version: '3.8'
services:
  db:
    image: postgres:14
    environment:
      POSTGRES_DB: user_management_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  backend:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - db
    environment:
      DB_HOST: db
      DB_PORT: 5432

volumes:
  postgres_data:
```

### Production Checklist

- [ ] Set strong JWT secret key
- [ ] Configure HTTPS/TLS
- [ ] Set up database backups
- [ ] Configure email service
- [ ] Set up monitoring and logging
- [ ] Configure CORS for production domain
- [ ] Set up rate limiting
- [ ] Enable database connection pooling
- [ ] Configure environment-specific settings

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Coding Standards

- Follow Go best practices and conventions
- Write clear, descriptive commit messages
- Add comments for complex logic
- Update documentation for API changes
- Test your changes thoroughly

---

## 📞 Support

For issues, questions, or feature requests:

- 📧 Email: support@zenithive.com
- 🐛 Issues: [GitHub Issues](https://github.com/sanjayk-eng/UserMenagmentSystem_Backend/issues)
- 📖 Documentation: [Complete API Docs](./COMPLETE_API_DOCUMENTATION.md)

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Gin Web Framework
- PostgreSQL Database
- Go Community
- All Contributors

---

## 📊 Project Status

- ✅ **Version**: 1.0
- ✅ **Status**: Production Ready
- ✅ **Last Updated**: November 2024
- ✅ **Maintained**: Yes

---

<div align="center">

**Built with ❤️ by the Zenithive Team**

[⬆ Back to Top](#zenithive---user-management-system-backend)

</div>