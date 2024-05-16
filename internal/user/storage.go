package user

import (
	"sort"
	"strconv"
	"sync"
	internalErrors "zpe-cloud-user-management-service/internal/msgs"
)

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

func GetUser(id string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := users[id]
	if !exists {
		return nil, internalErrors.ErrUserNotFound
	}
	return user, nil
}

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

func DeleteUser(id string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[id]; !exists {
		return internalErrors.ErrUserNotFound
	}

	delete(users, id)
	return nil
}
