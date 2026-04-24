package services

import (
	"fmt"
	"leave-api/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

const dateLayout = "2006-01-02"

// LeaveService manages leave requests and balances using in-memory storage
type LeaveService struct {
	mu       sync.RWMutex
	requests map[string]*models.LeaveRequest   // requestID -> request
	balances map[string]map[models.LeaveType]int // employeeID -> leaveType -> days
}

// NewLeaveService creates and seeds a new LeaveService
func NewLeaveService() *LeaveService {
	svc := &LeaveService{
		requests: make(map[string]*models.LeaveRequest),
		balances: make(map[string]map[models.LeaveType]int),
	}

	// Seed initial balances
	svc.balances["emp1"] = map[models.LeaveType]int{
		models.LeaveTypeAnnual: 15,
		models.LeaveTypeSick:   10,
	}
	svc.balances["emp2"] = map[models.LeaveType]int{
		models.LeaveTypeAnnual: 20,
		models.LeaveTypeSick:   5,
	}

	return svc
}

// GetBalance returns the leave balance for an employee
func (s *LeaveService) GetBalance(employeeID string) (*models.LeaveBalance, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bal, ok := s.balances[employeeID]
	if !ok {
		// Return zero balances for unknown employees
		bal = map[models.LeaveType]int{
			models.LeaveTypeAnnual: 0,
			models.LeaveTypeSick:   0,
		}
	}

	// Return a copy
	balCopy := make(map[models.LeaveType]int)
	for k, v := range bal {
		balCopy[k] = v
	}

	return &models.LeaveBalance{
		EmployeeID: employeeID,
		Balances:   balCopy,
	}, nil
}

// GetRequests returns all leave requests for an employee
func (s *LeaveService) GetRequests(employeeID string) ([]*models.LeaveRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*models.LeaveRequest
	for _, req := range s.requests {
		if req.EmployeeID == employeeID {
			// Return a copy
			copy := *req
			result = append(result, &copy)
		}
	}
	if result == nil {
		result = []*models.LeaveRequest{}
	}
	return result, nil
}

// CreateRequest creates a new leave request with overlap and validation checks
func (s *LeaveService) CreateRequest(input models.CreateLeaveRequestInput) (*models.LeaveRequest, error) {
	// Validate required fields
	if input.EmployeeID == "" {
		return nil, &ValidationError{Message: "employeeId is required"}
	}
	if input.LeaveType == "" {
		return nil, &ValidationError{Message: "leaveType is required"}
	}
	if input.LeaveType != models.LeaveTypeAnnual && input.LeaveType != models.LeaveTypeSick {
		return nil, &ValidationError{Message: "leaveType must be 'annual' or 'sick'"}
	}
	if input.StartDate == "" {
		return nil, &ValidationError{Message: "startDate is required"}
	}
	if input.EndDate == "" {
		return nil, &ValidationError{Message: "endDate is required"}
	}

	// Parse and validate dates
	start, err := time.Parse(dateLayout, input.StartDate)
	if err != nil {
		return nil, &ValidationError{Message: "startDate must be in YYYY-MM-DD format"}
	}
	end, err := time.Parse(dateLayout, input.EndDate)
	if err != nil {
		return nil, &ValidationError{Message: "endDate must be in YYYY-MM-DD format"}
	}
	if !end.Equal(start) && end.Before(start) {
		return nil, &ValidationError{Message: "endDate must be on or after startDate"}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for overlapping requests (only against pending/approved)
	for _, req := range s.requests {
		if req.EmployeeID != input.EmployeeID {
			continue
		}
		if req.Status == models.LeaveStatusRejected {
			continue
		}
		existStart, _ := time.Parse(dateLayout, req.StartDate)
		existEnd, _ := time.Parse(dateLayout, req.EndDate)

		if datesOverlap(start, end, existStart, existEnd) {
			return nil, &ConflictError{
				Message: fmt.Sprintf("leave request overlaps with existing request %s (%s to %s)", req.ID, req.StartDate, req.EndDate),
			}
		}
	}

	newReq := &models.LeaveRequest{
		ID:         uuid.New().String(),
		EmployeeID: input.EmployeeID,
		LeaveType:  input.LeaveType,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		Reason:     input.Reason,
		Status:     models.LeaveStatusPending,
		CreatedAt:  time.Now().UTC(),
	}

	s.requests[newReq.ID] = newReq

	// Return a copy
	copy := *newReq
	return &copy, nil
}

// CancelRequest cancels a pending leave request
func (s *LeaveService) CancelRequest(requestID, employeeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	req, ok := s.requests[requestID]
	if !ok {
		return &NotFoundError{Message: fmt.Sprintf("leave request %s not found", requestID)}
	}
	if req.EmployeeID != employeeID {
		return &NotFoundError{Message: fmt.Sprintf("leave request %s not found for employee %s", requestID, employeeID)}
	}
	if req.Status != models.LeaveStatusPending {
		return &ValidationError{Message: fmt.Sprintf("only pending requests can be cancelled; current status: %s", req.Status)}
	}

	delete(s.requests, requestID)
	return nil
}

// datesOverlap returns true if [s1,e1] and [s2,e2] overlap (inclusive)
func datesOverlap(s1, e1, s2, e2 time.Time) bool {
	return !s1.After(e2) && !s2.After(e1)
}

// --- Error types ---

// ValidationError is returned for bad input (HTTP 400)
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// ConflictError is returned for overlapping leave (HTTP 409)
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string { return e.Message }

// NotFoundError is returned when a resource doesn't exist (HTTP 404)
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string { return e.Message }
