package handler_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"GO-User_service/user-service/internal/handler"
	"GO-User_service/user-service/internal/usersdb"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetStatus(t *testing.T) {
	logger := &log.Logger{}
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	h := handler.NewHandler(nil, logger)

	h.GetStatus(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestHandler_GetUsers(t *testing.T) {
	logger := &log.Logger{}
	tests := []struct {
		name               string
		db                 handler.Database
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "successful_test",
			db:                 &mockDatabase{},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `[{"username":"username","password":"password","privileges":0}]`,
		},
		{
			name: "unsuccessful_test",
			db: &mockDatabase{
				expectFail: true,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `failed to fetch users from database:`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()
			h := handler.NewHandler(tc.db, logger)

			h.GetUsers(w, req)

			res := w.Result()
			defer res.Body.Close()
			resp, err := io.ReadAll(res.Body)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.Equal(t, res.StatusCode, tc.expectedStatusCode)
			assert.Contains(t, string(resp), tc.expectedResponse)
		})
	}
}

func TestHandler_GetUsersUsername(t *testing.T) {
	logger := &log.Logger{}
	tests := []struct {
		name               string
		db                 handler.Database
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "successful_test",
			db:                 &mockDatabase{},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"username":"username","password":"password","privileges":0}`,
		},
		{
			name: "unsuccessful_test",
			db: &mockDatabase{
				expectFail: true,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `failed to fetch user:`,
		},
		{
			name: "unsuccessful_test",
			db: &mockDatabase{
				expectFailedUsernameCheck: true,
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `failed to get user with username`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()
			h := handler.NewHandler(tc.db, logger)

			h.GetUsersUsername(w, req, "")

			res := w.Result()
			defer res.Body.Close()
			resp, err := io.ReadAll(res.Body)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.Equal(t, res.StatusCode, tc.expectedStatusCode)
			assert.Contains(t, string(resp), tc.expectedResponse)
		})
	}
}

//
//func TestHandler_PostUser(t *testing.T) {
//	type fields struct {
//		usersdb database
//		logger  *log.Logger
//	}
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &Handler{
//				usersdb: tt.fields.usersdb,
//				logger:  tt.fields.logger,
//			}
//			h.PostUser(tt.args.w, tt.args.r)
//		})
//	}
//}

type mockDatabase struct {
	expectFail                bool
	expectFailedUsernameCheck bool
	expectIncorrectUser       bool
}

func (m *mockDatabase) CreateUserIfNotExists(_ usersdb.User) error {
	if m.expectFail {
		return errors.New("error")
	}
	return nil
}

func (m *mockDatabase) GetUser(_ string) (*usersdb.User, error) {
	if m.expectFail {
		return nil, errors.New("error")
	}
	if m.expectIncorrectUser {
		return nil, nil
	}
	return &usersdb.User{
		Username:   "username",
		Password:   "password",
		Privileges: 0,
	}, nil
}

func (m *mockDatabase) CheckUsername(_ string) (bool, error) {
	if m.expectFailedUsernameCheck {
		return false, errors.New("error")
	}
	return true, nil
}

func (m *mockDatabase) GetAllUsers() ([]usersdb.User, error) {
	if m.expectFail {
		return nil, errors.New("error")
	}
	return []usersdb.User{
		{
			Username:   "username",
			Password:   "password",
			Privileges: 0,
		},
	}, nil
}

func (m *mockDatabase) Close() error {
	if m.expectFail {
		return errors.New("error")
	}
	return nil
}
