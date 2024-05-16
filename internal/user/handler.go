package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	internalMsgs "zpe-cloud-user-management-service/internal/msgs"
)

type RoleUpdateRequest struct {
	Roles []string `json:"roles"`
}

// HandleUsers handles HTTP requests for the /users endpoint.
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandleCreateUser(w, r)
	case http.MethodGet:
		HandleListUsers(w, r)
	default:
		errResponse(w, http.StatusMethodNotAllowed, internalMsgs.ErrMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}

// HandleUser handles HTTP requests for individual user operations at /users/{id}.
func HandleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandleGetUser(w, r)
	case http.MethodDelete:
		HandleDeleteUser(w, r)
	default:
		errResponse(w, http.StatusMethodNotAllowed, internalMsgs.ErrMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}

// HandleUserRoles handles HTTP requests for updating user roles at /users/roles/{id}.
func HandleUserRoles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		HandleUpdateUserRoles(w, r)
	default:
		errResponse(w, http.StatusMethodNotAllowed, internalMsgs.ErrMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	currentUserRole := r.Header.Get("X-User-Type")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errResponse(w, http.StatusBadRequest, internalMsgs.ErrInvalidRequestPayload)
		log.Printf("BadRequest: %v", err)
		return
	}

	if err := user.ValidateRequiredFields(); err != nil {
		errResponse(w, http.StatusBadRequest, err)
		log.Printf("BadRequest: %v", err)
		return
	}

	if !isValidCrudOperation(currentUserRole, user.Roles[0]) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrInsufficientPermissions)
		log.Printf("Forbidden: UserType=%s attempted to create a user", currentUserRole)
		return
	}
	if err := CreateUser(&user); err != nil {
		errResponse(w, http.StatusConflict, internalMsgs.ErrUserAlreadyExists)
		log.Printf("Conflict: %v", err)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]string{
		"id":      user.ID,
		"message": "User created successfully",
	})
	log.Printf("User created: %v", user)
}

func HandleListUsers(w http.ResponseWriter, r *http.Request) {
	currentUserRole := r.Header.Get("X-User-Type")
	if !isRoleExists(currentUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to list users", currentUserRole)
		return
	}
	users, err := ListUsers()
	if err != nil {
		errResponse(w, http.StatusInternalServerError, internalMsgs.ErrInternalServerError)
		log.Printf("InternalServerError: %v", err)
		return
	}

	jsonResponse(w, http.StatusOK, users)
	log.Printf("Users listed: %d users", len(users))
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	currentUserRole := r.Header.Get("X-User-Type")
	if !isRoleExists(currentUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to get a user", currentUserRole)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" {
		HandleListUsers(w, r)
		return
	}

	user, err := GetUser(id)
	if err != nil {
		users, err := ListUsers()
		if err != nil {
			jsonResponse(w, http.StatusInternalServerError, internalMsgs.ErrInternalServerError)
			log.Printf("InternalServerError: %v", err)
			return
		}

		if len(users) > 0 {
			jsonResponse(w, http.StatusOK, users)
			log.Printf("User not found. Returning list of all users")
		} else {
			jsonResponse(w, http.StatusNotFound, []*User{})
			log.Printf("User not found. No users in the system. Returning empty list")
		}
		return
	}

	jsonResponse(w, http.StatusOK, []*User{user})
	log.Printf("User retrieved: %v", *user)
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	currentUserRole := r.Header.Get("X-User-Type")
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	targetUserRole, err := getUserTypeByID(id)
	if err != nil {
		errResponse(w, http.StatusNotFound, internalMsgs.ErrUserNotFound)
		log.Printf("User not found: %s", id)
		return
	}

	if !checkPermission(w, currentUserRole, targetUserRole) {
		return
	}

	if err := DeleteUser(id); err != nil {
		errResponse(w, http.StatusNotFound, internalMsgs.ErrUserNotFound)
		log.Printf("NotFound: User %s", id)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("UserType=%s deleted user %s", currentUserRole, id)
}

func HandleUpdateUserRoles(w http.ResponseWriter, r *http.Request) {
	currentUserRole := r.Header.Get("X-User-Type")
	id := strings.TrimPrefix(r.URL.Path, "/users/roles/")

	var req RoleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errResponse(w, http.StatusBadRequest, internalMsgs.ErrInvalidRequestPayload)
		log.Printf("BadRequest: %v", err)
		return
	}

	if err := isValidRoleUpdate(req.Roles, currentUserRole); err != nil {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrInsufficientPermissions)
		log.Printf("Forbidden: %v", err)
		return
	}

	if err := UpdateUserRoles(id, req.Roles); err != nil {
		errResponse(w, http.StatusNotFound, internalMsgs.ErrUserNotFound)
		log.Printf("NotFound: %v", err)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "User roles updated successfully"})
	log.Printf("User roles updated: %s", id)
}
