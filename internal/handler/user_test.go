package handler

import (
	"encoding/json"
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
		repoErr    error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			body:       `{"username":"testuser","email":"test@email.com","password":"password123"}`,
			repoErr:    nil,
			wantStatus: http.StatusCreated,
			wantBody:   `{"id":1,"username":"testuser","email":"test@email.com"}`,
		},
		{
			name:       "invalid json",
			body:       `{"username":"testuser","email":"test@email.com"`,
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid json"}`,
		},
		{
			name:       "repository error",
			body:       `{"username":"testuser","email":"test@email.com","password":"password123"}`,
			repoErr:    fmt.Errorf("repo error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"internal server error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				createFunc: func(u *domain.User) error {
					if test.repoErr == nil {
						u.ID = 1
						u.CreatedAt = time.Now()
						u.UpdatedAt = time.Now()
					}
					return test.repoErr
				},
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
			if test.name == "success" {
				checkResponseFields(t, respBodyStr, map[string]any{
					"id":       1,
					"username": "testuser",
					"email":    "test@email.com",
				})
			} else {
				if strings.TrimSpace(respBodyStr) != strings.TrimSpace(test.wantBody) {
					t.Errorf("expected response body: %s, got: %s", test.wantBody, respBodyStr)
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
		handlerErr error
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			path: "/users/1",
			repoUser: &domain.User{
				ID:       1,
				Email:    "test@email.com",
				Username: "testuser",
			},
			handlerErr: nil,
			wantStatus: http.StatusOK,
			wantBody:   `{"id":1,"email":"test@email.com","username":"testuser"}`,
		},
		{
			name:       "user not found",
			path:       "/users/9999",
			repoUser:   nil,
			handlerErr: domain.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"resource not found"}`,
		},
		{
			name:       "invalid path",
			path:       "/invalid/path/1",
			repoUser:   nil,
			handlerErr: ErrInvalidPath,
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintf(`{"error":"%s"}`, ErrInvalidPath.Message),
		},
		{
			name:       "invalid id",
			path:       "/users/1b",
			repoUser:   nil,
			handlerErr: ErrInvalidID,
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintf(`{"error":"%s"}`, ErrInvalidID.Message),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				getByIDFunc: func(id int64) (*domain.User, error) {
					return test.repoUser, test.handlerErr
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
			if test.name == "success" {
				checkResponseFields(t, respBodyStr, map[string]any{
					"id":       1,
					"email":    "test@email.com",
					"username": "testuser",
				})
			} else {
				if strings.TrimSpace(respBodyStr) != strings.TrimSpace(test.wantBody) {
					t.Errorf("expected response body: %s, got: %s", test.wantBody, respBodyStr)
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
		handlerErr error
		body       string
		wantStatus int
		wantBody   string
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
			handlerErr: nil,
			body:       `{"email":"new@email.com"}`,
			wantStatus: http.StatusOK,
			wantBody:   `{"id":1,"email":"new@email.com","username":"olduser"}`,
			wantUpdate: &domain.UserUpdate{Email: strPtr("new@email.com")},
		},
		{
			name:       "unknown field in body",
			path:       "/users/1",
			repoUser:   nil,
			handlerErr: ErrInvalidJSON,
			body:       `{"mail":"new@email.com"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintf(`{"error":"%s"}`, ErrInvalidJSON.Message),
			wantUpdate: nil,
		},
		{
			name:       "user not found",
			path:       "/users/9999",
			repoUser:   nil,
			handlerErr: domain.ErrNotFound,
			body:       `{"email":"notfound@email.com"}`,
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"resource not found"}`,
			wantUpdate: &domain.UserUpdate{Email: strPtr("notfound@email.com")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotUpdate *domain.UserUpdate
			repo := &mockUserRepo{
				updateFunc: func(id int64, upd *domain.UserUpdate) error {
					gotUpdate = upd
					return test.handlerErr
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
			if test.name == "update email success" {
				checkResponseFields(t, respBodyStr, map[string]any{
					"id":       1,
					"email":    "new@email.com",
					"username": "olduser",
				})
			} else {
				if strings.TrimSpace(respBodyStr) != strings.TrimSpace(test.wantBody) {
					t.Errorf("expected response body: %s, got: %s", test.wantBody, respBodyStr)
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
		handlerErr error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			path:       "/users/1",
			handlerErr: nil,
			wantStatus: http.StatusNoContent,
			wantBody:   "",
		},
		{
			name:       "user not found",
			path:       "/users/9999",
			handlerErr: domain.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"resource not found"}`,
		},
		{
			name:       "invalid path",
			path:       "/invalid/path/1",
			handlerErr: ErrInvalidPath,
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintf(`{"error":"%s"}`, ErrInvalidPath.Message),
		},
		{
			name:       "invalid id",
			path:       "/users/1b",
			handlerErr: ErrInvalidID,
			wantStatus: http.StatusBadRequest,
			wantBody:   fmt.Sprintf(`{"error":"%s"}`, ErrInvalidID.Message),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				deleteFunc: func(id int64) error {
					return test.handlerErr
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
