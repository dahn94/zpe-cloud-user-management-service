package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	internalMsgs "zpe-cloud-user-management-service/internal/msgs"
)

var roleHierarchy = map[string][]string{
	"Admin":    {"Modifier", "Watcher"},
	"Modifier": {"Watcher"},
	"Watcher":  {},
}

func isRoleExists(role string) bool {
	_, exists := roleHierarchy[role]
	return exists
}

func isValidCrudOperation(currentUserRole, targetUserRole string) bool {
	if currentUserRole == "Admin" {
		return true
	}

	allowedRoles := roleHierarchy[currentUserRole]
	for _, role := range allowedRoles {
		if role == targetUserRole {
			return true
		}
	}
	return false
}

func checkPermission(w http.ResponseWriter, currentUserRole, requiredRole string) bool {
	if !isRoleExists(currentUserRole) || !isValidCrudOperation(currentUserRole, requiredRole) {
		errResponse(w, http.StatusForbidden, internalMsgs.ErrForbidden)
		log.Printf("Forbidden: UserType=%s does not have permission for role %s", currentUserRole, requiredRole)
		return false
	}
	return true
}

func isValidRoleUpdate(newRoles []string, sessionUserRole string) error {
	for _, role := range newRoles {
		if !isRoleExists(role) {
			return fmt.Errorf("%w: %s", internalMsgs.ErrInvalidRole, role)
		}
		if role == "Admin" && sessionUserRole != "Admin" {
			return fmt.Errorf("%w: %s", internalMsgs.ErrInsufficientPermissions, role)
		}
		if !isValidCrudOperation(sessionUserRole, role) {
			return fmt.Errorf("%w: %s", internalMsgs.ErrInsufficientPermissions, role)
		}
	}
	return nil
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func errResponse(w http.ResponseWriter, code int, err error) {
	jsonResponse(w, code, map[string]string{"message": err.Error()})
}

func getUserTypeByID(id string) (string, error) {
	user, err := GetUser(id)
	if err != nil {
		return "", internalMsgs.ErrUserNotFound
	}
	if len(user.Roles) > 0 {
		return user.Roles[0], nil
	}
	return "", errors.New("app has no roles")
}
