package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"manager-api/models"
	"manager-api/services"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	store *services.Store
}

// NewHandler creates a new Handler with the given store
func NewHandler(store *services.Store) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes registers all HTTP routes on the given router
func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/manager/requests", h.GetPendingRequests).Methods(http.MethodGet)
	r.HandleFunc("/manager/requests/{id}/approve", h.ApproveRequest).Methods(http.MethodPost)
	r.HandleFunc("/manager/requests/{id}/reject", h.RejectRequest).Methods(http.MethodPost)
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GetPendingRequests handles GET /manager/requests?managerId=X
func (h *Handler) GetPendingRequests(w http.ResponseWriter, r *http.Request) {
	managerID := r.URL.Query().Get("managerId")
	if managerID == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "managerId query parameter is required"})
		return
	}

	requests, err := h.store.GetPendingRequestsForManager(managerID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Return empty array rather than null
	if requests == nil {
		requests = []*models.LeaveRequest{}
	}
	writeJSON(w, http.StatusOK, requests)
}

// ApproveRequest handles POST /manager/requests/{id}/approve
func (h *Handler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	var body models.ApproveRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if body.ManagerID == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "managerId is required"})
		return
	}

	req, err := h.store.ApproveRequest(requestID, body.ManagerID, body.Note)
	if err != nil {
		statusCode := errorStatusCode(err.Error())
		writeJSON(w, statusCode, models.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, req)
}

// RejectRequest handles POST /manager/requests/{id}/reject
func (h *Handler) RejectRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	var body models.RejectRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if body.ManagerID == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "managerId is required"})
		return
	}

	req, err := h.store.RejectRequest(requestID, body.ManagerID, body.Note)
	if err != nil {
		statusCode := errorStatusCode(err.Error())
		writeJSON(w, statusCode, models.ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, req)
}

// writeJSON writes a JSON response with the given status code and body
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// errorStatusCode maps error messages to appropriate HTTP status codes
func errorStatusCode(msg string) int {
	// "not found" errors → 404
	if contains(msg, "not found") {
		return http.StatusNotFound
	}
	// Authorization errors → 403
	if contains(msg, "not authorized") {
		return http.StatusForbidden
	}
	// All others → 400
	return http.StatusBadRequest
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
