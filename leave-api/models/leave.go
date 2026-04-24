package models

import "time"

// LeaveType represents the category of leave
type LeaveType string

const (
	LeaveTypeAnnual LeaveType = "annual"
	LeaveTypeSick   LeaveType = "sick"
)

// LeaveStatus represents the status of a leave request
type LeaveStatus string

const (
	LeaveStatusPending  LeaveStatus = "pending"
	LeaveStatusApproved LeaveStatus = "approved"
	LeaveStatusRejected LeaveStatus = "rejected"
)

// LeaveRequest represents an employee leave request
type LeaveRequest struct {
	ID         string      `json:"id"`
	EmployeeID string      `json:"employeeId"`
	LeaveType  LeaveType   `json:"leaveType"`
	StartDate  string      `json:"startDate"` // YYYY-MM-DD
	EndDate    string      `json:"endDate"`   // YYYY-MM-DD
	Reason     string      `json:"reason"`
	Status     LeaveStatus `json:"status"`
	CreatedAt  time.Time   `json:"createdAt"`
}

// LeaveBalance represents an employee's leave balance per category
type LeaveBalance struct {
	EmployeeID string             `json:"employeeId"`
	Balances   map[LeaveType]int  `json:"balances"`
}

// CreateLeaveRequestInput is the payload for creating a leave request
type CreateLeaveRequestInput struct {
	EmployeeID string    `json:"employeeId"`
	LeaveType  LeaveType `json:"leaveType"`
	StartDate  string    `json:"startDate"`
	EndDate    string    `json:"endDate"`
	Reason     string    `json:"reason"`
}

// ErrorResponse is a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}
