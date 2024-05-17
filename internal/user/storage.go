package user

import (
	"sort"
	"strconv"
	"sync"
	internalErrors "zpe-cloud-user-management-service/internal/msgs"
)

// User represents a user in the system.
var (
	mu        sync.Mutex
	users     = make(map[string]*User)
	idCounter int
)

// InitializeStorage sets up the in-memory storage for users.
func InitializeStorage() {
	mu.Lock()
	defer mu.Unlock()
	users = make(map[string]*User)
	idCounter = 0
}

// CreateUser adds a new user to the storage.
func CreateUser(user *User) error {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if u.Email == user.Email {
			return internalErrors.ErrUserAlreadyExists
		}
	}

	// Assign a new unique ID to the user and add them to the storage.
	idCounter++
	user.ID = strconv.Itoa(idCounter)
	users[user.ID] = user
	return nil
}

// GetUser retrieves a user from the storage by ID.
func GetUser(id string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[id]
	if !exists {
		return nil, internalErrors.ErrUserNotFound
	}
	return user, nil
}

// ListUsers returns a list of all users in the storage.
func ListUsers() ([]*User, error) {
	mu.Lock()
	defer mu.Unlock()

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	sort.Slice(userList, func(i, j int) bool {
		return userList[i].ID < userList[j].ID
	})

	return userList, nil
}

// UpdateUserRoles updates the roles of a user in the storage.
func UpdateUserRoles(id string, roles []string) error {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[id]
	if !exists {
		return internalErrors.ErrUserNotFound
	}

	user.Roles = roles
	return nil
}

// DeleteUser removes a user from the storage by ID.
func DeleteUser(id string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[id]; !exists {
		return internalErrors.ErrUserNotFound
	}

	delete(users, id)
	return nil
}
