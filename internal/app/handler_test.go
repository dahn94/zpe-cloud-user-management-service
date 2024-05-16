package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func setupTestStorageWithUsers() {
	InitializeStorage()
	users := []*User{
		{Name: "Leia Organa", Email: "leia@example.com", Roles: []string{"Admin"}},
		{Name: "Obi-Wan Kenobi", Email: "obi-wan@example.com", Roles: []string{"Modifier"}},
		{Name: "R2-D2", Email: "r2-d2@example.com", Roles: []string{"Watcher"}},
		{Name: "Vegeta", Email: "vegeta@example.com", Roles: []string{"Modifier"}},
		{Name: "Gohan", Email: "gohan@example.com", Roles: []string{"Watcher"}},
		{Name: "Goku", Email: "goku@example", Roles: []string{"Admin"}},
	}

	for i, user := range users {
		if err := CreateUser(user); err != nil {
			log.Fatalf("Failed to create app%d: %v", i+1, err)
		}
	}
}

func TestHandleCreateUser(t *testing.T) {

	tests := []struct {
		name           string
		userType       string
		payload        User
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "Admin can create app at same level",
			userType: "Admin",
			payload: User{
				Name:  "Jar Jar Binks",
				Email: "Binks@example.com",
				Roles: []string{"Admin"},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"1","message":"User created successfully"}`,
		},
		{
			name:     "Admin can create app at modifier level",
			userType: "Admin",
			payload: User{
				Name:  "Darth Vader",
				Email: "vader@example.com",
				Roles: []string{"Modifier"},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"2","message":"User created successfully"}`,
		},
		{
			name:     "Admin can create app at watcher level",
			userType: "Admin",
			payload: User{
				Name:  "Yoda",
				Email: "yoda@example.com",
				Roles: []string{"Watcher"},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"3","message":"User created successfully"}`,
		},
		{
			name:     "Admin can create app with more than one role",
			userType: "Admin",
			payload: User{
				Name:  "Padmé Amidala",
				Email: "amidala@example.com",
				Roles: []string{"Modifier", "Watcher"},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"4","message":"User created successfully"}`,
		},
		{
			name:     "User creation fails if app already exists",
			userType: "Admin",
			payload: User{
				Name:  "Yoda",
				Email: "yoda@example.com",
				Roles: []string{"Watcher"},
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"message":"user already exists"}`,
		},
		{
			name:     "Modifier cannot create app at same level",
			userType: "Modifier",
			payload: User{
				Name:  "Chewbacca",
				Email: "chewbacca@example.com",
				Roles: []string{"Modifier"},
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"message":"insufficient permissions to assign role"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonPayload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-User-Type", tt.userType)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleCreateUser)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var actualBody, expectedBody map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &actualBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody); err != nil {
				t.Fatalf("Failed to unmarshal expected body: %v", err)
			}

			if !reflect.DeepEqual(actualBody, expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandleListUsers(t *testing.T) {
	setupTestStorageWithUsers()

	tests := []struct {
		name           string
		userType       string
		expectedStatus int
	}{
		{
			name:           "Admin can list users",
			userType:       "Admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Watcher can list users",
			userType:       "Watcher",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unknown app cannot list users",
			userType:       "Unknown",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/users", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-User-Type", tt.userType)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleListUsers)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandleGetUser(t *testing.T) {
	setupTestStorageWithUsers()

	tests := []struct {
		name           string
		userType       string
		userID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Admin can get app",
			userType:       "Admin",
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1","name":"Leia Organa","email":"leia@example.com","roles":["Admin"]}]`,
		},
		{
			name:           "Unknown cannot get app",
			userType:       "Unknown",
			userID:         "1",
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"message":"forbidden"}`,
		},
		{
			name:           "Get non-existent app",
			userType:       "Admin",
			userID:         "999",
			expectedStatus: http.StatusOK,
			expectedBody: `[{"id":"1","name":"Leia Organa","email":"leia@example.com","roles":["Admin"]},
                            {"id":"2","name":"Obi-Wan Kenobi","email":"obi-wan@example.com","roles":["Modifier"]},
                            {"id":"3","name":"R2-D2","email":"r2-d2@example.com","roles":["Watcher"]},
							{"id":"4","name":"Vegeta","email":"vegeta@example.com","roles":["Modifier"]},
							{"id":"5","name":"Gohan","email":"gohan@example.com","roles":["Watcher"]},
							{"id":"6","name":"Goku","email":"goku@example","roles":["Admin"]}]`,
		},
		{
			name:           "Get non-existent app with empty list",
			userType:       "Admin",
			userID:         "999",
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Get non-existent app with empty list" {
				InitializeStorage() // Clear storage for this test
			}

			req, err := http.NewRequest("GET", "/users/"+tt.userID, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-User-Type", tt.userType)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleGetUser)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var actualBody, expectedBody interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &actualBody); err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody); err != nil {
				t.Fatalf("Failed to unmarshal expected body: %v", err)
			}

			if !reflect.DeepEqual(actualBody, expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandleDeleteUser(t *testing.T) {
	setupTestStorageWithUsers()

	tests := []struct {
		name           string
		userType       string
		userID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Admin can delete modifier app",
			userType:       "Admin",
			userID:         "1",
			expectedStatus: http.StatusNoContent,
			expectedBody:   `{"message":"User deleted successfully"}`,
		},
		{
			name:           "Admin can delete watcher app",
			userType:       "Admin",
			userID:         "2",
			expectedStatus: http.StatusNoContent,
			expectedBody:   `{"message":"User deleted successfully"}`,
		},
		{
			name:           "Modifier can delete watcher app",
			userType:       "Modifier",
			userID:         "3",
			expectedStatus: http.StatusNoContent,
			expectedBody:   `{"message":"User deleted successfully"}`,
		},
		{
			name:           "Modifier cannot delete admin app",
			userType:       "Modifier",
			userID:         "4",
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"message":"forbidden"}`,
		},
		{
			name:           "Watcher cannot delete any app",
			userType:       "Watcher",
			userID:         "5", // ID of Goku (Admin)
			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"message":"forbidden"}`,
		},
		{
			name:           "Delete non-existent app",
			userType:       "Admin",
			userID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"user not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/users/"+tt.userID, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-User-Type", tt.userType)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleDeleteUser)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus != http.StatusNoContent {
				var actualBody map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&actualBody); err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				var expectedBody map[string]string
				if err := json.Unmarshal([]byte(tt.expectedBody), &expectedBody); err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}

				if !reflect.DeepEqual(actualBody, expectedBody) {
					t.Errorf("handler returned unexpected body: got %v want %v", actualBody, expectedBody)
				}
			}
		})
	}
}

func TestHandleUpdateUserRoles(t *testing.T) {
	setupTestStorageWithUsers()

	tests := []struct {
		name           string
		userType       string
		userID         string
		payload        RoleUpdateRequest
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:     "Admin can update roles",
			userType: "Admin",
			userID:   "1",
			payload: RoleUpdateRequest{
				Roles: []string{"Modifier"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "User roles updated successfully"},
		},
		{
			name:     "Modifier cannot update roles to Admin",
			userType: "Modifier",
			userID:   "1",
			payload: RoleUpdateRequest{
				Roles: []string{"Admin"},
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   map[string]string{"message": "insufficient permissions to assign role"},
		},
		{
			name:     "Update roles for non-existent user",
			userType: "Admin",
			userID:   "999",
			payload: RoleUpdateRequest{
				Roles: []string{"Modifier"},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"message": "user not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req, err := http.NewRequest("PUT", "/users/roles/"+tt.userID, bytes.NewBuffer(jsonPayload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-User-Type", tt.userType)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(HandleUpdateUserRoles)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var response map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if !reflect.DeepEqual(response, tt.expectedBody) {
				t.Errorf("handler returned unexpected message: got %v want %v", response, tt.expectedBody)
			}
		})
	}
}
