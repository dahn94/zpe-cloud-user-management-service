package app

import (
	"errors"
	"strings"
)

// User represents a user in the system with ID, Name, Email, and Roles.
type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// ValidateRequiredFields checks that the user has all necessary fields filled out.
func (u *User) ValidateRequiredFields() error {
	var missingFields []string

	if u.Name == "" {
		missingFields = append(missingFields, "name")
	}
	if u.Email == "" {
		missingFields = append(missingFields, "email")
	}
	if len(u.Roles) == 0 {
		missingFields = append(missingFields, "roles")
	}

	if len(missingFields) > 0 {
		return errors.New("fields required: " + strings.Join(missingFields, ", "))
	}
	return nil
}
