# Employee Dashboard API Documentation

## Overview
The Employee Dashboard provides a comprehensive view of an employee's leave information, including statistics, balances, upcoming leaves, and complete leave history with powerful filtering capabilities.

---

## Endpoint

**GET** `/api/leaves/dashboard`

**Authentication Required:** Yes (JWT Token)  
**Role Required:** EMPLOYEE only

---

## Features

✅ **Leave Statistics** - Total, pending, approved, rejected leaves  
✅ **Leave Balances** - Current balance for all leave types  
✅ **Upcoming Leaves** - Next 5 upcoming approved/pending leaves  
✅ **Complete Leave History** - All leaves with filtering  
✅ **Powerful Filters** - Filter by status, type, date range, year  

---

## Request

### Headers
```
Authorization: Bearer <employee_jwt_token>
```

### Query Parameters (All Optional)

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `status` | string | Filter by leave status | `Pending`, `APPROVED`, `REJECTED` |
| `leave_type` | string | Filter by leave type (partial match) | `Annual`, `Sick` |
| `start_date` | string | Filter leaves starting from date | `2024-01-01` |
| `end_date` | string | Filter leaves ending before date | `2024-12-31` |
| `year` | string | Filter by year | `2024` |

---

## Response Structure

### Success Response (200 OK)

```json
{
  "employee_id": "UUID",
  "statistics": {
    "total_leaves": 0,
    "pending_leaves": 0,
    "approved_leaves": 0,
    "rejected_leaves": 0,
    "total_days_used": 0.0,
    "current_year": 2024
  },
  "leave_balances": [
    {
      "leave_type": "string",
      "used": 0.0,
      "total": 0,
      "available": 0.0
    }
  ],
  "upcoming_leaves": [
    {
      "id": "UUID",
      "employee": "string",
      "leave_type": "string",
      "start_date": "timestamp",
      "end_date": "timestamp",
      "days": 0,
      "status": "string",
      "applying_date": "timestamp"
    }
  ],
  "leaves": {
    "total": 0,
    "data": [
      {
        "id": "UUID",
        "employee": "string",
        "leave_type": "string",
        "start_date": "timestamp",
        "end_date": "timestamp",
        "days": 0,
        "status": "string",
        "applying_date": "timestamp"
      }
    ]
  },
  "filters_applied": {
    "status": "string",
    "leave_type": "string",
    "start_date": "string",
    "end_date": "string",
    "year": "string"
  }
}
```

---

## Response Fields Explained

### Statistics Object
- **total_leaves**: Total number of leave applications in current year
- **pending_leaves**: Number of leaves awaiting approval
- **approved_leaves**: Number of approved leaves
- **rejected_leaves**: Number of rejected leaves
- **total_days_used**: Total days of approved leaves taken
- **current_year**: Current year for statistics

### Leave Balances Array
- **leave_type**: Name of leave type (Annual, Sick, etc.)
- **used**: Number of days used
- **total**: Total entitlement for the year
- **available**: Remaining balance (total - used + adjustments)

### Upcoming Leaves Array
- Shows next 5 upcoming leaves (start_date >= today)
- Only includes Pending and APPROVED leaves
- Sorted by start_date (earliest first)

### Leaves Object
- **total**: Count of leaves matching filters
- **data**: Array of leave records matching filters

### Filters Applied Object
- Shows which filters were applied in the request
- Empty string means filter was not applied

---

## Example Responses

### Example 1: Basic Dashboard (No Filters)

**Request:**
```bash
curl -X GET http://localhost:8080/api/leaves/dashboard \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response:**
```json
{
  "employee_id": "550e8400-e29b-41d4-a716-446655440000",
  "statistics": {
    "total_leaves": 12,
    "pending_leaves": 3,
    "approved_leaves": 8,
    "rejected_leaves": 1,
    "total_days_used": 25.0,
    "current_year": 2024
  },
  "leave_balances": [
    {
      "leave_type": "Annual Leave",
      "used": 15.0,
      "total": 20,
      "available": 5.0
    },
    {
      "leave_type": "Sick Leave",
      "used": 5.0,
      "total": 10,
      "available": 5.0
    },
    {
      "leave_type": "Casual Leave",
      "used": 5.0,
      "total": 12,
      "available": 7.0
    }
  ],
  "upcoming_leaves": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "employee": "John Doe",
      "leave_type": "Annual Leave",
      "start_date": "2024-12-15T00:00:00Z",
      "end_date": "2024-12-20T00:00:00Z",
      "days": 5,
      "status": "APPROVED",
      "applying_date": "2024-11-20T10:00:00Z"
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "employee": "John Doe",
      "leave_type": "Sick Leave",
      "start_date": "2024-12-25T00:00:00Z",
      "end_date": "2024-12-26T00:00:00Z",
      "days": 2,
      "status": "Pending",
      "applying_date": "2024-11-25T14:30:00Z"
    }
  ],
  "leaves": {
    "total": 12,
    "data": [
      {
        "id": "990e8400-e29b-41d4-a716-446655440004",
        "employee": "John Doe",
        "leave_type": "Annual Leave",
        "start_date": "2024-11-01T00:00:00Z",
        "end_date": "2024-11-05T00:00:00Z",
        "days": 5,
        "status": "APPROVED",
        "applying_date": "2024-10-25T09:00:00Z"
      }
    ]
  },
  "filters_applied": {
    "status": "",
    "leave_type": "",
    "start_date": "",
    "end_date": "",
    "year": ""
  }
}
```

---

### Example 2: Filter by Status (Approved Only)

**Request:**
```bash
curl -X GET "http://localhost:8080/api/leaves/dashboard?status=APPROVED" \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "employee_id": "550e8400-e29b-41d4-a716-446655440000",
  "statistics": {
    "total_leaves": 12,
    "pending_leaves": 3,
    "approved_leaves": 8,
    "rejected_leaves": 1,
    "total_days_used": 25.0,
    "current_year": 2024
  },
  "leave_balances": [...],
  "upcoming_leaves": [...],
  "leaves": {
    "total": 8,
    "data": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "employee": "John Doe",
        "leave_type": "Annual Leave",
        "start_date": "2024-12-15T00:00:00Z",
        "end_date": "2024-12-20T00:00:00Z",
        "days": 5,
        "status": "APPROVED",
        "applying_date": "2024-11-20T10:00:00Z"
      }
    ]
  },
  "filters_applied": {
    "status": "APPROVED",
    "leave_type": "",
    "start_date": "",
    "end_date": "",
    "year": ""
  }
}
```

---

### Example 3: Filter by Leave Type

**Request:**
```bash
curl -X GET "http://localhost:8080/api/leaves/dashboard?leave_type=Annual" \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "employee_id": "550e8400-e29b-41d4-a716-446655440000",
  "statistics": {...},
  "leave_balances": [...],
  "upcoming_leaves": [...],
  "leaves": {
    "total": 6,
    "data": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "employee": "John Doe",
        "leave_type": "Annual Leave",
        "start_date": "2024-12-15T00:00:00Z",
        "end_date": "2024-12-20T00:00:00Z",
        "days": 5,
        "status": "APPROVED",
        "applying_date": "2024-11-20T10:00:00Z"
      }
    ]
  },
  "filters_applied": {
    "status": "",
    "leave_type": "Annual",
    "start_date": "",
    "end_date": "",
    "year": ""
  }
}
```

---

### Example 4: Filter by Date Range

**Request:**
```bash
curl -X GET "http://localhost:8080/api/leaves/dashboard?start_date=2024-06-01&end_date=2024-12-31" \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "employee_id": "550e8400-e29b-41d4-a716-446655440000",
  "statistics": {...},
  "leave_balances": [...],
  "upcoming_leaves": [...],
  "leaves": {
    "total": 7,
    "data": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "employee": "John Doe",
        "leave_type": "Annual Leave",
        "start_date": "2024-12-15T00:00:00Z",
        "end_date": "2024-12-20T00:00:00Z",
        "days": 5,
        "status": "APPROVED",
        "applying_date": "2024-11-20T10:00:00Z"
      }
    ]
  },
  "filters_applied": {
    "status": "",
    "leave_type": "",
    "start_date": "2024-06-01",
    "end_date": "2024-12-31",
    "year": ""
  }
}
```

---

### Example 5: Multiple Filters

**Request:**
```bash
curl -X GET "http://localhost:8080/api/leaves/dashboard?status=APPROVED&year=2024&leave_type=Annual" \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
{
  "employee_id": "550e8400-e29b-41d4-a716-446655440000",
  "statistics": {...},
  "leave_balances": [...],
  "upcoming_leaves": [...],
  "leaves": {
    "total": 4,
    "data": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "employee": "John Doe",
        "leave_type": "Annual Leave",
        "start_date": "2024-12-15T00:00:00Z",
        "end_date": "2024-12-20T00:00:00Z",
        "days": 5,
        "status": "APPROVED",
        "applying_date": "2024-11-20T10:00:00Z"
      }
    ]
  },
  "filters_applied": {
    "status": "APPROVED",
    "leave_type": "Annual",
    "start_date": "",
    "end_date": "",
    "year": "2024"
  }
}
```

---

## Error Responses

### 403 Forbidden
**Scenario:** Non-employee tries to access dashboard

```json
{
  "error": {
    "code": 403,
    "message": "This dashboard is only for employees"
  }
}
```

### 401 Unauthorized
**Scenario:** Missing or invalid JWT token

```json
{
  "error": {
    "code": 401,
    "message": "Missing Authorization header"
  }
}
```

### 500 Internal Server Error
**Scenario:** Database error

```json
{
  "error": {
    "code": 500,
    "message": "Failed to fetch leaves: <error details>"
  }
}
```

---

## Use Cases

### 1. Employee Login Dashboard
Show employee their complete leave overview on login:
```bash
GET /api/leaves/dashboard
```

### 2. View Pending Leaves
Employee wants to see all pending leave requests:
```bash
GET /api/leaves/dashboard?status=Pending
```

### 3. Check Annual Leave History
Employee wants to see all annual leaves taken:
```bash
GET /api/leaves/dashboard?leave_type=Annual&status=APPROVED
```

### 4. Year-End Review
Employee wants to review all leaves for the year:
```bash
GET /api/leaves/dashboard?year=2024
```

### 5. Check Recent Leaves
Employee wants to see leaves from last 3 months:
```bash
GET /api/leaves/dashboard?start_date=2024-09-01&end_date=2024-11-30
```

---

## Frontend Integration Examples

### React/JavaScript Example

```javascript
// Fetch employee dashboard
const fetchDashboard = async (filters = {}) => {
  const queryParams = new URLSearchParams(filters).toString();
  const url = `http://localhost:8080/api/leaves/dashboard${queryParams ? '?' + queryParams : ''}`;
  
  const response = await fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`,
      'Content-Type': 'application/json'
    }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch dashboard');
  }
  
  return await response.json();
};

// Usage examples
// No filters
const dashboard = await fetchDashboard();

// With filters
const filteredDashboard = await fetchDashboard({
  status: 'APPROVED',
  year: '2024'
});
```

### React Component Example

```jsx
import React, { useState, useEffect } from 'react';

function EmployeeDashboard() {
  const [dashboard, setDashboard] = useState(null);
  const [filters, setFilters] = useState({
    status: '',
    leave_type: '',
    year: new Date().getFullYear()
  });

  useEffect(() => {
    fetchDashboard(filters).then(setDashboard);
  }, [filters]);

  return (
    <div className="dashboard">
      {/* Statistics Cards */}
      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Leaves</h3>
          <p>{dashboard?.statistics.total_leaves}</p>
        </div>
        <div className="stat-card">
          <h3>Pending</h3>
          <p>{dashboard?.statistics.pending_leaves}</p>
        </div>
        <div className="stat-card">
          <h3>Approved</h3>
          <p>{dashboard?.statistics.approved_leaves}</p>
        </div>
        <div className="stat-card">
          <h3>Days Used</h3>
          <p>{dashboard?.statistics.total_days_used}</p>
        </div>
      </div>

      {/* Leave Balances */}
      <div className="balances">
        <h2>Leave Balances</h2>
        {dashboard?.leave_balances.map(balance => (
          <div key={balance.leave_type} className="balance-item">
            <span>{balance.leave_type}</span>
            <span>{balance.available} / {balance.total} available</span>
          </div>
        ))}
      </div>

      {/* Filters */}
      <div className="filters">
        <select 
          value={filters.status} 
          onChange={(e) => setFilters({...filters, status: e.target.value})}
        >
          <option value="">All Status</option>
          <option value="Pending">Pending</option>
          <option value="APPROVED">Approved</option>
          <option value="REJECTED">Rejected</option>
        </select>
      </div>

      {/* Leave List */}
      <div className="leaves-list">
        {dashboard?.leaves.data.map(leave => (
          <div key={leave.id} className="leave-item">
            <span>{leave.leave_type}</span>
            <span>{leave.start_date} - {leave.end_date}</span>
            <span className={`status-${leave.status.toLowerCase()}`}>
              {leave.status}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## Performance Considerations

- Dashboard data is fetched in a single request (efficient)
- Filters are applied at database level (fast)
- Upcoming leaves limited to 5 records
- All queries use indexed columns for performance
- Statistics calculated in single query using aggregation

---

## Security

- ✅ JWT authentication required
- ✅ Role-based access (EMPLOYEE only)
- ✅ Employee can only see their own data
- ✅ SQL injection protected (parameterized queries)
- ✅ No sensitive data exposed

---

**Last Updated**: November 26, 2024  
**Version**: 1.0  
**Endpoint**: `/api/leaves/dashboard`
