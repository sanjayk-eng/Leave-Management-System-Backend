# Leave API Usage Examples

## Testing the Enhanced Leave Filtering

### 1. Get Current Month's Leaves (Default Behavior)
```bash
curl -X GET "http://localhost:8080/api/leaves/all" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "message": "Leaves fetched successfully",
  "total": 3,
  "role": "ADMIN",
  "month": 1,
  "year": 2025,
  "data": [...]
}
```

### 2. Get Leaves for December 2024
```bash
curl -X GET "http://localhost:8080/api/leaves/all?month=12&year=2024" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Get Leaves for November (Current Year)
```bash
curl -X GET "http://localhost:8080/api/leaves/all?month=11" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 4. Get Leaves for 2023 (Current Month)
```bash
curl -X GET "http://localhost:8080/api/leaves/all?year=2023" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Error Cases

**Invalid Month:**
```bash
curl -X GET "http://localhost:8080/api/leaves/all?month=13" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
Response: `400 Bad Request - "Invalid month. Must be between 1-12"`

**Invalid Year:**
```bash
curl -X GET "http://localhost:8080/api/leaves/all?year=1999" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
Response: `400 Bad Request - "Invalid year. Must be between 2000-2100"`

## Role-Based Results

### Employee Role
- Only sees their own leaves for the specified month/year
- Cannot see other employees' leaves

### Manager Role  
- Sees their own leaves + team members' leaves for the specified month/year
- Limited to their reporting hierarchy

### Admin/HR/SuperAdmin Role
- Sees all company leaves for the specified month/year
- Full access to all employee data

## Frontend Integration Examples

### JavaScript/Fetch
```javascript
// Get current month leaves
const getCurrentMonthLeaves = async () => {
  const response = await fetch('/api/leaves/all', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};

// Get specific month/year leaves
const getLeavesByMonthYear = async (month, year) => {
  const response = await fetch(`/api/leaves/all?month=${month}&year=${year}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};
```

### React Hook Example
```javascript
const useLeaves = (month, year) => {
  const [leaves, setLeaves] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLeaves = async () => {
      try {
        setLoading(true);
        const params = new URLSearchParams();
        if (month) params.append('month', month);
        if (year) params.append('year', year);
        
        const response = await fetch(`/api/leaves/all?${params}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        
        const data = await response.json();
        setLeaves(data.data);
      } catch (error) {
        console.error('Failed to fetch leaves:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchLeaves();
  }, [month, year]);

  return { leaves, loading };
};
```

## Database Query Performance

The filtering uses PostgreSQL's `EXTRACT()` function which is efficient for date-based filtering:

```sql
-- Optimized query with proper indexing
WHERE EXTRACT(MONTH FROM l.start_date) = $month
AND EXTRACT(YEAR FROM l.start_date) = $year
```

**Recommended Index:**
```sql
CREATE INDEX idx_leave_start_date_month_year 
ON Tbl_Leave (EXTRACT(YEAR FROM start_date), EXTRACT(MONTH FROM start_date));
```