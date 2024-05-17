package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	internalMsgs "zpe-cloud-user-management-service/internal/msgs"
)

// RoleUpdateRequest represents the request payload for updating a user's roles.
type RoleUpdateRequest struct {
	Roles []string `json:"roles"`
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

	if !isRoleExists(currentUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to create users", currentUserRole)
		return
	}

	if !isValidCrudOperation(currentUserRole, user.Roles...) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
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
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	if !isRoleExists(currentUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to get a user", currentUserRole)
		return
	}

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
		log.Printf("NotFound: %v", err)
		return
	}

	if !isRoleExists(currentUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to delete a user", currentUserRole)
		return
	}

	if !isValidCrudOperation(currentUserRole, targetUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to delete a user with role %s", currentUserRole, targetUserRole)
		return
	}

	if err := DeleteUser(id); err != nil {
		errResponse(w, http.StatusInternalServerError, err)
		log.Printf("InternalServerError: %v", err)
		return
	}

	jsonResponse(w, http.StatusNoContent, nil)
	log.Printf("User deleted: %s", id)
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

	if !isRoleExists(currentUserRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s attempted to update roles of a user", currentUserRole)
		return
	}

	if !isValidCrudOperation(currentUserRole, req.Roles...) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: %v", currentUserRole)
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
