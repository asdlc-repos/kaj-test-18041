package models

// Manager represents a manager in the organizational hierarchy
type Manager struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	DirectReports []string `json:"directReports"`
}

// Employee represents an employee with leave balance
type Employee struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

// LeaveRequest represents a leave request from an employee
type LeaveRequest struct {
	ID          string `json:"id"`
	EmployeeID  string `json:"employeeId"`
	Type        string `json:"type"`
	Days        int    `json:"days"`
	Status      string `json:"status"`
	ManagerNote string `json:"managerNote,omitempty"`
}

// ApproveRequest is the body for approving a leave request
type ApproveRequest struct {
	ManagerID string `json:"managerId"`
	Note      string `json:"note,omitempty"`
}

// RejectRequest is the body for rejecting a leave request
type RejectRequest struct {
	ManagerID string `json:"managerId"`
	Note      string `json:"note,omitempty"`
}

// ErrorResponse represents an error response body
type ErrorResponse struct {
	Error string `json:"error"`
}
