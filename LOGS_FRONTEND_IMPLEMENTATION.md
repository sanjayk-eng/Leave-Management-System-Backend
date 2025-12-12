# Logs Frontend Implementation Guide

## Overview
This document provides a complete implementation guide for integrating the logs API endpoint into a React frontend application.

## API Endpoint Details

### Base URL
```
GET /api/logs/
```

### Authentication
- Requires JWT token in Authorization header
- Only accessible by users with `SUPERADMIN` role

### Query Parameters
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| days | integer | No | 7 | Number of days to filter logs (must be positive) |

### Request Examples
```bash
# Get logs from last 7 days (default)
GET /api/logs/

# Get logs from last 30 days
GET /api/logs/?days=30

# Get logs from last 1 day
GET /api/logs/?days=1
```

### Response Format
```json
{
  "message": "Logs retrieved successfully",
  "data": {
    "logs": [
      {
        "id": 1,
        "user_name": "John Doe",
        "action": "CREATE_EMPLOYEE",
        "component": "EMPLOYEE",
        "created_at": "2024-12-12T10:30:00Z"
      }
    ],
    "total_count": 1,
    "days_filter": 7,
    "date_from": "2024-12-05"
  }
}
```

### Error Responses
```json
// 403 Forbidden (Non-superadmin user)
{
  "error": "Access denied. Only SUPERADMIN can view logs"
}

// 400 Bad Request (Invalid days parameter)
{
  "error": "Invalid days parameter. Must be a positive integer"
}

// 500 Internal Server Error
{
  "error": "Failed to fetch logs"
}
```

---

## React Implementation

### 1. API Service Layer

Create `src/services/logsService.js`:

```javascript
import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to include auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Logs API functions
export const logsService = {
  // Get logs with optional days filter
  getLogs: async (days = 7) => {
    try {
      const response = await apiClient.get(`/api/logs/`, {
        params: { days }
      });
      return response.data;
    } catch (error) {
      throw error.response?.data || error.message;
    }
  },
};

export default logsService;
```

### 2. Custom Hook for Logs

Create `src/hooks/useLogs.js`:

```javascript
import { useState, useEffect, useCallback } from 'react';
import { logsService } from '../services/logsService';

export const useLogs = (initialDays = 7) => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [totalCount, setTotalCount] = useState(0);
  const [daysFilter, setDaysFilter] = useState(initialDays);
  const [dateFrom, setDateFrom] = useState('');

  const fetchLogs = useCallback(async (days = daysFilter) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await logsService.getLogs(days);
      setLogs(response.data.logs);
      setTotalCount(response.data.total_count);
      setDaysFilter(response.data.days_filter);
      setDateFrom(response.data.date_from);
    } catch (err) {
      setError(err.error || 'Failed to fetch logs');
      setLogs([]);
      setTotalCount(0);
    } finally {
      setLoading(false);
    }
  }, [daysFilter]);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  const refetch = useCallback((days) => {
    fetchLogs(days);
  }, [fetchLogs]);

  return {
    logs,
    loading,
    error,
    totalCount,
    daysFilter,
    dateFrom,
    refetch,
  };
};
```

### 3. Logs Component

Create `src/components/Logs/LogsPage.jsx`:

```javascript
import React, { useState } from 'react';
import { useLogs } from '../../hooks/useLogs';
import LogsTable from './LogsTable';
import LogsFilter from './LogsFilter';
import LoadingSpinner from '../common/LoadingSpinner';
import ErrorMessage from '../common/ErrorMessage';
import './LogsPage.css';

const LogsPage = () => {
  const [selectedDays, setSelectedDays] = useState(7);
  const { logs, loading, error, totalCount, daysFilter, dateFrom, refetch } = useLogs(selectedDays);

  const handleDaysChange = (days) => {
    setSelectedDays(days);
    refetch(days);
  };

  const handleRefresh = () => {
    refetch(selectedDays);
  };

  if (loading) {
    return <LoadingSpinner message="Loading logs..." />;
  }

  return (
    <div className="logs-page">
      <div className="logs-header">
        <h1>System Logs</h1>
        <p className="logs-subtitle">
          View and monitor system activities and user actions
        </p>
      </div>

      <LogsFilter
        selectedDays={selectedDays}
        onDaysChange={handleDaysChange}
        onRefresh={handleRefresh}
        totalCount={totalCount}
        dateFrom={dateFrom}
      />

      {error && (
        <ErrorMessage 
          message={error} 
          onRetry={handleRefresh}
        />
      )}

      {!error && (
        <LogsTable 
          logs={logs} 
          loading={loading}
          totalCount={totalCount}
        />
      )}
    </div>
  );
};

export default LogsPage;
```

### 4. Logs Filter Component

Create `src/components/Logs/LogsFilter.jsx`:

```javascript
import React from 'react';
import './LogsFilter.css';

const LogsFilter = ({ 
  selectedDays, 
  onDaysChange, 
  onRefresh, 
  totalCount, 
  dateFrom 
}) => {
  const dayOptions = [
    { value: 1, label: 'Last 24 hours' },
    { value: 7, label: 'Last 7 days' },
    { value: 30, label: 'Last 30 days' },
    { value: 90, label: 'Last 90 days' },
  ];

  const handleCustomDays = (e) => {
    const value = parseInt(e.target.value);
    if (value > 0) {
      onDaysChange(value);
    }
  };

  return (
    <div className="logs-filter">
      <div className="filter-section">
        <div className="filter-group">
          <label htmlFor="days-select">Time Period:</label>
          <select
            id="days-select"
            value={selectedDays}
            onChange={(e) => onDaysChange(parseInt(e.target.value))}
            className="days-select"
          >
            {dayOptions.map(option => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>

        <div className="filter-group">
          <label htmlFor="custom-days">Custom Days:</label>
          <input
            id="custom-days"
            type="number"
            min="1"
            max="365"
            placeholder="Enter days"
            onBlur={handleCustomDays}
            onKeyPress={(e) => e.key === 'Enter' && handleCustomDays(e)}
            className="custom-days-input"
          />
        </div>

        <button 
          onClick={onRefresh}
          className="refresh-btn"
          title="Refresh logs"
        >
          🔄 Refresh
        </button>
      </div>

      <div className="filter-info">
        <span className="total-count">
          Total Logs: <strong>{totalCount}</strong>
        </span>
        <span className="date-range">
          From: <strong>{dateFrom}</strong>
        </span>
      </div>
    </div>
  );
};

export default LogsFilter;
```

### 5. Logs Table Component

Create `src/components/Logs/LogsTable.jsx`:

```javascript
import React from 'react';
import { formatDate, formatTime } from '../../utils/dateUtils';
import './LogsTable.css';

const LogsTable = ({ logs, loading, totalCount }) => {
  const getActionBadge = (action) => {
    const actionTypes = {
      'CREATE': 'badge-success',
      'UPDATE': 'badge-warning',
      'DELETE': 'badge-danger',
      'LOGIN': 'badge-info',
      'LOGOUT': 'badge-secondary',
    };

    const actionType = action.split('_')[0];
    const badgeClass = actionTypes[actionType] || 'badge-primary';

    return (
      <span className={`action-badge ${badgeClass}`}>
        {action.replace('_', ' ')}
      </span>
    );
  };

  const getComponentIcon = (component) => {
    const icons = {
      'EMPLOYEE': '👤',
      'LEAVE': '📅',
      'PAYROLL': '💰',
      'SETTINGS': '⚙️',
      'AUTH': '🔐',
      'DESIGNATION': '🏷️',
    };
    return icons[component] || '📋';
  };

  if (totalCount === 0 && !loading) {
    return (
      <div className="no-logs">
        <div className="no-logs-icon">📋</div>
        <h3>No Logs Found</h3>
        <p>No system activities found for the selected time period.</p>
      </div>
    );
  }

  return (
    <div className="logs-table-container">
      <div className="table-wrapper">
        <table className="logs-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>User</th>
              <th>Action</th>
              <th>Component</th>
              <th>Date</th>
              <th>Time</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log) => (
              <tr key={log.id}>
                <td className="log-id">#{log.id}</td>
                <td className="user-name">
                  <div className="user-info">
                    <span className="user-avatar">👤</span>
                    {log.user_name}
                  </div>
                </td>
                <td className="action-cell">
                  {getActionBadge(log.action)}
                </td>
                <td className="component-cell">
                  <div className="component-info">
                    <span className="component-icon">
                      {getComponentIcon(log.component)}
                    </span>
                    {log.component}
                  </div>
                </td>
                <td className="date-cell">
                  {formatDate(log.created_at)}
                </td>
                <td className="time-cell">
                  {formatTime(log.created_at)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default LogsTable;
```

### 6. Utility Functions

Create `src/utils/dateUtils.js`:

```javascript
export const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

export const formatTime = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

export const formatDateTime = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};
```

### 7. CSS Styles

Create `src/components/Logs/LogsPage.css`:

```css
.logs-page {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.logs-header {
  margin-bottom: 32px;
}

.logs-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.logs-subtitle {
  color: #6b7280;
  font-size: 1rem;
}

/* Filter Styles */
.logs-filter {
  background: white;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.filter-section {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.filter-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.days-select,
.custom-days-input {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.875rem;
}

.custom-days-input {
  width: 120px;
}

.refresh-btn {
  padding: 8px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background-color 0.2s;
}

.refresh-btn:hover {
  background: #2563eb;
}

.filter-info {
  display: flex;
  gap: 24px;
  font-size: 0.875rem;
  color: #6b7280;
}

/* Table Styles */
.logs-table-container {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.table-wrapper {
  overflow-x: auto;
}

.logs-table {
  width: 100%;
  border-collapse: collapse;
}

.logs-table th {
  background: #f9fafb;
  padding: 12px 16px;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

.logs-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #f3f4f6;
}

.logs-table tr:hover {
  background: #f9fafb;
}

.log-id {
  font-family: monospace;
  color: #6b7280;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-avatar {
  font-size: 1.2rem;
}

.action-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
}

.badge-success { background: #dcfce7; color: #166534; }
.badge-warning { background: #fef3c7; color: #92400e; }
.badge-danger { background: #fee2e2; color: #991b1b; }
.badge-info { background: #dbeafe; color: #1e40af; }
.badge-secondary { background: #f3f4f6; color: #374151; }
.badge-primary { background: #e0e7ff; color: #3730a3; }

.component-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.component-icon {
  font-size: 1.1rem;
}

.date-cell,
.time-cell {
  font-family: monospace;
  font-size: 0.875rem;
  color: #6b7280;
}

.no-logs {
  text-align: center;
  padding: 48px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.no-logs-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.no-logs h3 {
  color: #374151;
  margin-bottom: 8px;
}

.no-logs p {
  color: #6b7280;
}

/* Responsive */
@media (max-width: 768px) {
  .logs-page {
    padding: 16px;
  }
  
  .filter-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-info {
    flex-direction: column;
    gap: 8px;
  }
}
```

### 8. Route Configuration

Add to your `src/App.js` or routing configuration:

```javascript
import { Routes, Route } from 'react-router-dom';
import LogsPage from './components/Logs/LogsPage';
import ProtectedRoute from './components/auth/ProtectedRoute';

function App() {
  return (
    <Routes>
      {/* Other routes */}
      <Route 
        path="/logs" 
        element={
          <ProtectedRoute requiredRole="SUPERADMIN">
            <LogsPage />
          </ProtectedRoute>
        } 
      />
    </Routes>
  );
}
```

### 9. Navigation Menu Item

Add to your navigation component:

```javascript
// Only show for SUPERADMIN users
{userRole === 'SUPERADMIN' && (
  <NavLink to="/logs" className="nav-item">
    📋 System Logs
  </NavLink>
)}
```

---

## Features Included

✅ **Complete API Integration**  
✅ **Role-based Access Control**  
✅ **Flexible Date Filtering**  
✅ **Real-time Data Refresh**  
✅ **Responsive Design**  
✅ **Error Handling**  
✅ **Loading States**  
✅ **Custom Hooks**  
✅ **Reusable Components**  
✅ **Professional UI/UX**

## Usage Examples

```javascript
// Basic usage
<LogsPage />

// With custom initial days
<LogsPage initialDays={30} />

// Using the hook directly
const { logs, loading, error, refetch } = useLogs(7);
```

This implementation provides a complete, production-ready logs management interface for your React application!