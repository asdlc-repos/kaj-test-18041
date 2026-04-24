package handlers

import (
	"encoding/json"
	"errors"
	"leave-api/models"
	"leave-api/services"
	"net/http"
	"strings"
)

// LeaveHandler holds the leave service
type LeaveHandler struct {
	svc *services.LeaveService
}

// NewLeaveHandler creates a new LeaveHandler
func NewLeaveHandler(svc *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{svc: svc}
}

// RegisterRoutes registers all leave-related routes on the mux
func (h *LeaveHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/leave/balance", h.GetBalance)
	mux.HandleFunc("/leave/requests", h.LeaveRequests)
	mux.HandleFunc("/leave/requests/", h.DeleteRequest)
}

// Health handles GET /health
func (h *LeaveHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetBalance handles GET /leave/balance?employeeId=X
func (h *LeaveHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	employeeID := r.URL.Query().Get("employeeId")
	if employeeID == "" {
		writeError(w, http.StatusBadRequest, "employeeId query parameter is required")
		return
	}

	balance, err := h.svc.GetBalance(employeeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, balance)
}

// LeaveRequests handles GET and POST /leave/requests
func (h *LeaveHandler) LeaveRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getRequests(w, r)
	case http.MethodPost:
		h.createRequest(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// getRequests handles GET /leave/requests?employeeId=X
func (h *LeaveHandler) getRequests(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("employeeId")
	if employeeID == "" {
		writeError(w, http.StatusBadRequest, "employeeId query parameter is required")
		return
	}

	requests, err := h.svc.GetRequests(employeeID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, requests)
}

// createRequest handles POST /leave/requests
func (h *LeaveHandler) createRequest(w http.ResponseWriter, r *http.Request) {
	var input models.CreateLeaveRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req, err := h.svc.CreateRequest(input)
	if err != nil {
		var valErr *services.ValidationError
		var conflictErr *services.ConflictError
		if errors.As(err, &valErr) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else if errors.As(err, &conflictErr) {
			writeError(w, http.StatusConflict, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, req)
}

// DeleteRequest handles DELETE /leave/requests/{id}?employeeId=X
func (h *LeaveHandler) DeleteRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract ID from path: /leave/requests/{id}
	path := strings.TrimPrefix(r.URL.Path, "/leave/requests/")
	requestID := strings.TrimSuffix(path, "/")
	if requestID == "" {
		writeError(w, http.StatusBadRequest, "request ID is required in path")
		return
	}

	employeeID := r.URL.Query().Get("employeeId")
	if employeeID == "" {
		writeError(w, http.StatusBadRequest, "employeeId query parameter is required")
		return
	}

	if err := h.svc.CancelRequest(requestID, employeeID); err != nil {
		var notFoundErr *services.NotFoundError
		var valErr *services.ValidationError
		if errors.As(err, &notFoundErr) {
			writeError(w, http.StatusNotFound, err.Error())
		} else if errors.As(err, &valErr) {
			writeError(w, http.StatusBadRequest, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}
