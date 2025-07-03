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

var fixedTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		repoErr    error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			method:     http.MethodPost,
			body:       `{"username": "test", "email": "test@email.net", "password": "mypw123"}`,
			repoErr:    nil,
			wantStatus: http.StatusCreated,
			wantBody: fmt.Sprintf(`{"id":1,"email":"test@email.net","username":"test","created_at":"%s","updated_at":"%s"}`,
				fixedTime.Format(time.RFC3339), fixedTime.Format(time.RFC3339)),
		},
		{
			name:       "invalid method",
			method:     http.MethodGet,
			body:       "",
			repoErr:    nil,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "Method not allowed\n",
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       `{bad json}`,
			repoErr:    nil,
			wantStatus: http.StatusBadRequest,
			wantBody:   "Invalid JSON\n",
		},
		{
			name:       "repository error",
			method:     http.MethodPost,
			body:       `{"username": "test", "email": "test@email.net"}`,
			repoErr:    errors.New("repo error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "Failed to create user: repo error\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &mockUserRepo{
				createFunc: func(*domain.User) error {
					return test.repoErr
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

			respBody, _ := io.ReadAll(resp.Body)
			respBodyStr := string(respBody)

			var got, want map[string]any
			json.Unmarshal([]byte(respBodyStr), &got)
			json.Unmarshal([]byte(test.wantBody), &want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("expected response: %v, got: %v", want, got)
			}
		})
	}
}

type mockUserRepo struct {
	createFunc  func(*domain.User) error
	getByIDFunc func(int64) (*domain.User, error)
	getByEmail  func(string) (*domain.User, error)
	updateFunc  func(*domain.User) error
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

func (m *mockUserRepo) Update(u *domain.User) error {
	if m.updateFunc != nil {
		return m.updateFunc(u)
	}
	return nil
}

func (m *mockUserRepo) Delete(id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error) {
	if m.getByEmail != nil {
		return m.getByEmail(email)
	}
	return nil, nil
}
