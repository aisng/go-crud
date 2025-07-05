package handler

import (
	"encoding/json"
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

func checkResponseFields(t *testing.T, respBodyStr string, wantFields map[string]any) {
	var got map[string]any
	if err := json.Unmarshal([]byte(respBodyStr), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	for field, want := range wantFields {
		gotVal, ok := got[field]
		if !ok {
			t.Errorf("response missing expected field %s: %v", field, got)
			continue
		}
		if fmt.Sprint(gotVal) != fmt.Sprint(want) {
			t.Errorf("field %s: expected %v, got %v", field, want, gotVal)
		}
	}
}

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantErr    error
		wantStatus int
		wantFields map[string]any
	}{
		{
			name:       "success",
			body:       `{"username": "test", "email": "test@email.net", "password": "mypw123"}`,
			wantErr:    nil,
			wantStatus: http.StatusCreated,
			wantFields: map[string]any{
				"id":       1,
				"email":    "test@email.net",
				"username": "test",
			},
		},
		{
			name:       "invalid json",
			body:       `{bad json}`,
			wantErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantFields: nil,
		},
		{
			name:       "repository error",
			body:       `{"username": "test", "email": "test@email.net"}`,
			wantErr:    errors.New("repo error"),
			wantStatus: http.StatusInternalServerError,
			wantFields: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var repo *mockUserRepo
			if test.wantErr != nil {
				repo = &mockUserRepo{
					createFunc: func(*domain.User) error {
						return test.wantErr
					},
				}
			} else {
				repo = &mockUserRepo{}
			}

			handler := NewUserHandler(repo)

			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(test.body))
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

			if test.wantFields != nil {
				checkResponseFields(t, respBodyStr, test.wantFields)
			} else {
				var expectedBody string
				switch test.name {
				case "invalid json":
					expectedBody = "Invalid JSON\n"
				case "repository error":
					expectedBody = "Failed to create user: repo error\n"
				}
				if strings.TrimSpace(respBodyStr) != strings.TrimSpace(expectedBody) {
					t.Errorf("expected response body: %s, got: %s", expectedBody, respBodyStr)
				}
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		repoUser   *domain.User
		wantErr    error
		wantStatus int
		wantFields map[string]any
	}{
		{
			name: "success",
			path: "/users/1",
			repoUser: &domain.User{
				ID:       1,
				Email:    "test@email.com",
				Username: "testuser",
			},
			wantErr:    nil,
			wantStatus: http.StatusOK,
			wantFields: map[string]any{
				"id":       1,
				"email":    "test@email.com",
				"username": "testuser",
			},
		},
		{
			name:       "user not found",
			path:       "/users/9999",
			repoUser:   nil,
			wantErr:    fmt.Errorf("user not found"),
			wantStatus: http.StatusNotFound,
			wantFields: nil,
		},
		{
			name:       "invalid path",
			path:       "/invalid/path/1",
			repoUser:   nil,
			wantErr:    fmt.Errorf("invalid path format"),
			wantStatus: http.StatusBadRequest,
			wantFields: nil,
		},
		{
			name:       "invalid id",
			path:       "/users/1b",
			repoUser:   nil,
			wantErr:    fmt.Errorf("Invalid parameter 'id'\n"),
			wantStatus: http.StatusBadRequest,
			wantFields: nil,
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

			req := httptest.NewRequest(http.MethodGet, test.path, nil)
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

			if test.wantFields != nil {
				checkResponseFields(t, respBodyStr, test.wantFields)
			} else {
				var expectedBody string
				switch test.name {
				case "user not found":
					expectedBody = "User not found\n"
				case "invalid path":
					expectedBody = "invalid path format\n"
				case "invalid id":
					expectedBody = "Invalid parameter 'id'\n"
				}
				if strings.TrimSpace(respBodyStr) != strings.TrimSpace(expectedBody) {
					t.Errorf("expected response body: %s, got: %s", expectedBody, respBodyStr)
				}
			}
		})
	}
}

func strPtr(s string) *string { return &s }

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		repoUser   *domain.User
		updateErr  error
		body       string
		wantStatus int
		wantFields map[string]any
		wantUpdate *domain.UserUpdate
	}{
		{
			name: "update email success",
			path: "/users/1",
			repoUser: &domain.User{
				ID:       1,
				Username: "olduser",
				Email:    "new@email.com",
			},
			updateErr:  nil,
			body:       `{"email":"new@email.com"}`,
			wantStatus: http.StatusOK,
			wantFields: map[string]any{
				"id":       1,
				"email":    "new@email.com",
				"username": "olduser",
			},
			wantUpdate: &domain.UserUpdate{Email: strPtr("new@email.com")},
		},
		{
			name:       "unknown field in body",
			path:       "/users/1",
			repoUser:   nil,
			updateErr:  nil,
			body:       `{"mail":"new@email.com"}`,
			wantStatus: http.StatusBadRequest,
			wantFields: nil,
			wantUpdate: nil,
		},
		{
			name:       "user not found",
			path:       "/users/9999",
			repoUser:   nil,
			updateErr:  fmt.Errorf("user not found"),
			body:       `{"email":"notfound@email.com"}`,
			wantStatus: http.StatusNotFound,
			wantFields: nil,
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
			req := httptest.NewRequest(http.MethodPut, test.path, strings.NewReader(test.body))

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

			if test.wantFields != nil {
				checkResponseFields(t, respBodyStr, test.wantFields)
			} else {
				var expectedBody string
				switch test.name {
				case "unknown field in body":
					expectedBody = "Invalid JSON\n"
				case "user not found":
					expectedBody = "User not found\n"
				}
				if strings.TrimSpace(respBodyStr) != strings.TrimSpace(expectedBody) {
					t.Errorf("expected response body: %s, got: %s", expectedBody, respBodyStr)
				}
			}

			if test.wantUpdate != nil && gotUpdate != nil {
				if !reflect.DeepEqual(test.wantUpdate, gotUpdate) {
					t.Errorf("expected update: %+v, got: %+v", test.wantUpdate, gotUpdate)
				}
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		deleteErr  error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			path:       "/users/1",
			deleteErr:  nil,
			wantStatus: http.StatusNoContent,
			wantBody:   "",
		},
		{
			name:       "user not found",
			path:       "/users/9999",
			deleteErr:  fmt.Errorf("user not found"),
			wantStatus: http.StatusNotFound,
			wantBody:   "User not found\n",
		},
		{
			name:       "invalid path",
			path:       "/invalid/path/1",
			deleteErr:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid path format\n",
		},
		{
			name:       "invalid id",
			path:       "/users/1b",
			deleteErr:  nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid parameter 'id'\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				deleteFunc: func(id int64) error {
					return test.deleteErr
				},
			}
			handler := NewUserHandler(repo)

			req := httptest.NewRequest(http.MethodDelete, test.path, nil)
			w := httptest.NewRecorder()

			handler.Delete(w, req)

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
		return m.createFunc(u)
	}
	u.ID = 1
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
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
