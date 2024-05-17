package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	internalMsgs "zpe-cloud-user-management-service/internal/msgs"
)

// roleHierarchy defines the hierarchy of roles and their allowed subordinate roles.
var roleHierarchy = map[string][]string{
	"Admin":    {"Admin", "Modifier", "Watcher"},
	"Modifier": {"Watcher"},
	"Watcher":  {},
}

// isRoleExists checks if a role exists in the role hierarchy.
func isRoleExists(role string) bool {
	_, exists := roleHierarchy[role]
	return exists
}

// isValidCrudOperation checks if the current user has permission to perform actions on the target user's role.
func isValidCrudOperation(sessionUserRole string, targetRoles ...string) bool {
	for _, targetRole := range targetRoles {
		allowedRoles := roleHierarchy[sessionUserRole]
		isAllowed := false
		for _, role := range allowedRoles {
			if role == targetRole {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return false
		}
	}
	return true
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
	return "", errors.New("user has no roles")
}
