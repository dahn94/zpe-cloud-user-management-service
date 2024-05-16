package app

import "errors"

// User represents a user in the system with ID, Name, Email, and Roles.
type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// Validate checks that the user has all necessary fields filled out.
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if len(u.Roles) == 0 {
		return errors.New("roles are required")
	}
	return nil
}
