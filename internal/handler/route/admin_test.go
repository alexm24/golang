package route

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/service"
	mockService "github.com/alexm24/golang/internal/service/mocks"
)

func TestRoute_CheckAdmin(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIAdmin, u api.SUsername)

	user := "test"

	isAdministrator := true

	isAdmin := api.SAdmin{IsAdmin: &isAdministrator}

	username := api.SUsername{Username: &user}

	jsonUsername, _ := json.Marshal(username)
	jsonIsAdmin, _ := json.Marshal(isAdmin)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            api.SUsername
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: string(jsonUsername),
			inputUser: username,
			mockBehavior: func(r *mockService.MockIAdmin, u api.SUsername) {
				r.EXPECT().CheckAdminUser(u).Return(isAdmin, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonIsAdmin) + "\n",
		},
		{
			name:                 "Username field is empty",
			inputBody:            `{}`,
			inputUser:            api.SUsername{},
			mockBehavior:         func(r *mockService.MockIAdmin, username api.SUsername) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonUsername),
			inputUser: username,
			mockBehavior: func(r *mockService.MockIAdmin, u api.SUsername) {
				r.EXPECT().CheckAdminUser(u).Return(isAdmin, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceCheckAdminUser + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIAdmin := mockService.NewMockIAdmin(c)
			test.mockBehavior(mockIAdmin, test.inputUser)

			services := &service.Service{IAdmin: mockIAdmin}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/admin"
			r.Post(path, handler.CheckAdmin)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}

}

func TestRoute_PostUserGetToken(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIAdmin, u api.SUsername)

	user := "test"
	token := "eyJhbGciOiJIUzI1.eyJleHAiOjE2NTQxNzg3NzYsInN1Y.St6KshOEaO3q-uSTJpSipftKEH3rD"
	exp := time.Now()

	objToken := api.SToken{
		Exp:   &exp,
		Token: &token,
	}

	username := api.SUsername{Username: &user}

	jsonUsername, _ := json.Marshal(username)
	jsonToken, _ := json.Marshal(objToken)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            api.SUsername
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: string(jsonUsername),
			inputUser: username,
			mockBehavior: func(r *mockService.MockIAdmin, u api.SUsername) {
				r.EXPECT().GetToken(u).Return(objToken, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonToken) + "\n",
		},
		{
			name:                 "Username field is empty",
			inputBody:            `{}`,
			inputUser:            api.SUsername{},
			mockBehavior:         func(r *mockService.MockIAdmin, username api.SUsername) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonUsername),
			inputUser: username,
			mockBehavior: func(r *mockService.MockIAdmin, u api.SUsername) {
				r.EXPECT().GetToken(u).Return(objToken, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetToken + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIAdmin := mockService.NewMockIAdmin(c)
			test.mockBehavior(mockIAdmin, test.inputUser)

			services := &service.Service{IAdmin: mockIAdmin}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/token"
			r.Post(path, handler.PostUserGetToken)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}

}
