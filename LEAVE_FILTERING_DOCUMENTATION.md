# Leave Filtering Implementation

## Overview
Enhanced the `GetAllLeaves` endpoint to support month and year filtering with default values.

## API Endpoints

### Get All Leaves with Filtering
**Endpoint:** `GET /api/leaves/all`

**Query Parameters:**
- `month` (optional): Filter leaves by month (1-12). Default: current month
- `year` (optional): Filter leaves by year (2000-2100). Default: current year

**Examples:**

1. **Get leaves for current month and year (default):**
   ```
   GET /api/leaves/all
   ```

2. **Get leaves for specific month and year:**
   ```
   GET /api/leaves/all?month=12&year=2024
   ```

3. **Get leaves for specific month (current year):**
   ```
   GET /api/leaves/all?month=11
   ```

4. **Get leaves for specific year (current month):**
   ```
   GET /api/leaves/all?year=2023
   ```

## Response Format

```json
{
  "message": "Leaves fetched successfully",
  "total": 5,
  "role": "ADMIN",
  "month": 12,
  "year": 2024,
  "data": [
    {
      "id": "uuid",
      "employee": "John Doe",
      "leave_type": "Annual Leave",
      "is_paid": true,
      "leave_timing_type": "FULL",
      "leave_timing": "Full Day",
      "start_date": "2024-12-15T00:00:00Z",
      "end_date": "2024-12-17T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "status": "APPROVED",
      "applied_at": "2024-12-01T10:30:00Z"
    }
  ]
}
```

## Role-Based Access Control

The filtering works with existing role-based access:

- **EMPLOYEE**: Can see only their own leaves for the specified month/year
- **MANAGER**: Can see their own leaves + team members' leaves for the specified month/year
- **ADMIN/HR/SUPERADMIN**: Can see all leaves for the specified month/year

## Database Changes

### New Repository Methods Added:

1. `GetAllEmployeeLeaveByMonthYear(userID, month, year)` - Employee leaves filtered by month/year
2. `GetAllleavebaseonassignManagerByMonthYear(userID, month, year)` - Manager team leaves filtered by month/year  
3. `GetAllLeaveByMonthYear(month, year)` - All leaves filtered by month/year (Admin/HR/SuperAdmin)

### SQL Filtering Logic:
```sql
WHERE EXTRACT(MONTH FROM l.start_date) = $month
AND EXTRACT(YEAR FROM l.start_date) = $year
```

## Validation

- **Month**: Must be between 1-12
- **Year**: Must be between 2000-2100
- Invalid parameters return 400 Bad Request with descriptive error message

## Backward Compatibility

- Existing API calls without query parameters work exactly as before
- Default behavior shows current month and year data
- No breaking changes to existing functionality

## Error Handling

- Invalid month: `"Invalid month. Must be between 1-12"`
- Invalid year: `"Invalid year. Must be between 2000-2100"`
- Database errors: `"Failed to fetch leaves: [error details]"`

## Implementation Details

- Uses Go's `time.Now()` to get current month/year as defaults
- Leverages PostgreSQL's `EXTRACT()` function for date filtering
- Maintains existing role-based security and access control
- Preserves all existing sorting (ORDER BY created_at DESC)