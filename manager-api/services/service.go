package services

import (
	"fmt"
	"sync"

	"manager-api/models"
)

// Store holds the in-memory state for managers, employees, and leave requests
type Store struct {
	mu           sync.RWMutex
	managers     map[string]*models.Manager
	employees    map[string]*models.Employee
	leaveRequests map[string]*models.LeaveRequest
	nextID       int
}

// NewStore initializes the in-memory store with seed data
func NewStore() *Store {
	s := &Store{
		managers:      make(map[string]*models.Manager),
		employees:     make(map[string]*models.Employee),
		leaveRequests: make(map[string]*models.LeaveRequest),
		nextID:        1,
	}
	s.seed()
	return s
}

// seed populates the store with initial test data
func (s *Store) seed() {
	// Seed manager: mgr1 oversees emp1 and emp2
	s.managers["mgr1"] = &models.Manager{
		ID:            "mgr1",
		Name:          "Manager One",
		DirectReports: []string{"emp1", "emp2"},
	}

	// Seed employees with leave balance
	s.employees["emp1"] = &models.Employee{
		ID:      "emp1",
		Name:    "Employee One",
		Balance: 20,
	}
	s.employees["emp2"] = &models.Employee{
		ID:      "emp2",
		Name:    "Employee Two",
		Balance: 15,
	}

	// Seed pending leave requests
	s.leaveRequests["req1"] = &models.LeaveRequest{
		ID:         "req1",
		EmployeeID: "emp1",
		Type:       "annual",
		Days:       3,
		Status:     "pending",
	}
	s.leaveRequests["req2"] = &models.LeaveRequest{
		ID:         "req2",
		EmployeeID: "emp2",
		Type:       "sick",
		Days:       2,
		Status:     "pending",
	}
}

// GetPendingRequestsForManager returns all pending leave requests for a manager's direct reports
func (s *Store) GetPendingRequestsForManager(managerID string) ([]*models.LeaveRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mgr, ok := s.managers[managerID]
	if !ok {
		return nil, fmt.Errorf("manager %q not found", managerID)
	}

	// Build a set of direct reports for quick lookup
	directReportSet := make(map[string]bool)
	for _, empID := range mgr.DirectReports {
		directReportSet[empID] = true
	}

	var requests []*models.LeaveRequest
	for _, req := range s.leaveRequests {
		if req.Status == "pending" && directReportSet[req.EmployeeID] {
			requests = append(requests, req)
		}
	}
	return requests, nil
}

// ApproveRequest approves a leave request after validating manager authorization
func (s *Store) ApproveRequest(requestID, managerID, note string) (*models.LeaveRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, ok := s.leaveRequests[requestID]
	if !ok {
		return nil, fmt.Errorf("request %q not found", requestID)
	}

	if req.Status != "pending" {
		return nil, fmt.Errorf("request %q is not pending (current status: %s)", requestID, req.Status)
	}

	// Validate manager authorization
	if err := s.validateManagerAuthority(managerID, req.EmployeeID); err != nil {
		return nil, err
	}

	// Check employee has sufficient balance
	emp, ok := s.employees[req.EmployeeID]
	if !ok {
		return nil, fmt.Errorf("employee %q not found", req.EmployeeID)
	}
	if emp.Balance < req.Days {
		return nil, fmt.Errorf("employee %q has insufficient leave balance (%d days) for request of %d days", req.EmployeeID, emp.Balance, req.Days)
	}

	// Deduct balance and approve
	emp.Balance -= req.Days
	req.Status = "approved"
	req.ManagerNote = note

	return req, nil
}

// RejectRequest rejects a leave request after validating manager authorization
func (s *Store) RejectRequest(requestID, managerID, note string) (*models.LeaveRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, ok := s.leaveRequests[requestID]
	if !ok {
		return nil, fmt.Errorf("request %q not found", requestID)
	}

	if req.Status != "pending" {
		return nil, fmt.Errorf("request %q is not pending (current status: %s)", requestID, req.Status)
	}

	// Validate manager authorization
	if err := s.validateManagerAuthority(managerID, req.EmployeeID); err != nil {
		return nil, err
	}

	// Reject without balance deduction
	req.Status = "rejected"
	req.ManagerNote = note

	return req, nil
}

// validateManagerAuthority checks that the manager is authorized to act on the employee's requests
func (s *Store) validateManagerAuthority(managerID, employeeID string) error {
	mgr, ok := s.managers[managerID]
	if !ok {
		return fmt.Errorf("manager %q not found or not authorized", managerID)
	}

	for _, empID := range mgr.DirectReports {
		if empID == employeeID {
			return nil
		}
	}
	return fmt.Errorf("manager %q is not authorized to act on requests for employee %q", managerID, employeeID)
}
