package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/models"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/internal/service"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/pkg/response"
	"github.com/johnmerga/realgaming-marketplace-backend/marketplace-backend/pkg/validator"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validator
	logger      zerolog.Logger
}

func NewUserHandler(userService service.UserService, validator *validator.Validator, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
		logger:      logger,
	}
}

// CreateUser creates a new user
// POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	// Validate and parse JSON
	if err := h.validator.ValidateAndParseJSON(r, &req); err != nil {
		h.logger.Error().Err(err).Msg("validation failed")
		response.JSON(w, http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	// Create user
	user, err := h.userService.CreateUser(r.Context(), &req)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create user")
		if strings.Contains(err.Error(), "already exists") {
			response.JSON(w, http.StatusConflict, response.Error(err.Error()))
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	h.logger.Info().Str("user_id", user.ID.String()).Msg("user created successfully")
	response.JSON(w, http.StatusCreated, response.SuccessWithMessage(user, "User created successfully"))
}

// GetUser gets a user by ID
// GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.Error("invalid user ID"))
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", id.String()).Msg("failed to get user")
		if strings.Contains(err.Error(), "not found") {
			response.JSON(w, http.StatusNotFound, response.Error("User not found"))
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	response.JSON(w, http.StatusOK, response.Success(user))
}

// UpdateUser updates a user
// PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.Error("invalid user ID"))
		return
	}

	var req models.UpdateUserRequest

	// Validate and parse JSON
	if err := h.validator.ValidateAndParseJSON(r, &req); err != nil {
		h.logger.Error().Err(err).Msg("validation failed")
		response.JSON(w, http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", id.String()).Msg("failed to update user")
		if strings.Contains(err.Error(), "not found") {
			response.JSON(w, http.StatusNotFound, response.Error("User not found"))
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	h.logger.Info().Str("user_id", user.ID.String()).Msg("user updated successfully")
	response.JSON(w, http.StatusOK, response.SuccessWithMessage(user, "User updated successfully"))
}

// DeleteUser (actually updates status to inactive)
// DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.Error("invalid user ID"))
		return
	}

	err = h.userService.UpdateUserStatus(r.Context(), id, models.StatusInactive)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", id.String()).Msg("failed to delete user")
		if strings.Contains(err.Error(), "not found") {
			response.JSON(w, http.StatusNotFound, response.Error("User not found"))
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	h.logger.Info().Str("user_id", id.String()).Msg("user deleted successfully")
	response.JSON(w, http.StatusOK, response.SuccessWithMessage(nil, "User deleted successfully"))
}

// ListUsers lists users with optional filters
// GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Pagination
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Filters
	var role *models.UserRole
	var status *models.UserStatus

	if roleStr := query.Get("role"); roleStr != "" {
		roleVal := models.UserRole(roleStr)
		role = &roleVal
	}

	if statusStr := query.Get("status"); statusStr != "" {
		statusVal := models.UserStatus(statusStr)
		status = &statusVal
	}

	users, err := h.userService.ListUsers(r.Context(), role, status, page, limit)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list users")
		response.JSON(w, http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	// For pagination, you might want to get total count as well
	// This is simplified - you'd typically need a separate count query
	total := len(users)

	response.JSON(w, http.StatusOK, response.Paginated(users, page, limit, total))
}

// Login authenticates a user
// POST /api/v1/auth/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	// Validate and parse JSON
	if err := h.validator.ValidateAndParseJSON(r, &req); err != nil {
		h.logger.Error().Err(err).Msg("validation failed")
		response.JSON(w, http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	user, err := h.userService.Login(r.Context(), &req)
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("login failed")
		if strings.Contains(err.Error(), "credentials") || strings.Contains(err.Error(), "inactive") {
			response.JSON(w, http.StatusUnauthorized, response.Error(err.Error()))
			return
		}
		response.JSON(w, http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	h.logger.Info().Str("user_id", user.ID.String()).Msg("user logged in successfully")
	response.JSON(w, http.StatusOK, response.SuccessWithMessage(user, "Login successful"))
}
