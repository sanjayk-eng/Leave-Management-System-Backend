# Get My Leave API Documentation

## Overview
The `GetAllMyLeave` endpoint allows users to retrieve their own leave applications with optional month and year filtering.

## API Endpoint

### GET /api/leaves/my

**Description:** Get current user's own leaves with month/year filtering

**Authentication:** Required (Bearer Token)

**Authorization:** All authenticated users can access their own leaves

## Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `month` | integer | No | Current month | Month to filter (1-12) |
| `year` | integer | No | Current year | Year to filter (2000-2100) |

## Request Examples

### Get current month's leaves
```bash
GET /api/leaves/my
```

### Get leaves for specific month and year
```bash
GET /api/leaves/my?month=11&year=2024
```

### Get leaves for specific month (current year)
```bash
GET /api/leaves/my?month=12
```

### Get leaves for specific year (current month)
```bash
GET /api/leaves/my?year=2023
```

## Response Format

### Success Response (200 OK)
```json
{
  "message": "My leaves fetched successfully",
  "total": 3,
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "month": 11,
  "year": 2024,
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "employee": "John Doe",
      "leave_type": "Annual Leave",
      "is_paid": true,
      "leave_timing_type": "FULL",
      "leave_timing": "Full Day",
      "start_date": "2024-11-15T00:00:00Z",
      "end_date": "2024-11-17T00:00:00Z",
      "days": 3.0,
      "reason": "Family vacation",
      "status": "APPROVED",
      "applied_at": "2024-11-01T10:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174002",
      "employee": "John Doe",
      "leave_type": "Sick Leave",
      "is_paid": true,
      "leave_timing_type": "HALF",
      "leave_timing": "First Half",
      "start_date": "2024-11-05T00:00:00Z",
      "end_date": "2024-11-05T00:00:00Z",
      "days": 0.5,
      "reason": "Medical appointment",
      "status": "APPROVED",
      "applied_at": "2024-11-04T09:15:00Z"
    }
  ]
}
```

### Error Responses

#### 401 Unauthorized
```json
{
  "error": "User ID not found in context"
}
```

#### 400 Bad Request
```json
{
  "error": "Invalid month. Must be between 1-12"
}
```

```json
{
  "error": "Invalid year. Must be between 2000-2100"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to fetch my leaves: database connection error"
}
```

## Features

1. **User-Specific**: Only returns leaves for the authenticated user
2. **Month/Year Filtering**: Filter leaves by specific month and year
3. **Default Values**: Uses current month/year if not specified
4. **Comprehensive Data**: Returns complete leave information including timing details
5. **Sorted Results**: Results are ordered by application date (newest first)

## Differences from GetAllLeaves

| Feature | GetAllLeaves | GetAllMyLeave |
|---------|--------------|---------------|
| **Scope** | Role-based (can see team/all leaves) | User-specific only |
| **Authorization** | Role-dependent access control | Any authenticated user |
| **Data Returned** | Varies by role (own/team/all) | Only current user's leaves |
| **Use Case** | Admin/Manager oversight | Personal leave history |

## Security

- **Authentication Required**: Must provide valid Bearer token
- **User Isolation**: Users can only see their own leave records
- **Input Validation**: Month and year parameters are validated
- **SQL Injection Protection**: Uses parameterized queries

## Performance

- **Indexed Queries**: Uses employee_id and date-based filtering for optimal performance
- **Limited Scope**: Only queries user's own records, reducing data load
- **Efficient Joins**: Optimized JOIN operations with proper table relationships