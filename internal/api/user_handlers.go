package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/NetMob-Mobile-Networks-Project/polaris-server/internal/storage/db"
)

// Request/Response structures
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

// convertRole converts string role to db.UsersRole
func convertRole(role string) db.UsersRole {
	if role == "admin" {
		return db.UsersRoleAdmin
	}
	return db.UsersRoleUser
}

// ListUsersHandler returns a paginated list of users
func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// Get users with pagination
	params := db.ListUsersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	}
	users, err := h.db.ListUsers(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Get total count
	count, err := h.db.CountUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to count users", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:        int(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Time.Format(time.RFC3339),
		}
	}

	// Return response
	json.NewEncoder(w).Encode(ListUsersResponse{
		Users:      userResponses,
		TotalCount: count,
		Page:       page,
		PageSize:   pageSize,
	})
}

// CreateUserHandler creates a new user
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Validate role
	role := convertRole(req.Role)

	// Check if user already exists
	_, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user
	params := db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	err = h.db.CreateUser(r.Context(), params)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Get the created user
	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Failed to fetch created user", http.StatusInternalServerError)
		return
	}

	// Return created user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:        int(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Time.Format(time.RFC3339),
	})
}

// DeleteUserHandler deletes a user
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user exists
	_, err = h.db.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Prevent self-deletion
	currentUserID := r.Context().Value("user_id").(int)
	if currentUserID == userID {
		http.Error(w, "Cannot delete your own account", http.StatusForbidden)
		return
	}

	// Delete user
	err = h.db.DeleteUser(r.Context(), int32(userID))
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
