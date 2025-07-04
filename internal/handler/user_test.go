package handler

import (
	"errors"
	"fmt"
	"go-crud/internal/domain"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

var fixedTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantErr    error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			method:     http.MethodPost,
			body:       `{"username": "test", "email": "test@email.net", "password": "mypw123"}`,
			wantErr:    nil,
			wantStatus: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"id":1,"email":"test@email.net","username":"test","created_at":"%s","updated_at":"%s"}`,
				fixedTime.Format(time.RFC3339), fixedTime.Format(time.RFC3339)),
		},
		{
			name:       "invalid method",
			method:     http.MethodGet,
			body:       "",
			wantErr:    nil,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed\n",
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       `{bad json}`,
			wantErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid JSON\n",
		},
		{
			name:       "repository error",
			method:     http.MethodPost,
			body:       `{"username": "test", "email": "test@email.net"}`,
			wantErr:    errors.New("repo error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Failed to create user: repo error\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				createFunc: func(*domain.User) error {
					return test.wantErr
				},
			}

			handler := NewUserHandler(repo)

			req := httptest.NewRequest(test.method, "/users", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			handler.Create(w, req)

			resp := w.Result()

			defer resp.Body.Close()

			if test.wantStatus != resp.StatusCode {
				t.Errorf("expected status code: %v, got: %v", test.wantStatus, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			respBodyStr := string(respBody)

			if strings.TrimSpace(string(respBody)) != strings.TrimSpace(test.wantBody) {
				t.Errorf("expected response body: %s, got: %s", test.wantBody, respBodyStr)
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		repoUser   *domain.User
		wantErr    error
		wantStatus int
		wantBody   string
	}{
		{
			name:   "success",
			method: http.MethodGet,
			path:   "/users/1",
			repoUser: &domain.User{
				ID:        1,
				Email:     "test@email.com",
				Username:  "testuser",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			wantErr:    nil,
			wantStatus: http.StatusOK,
			wantBody: fmt.Sprintf(`{"id":1,"email":"test@email.com","username":"testuser","created_at":"%s","updated_at":"%s"}`,
				fixedTime.Format(time.RFC3339), fixedTime.Format(time.RFC3339)),
		},
		{
			name:       "user not found",
			method:     http.MethodGet,
			path:       "/users/9999",
			repoUser:   nil,
			wantErr:    fmt.Errorf("user not found"),
			wantStatus: http.StatusNotFound,
			wantBody:   "User not found\n",
		},
		{
			name:       "method not allowed",
			method:     http.MethodPost,
			path:       "/users/1",
			repoUser:   nil,
			wantErr:    nil,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed\n",
		},
		{
			name:       "invalid path",
			method:     http.MethodGet,
			path:       "/invalid/path/1",
			repoUser:   nil,
			wantErr:    fmt.Errorf("invalid path format"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid path format\n",
		},
		{
			name:       "invalid id",
			method:     http.MethodGet,
			path:       "/users/1b",
			repoUser:   nil,
			wantErr:    fmt.Errorf("Invalid parameter 'id'\n"),
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid parameter 'id'\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				getByIDFunc: func(id int64) (*domain.User, error) {
					return test.repoUser, test.wantErr
				},
			}
			handler := NewUserHandler(repo)

			req := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()
			handler.GetByID(w, req)

			resp := w.Result()

			defer resp.Body.Close()

			if test.wantStatus != resp.StatusCode {
				t.Errorf("expected status code: %v, got: %v", test.wantStatus, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			respBodyStr := string(respBody)

			if strings.TrimSpace(string(respBody)) != strings.TrimSpace(test.wantBody) {
				t.Errorf("expected response body: %s, got: %s", test.wantBody, respBodyStr)
			}
		})
	}
}

func strPtr(s string) *string { return &s }

func TestUserHandler_Update(t *testing.T) {
	updatedTime := fixedTime.Add(time.Hour)
	tests := []struct {
		name       string
		method     string
		path       string
		repoUser   *domain.User
		updateErr  error
		body       string
		wantStatus int
		wantBody   string
		wantUpdate *domain.UserUpdate
	}{
		{
			name:   "update email success",
			method: http.MethodPut,
			path:   "/users/1",
			repoUser: &domain.User{
				ID:        1,
				Username:  "olduser",
				Email:     "new@email.com",
				CreatedAt: fixedTime,
				UpdatedAt: updatedTime,
			},
			updateErr:  nil,
			body:       `{"email":"new@email.com"}`,
			wantStatus: http.StatusOK,
			wantBody: fmt.Sprintf(
				`{"id":1,"email":"new@email.com","username":"olduser","created_at":"%s","updated_at":"%s"}`,
				fixedTime.Format(time.RFC3339), updatedTime.Format(time.RFC3339)),
			wantUpdate: &domain.UserUpdate{Email: strPtr("new@email.com")},
		},
		{
			name:       "unknown field in body",
			method:     http.MethodPut,
			path:       "/users/1",
			repoUser:   nil,
			updateErr:  nil,
			body:       `{"mail":"new@email.com"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid JSON\n",
			wantUpdate: nil,
		},
		{
			name:       "invalid method",
			method:     http.MethodPatch,
			path:       "/users/1",
			repoUser:   nil,
			updateErr:  nil,
			body:       `{"email":"new@email.com"}`,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed\n",
			wantUpdate: nil,
		},
		{
			name:       "user not found",
			method:     http.MethodPut,
			path:       "/users/9999",
			repoUser:   nil,
			updateErr:  fmt.Errorf("user not found"),
			body:       `{"email":"notfound@email.com"}`,
			wantStatus: http.StatusNotFound,
			wantBody:   "User not found\n",
			wantUpdate: &domain.UserUpdate{Email: strPtr("notfound@email.com")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotUpdate *domain.UserUpdate

			repo := &mockUserRepo{
				updateFunc: func(id int64, upd *domain.UserUpdate) error {
					gotUpdate = upd
					return test.updateErr
				},
				getByIDFunc: func(id int64) (*domain.User, error) {
					return test.repoUser, nil
				},
			}
			handler := NewUserHandler(repo)
			req := httptest.NewRequest(test.method, test.path, strings.NewReader(test.body))

			w := httptest.NewRecorder()
			handler.Update(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if test.wantStatus != resp.StatusCode {
				t.Errorf("expected status code: %v, got: %v", test.wantStatus, resp.StatusCode)
			}

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			respBodyStr := string(respBody)

			if strings.TrimSpace(respBodyStr) != strings.TrimSpace(test.wantBody) {
				t.Errorf("expected response body: %s, got: %s", test.wantBody, respBodyStr)
			}

			if test.wantUpdate != nil && gotUpdate != nil {
				if !reflect.DeepEqual(test.wantUpdate, gotUpdate) {
					t.Errorf("expected update: %+v, got: %+v", test.wantUpdate, gotUpdate)
				}
			}
		})
	}
}

type mockUserRepo struct {
	createFunc  func(*domain.User) error
	getByIDFunc func(int64) (*domain.User, error)
	updateFunc  func(int64, *domain.UserUpdate) error
	deleteFunc  func(int64) error
}

func (m *mockUserRepo) Create(u *domain.User) error {
	if m.createFunc != nil {
		u.ID = 1
		u.CreatedAt = fixedTime
		u.UpdatedAt = fixedTime
		return m.createFunc(u)
	}
	return nil
}

func (m *mockUserRepo) GetByID(id int64) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return nil, nil
}

func (m *mockUserRepo) Update(id int64, upd *domain.UserUpdate) error {
	if m.updateFunc != nil {
		return m.updateFunc(id, upd)
	}
	return nil
}

func (m *mockUserRepo) Delete(id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}
