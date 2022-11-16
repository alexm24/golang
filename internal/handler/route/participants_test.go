package route

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/service"
	mockService "github.com/alexm24/golang/internal/service/mocks"
)

func TestRoute_GetParticipantsByChannel(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIParticipants, channel string)

	id := uuid.New()
	username := "test"
	fullname := "test test"
	email := "test@vp.ru"

	participants := []models.Participant{
		{
			SIdentifier: api.SIdentifier{Id: &id},
			SUsername:   api.SUsername{Username: &username},
			SFullname:   api.SFullname{Fullname: &fullname},
			SEMail:      api.SEMail{Email: &email},
		},
	}

	jsonParticipants, _ := json.Marshal(participants)

	tests := []struct {
		name                 string
		channel              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			channel: "test",
			mockBehavior: func(r *mockService.MockIParticipants, channel string) {
				r.EXPECT().GetParticipants(channel).Return(participants, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonParticipants) + "\n",
		},
		{
			name:    "Service failure",
			channel: "test",
			mockBehavior: func(r *mockService.MockIParticipants, channel string) {
				r.EXPECT().GetParticipants(channel).Return(participants, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetParticipants + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIParticipants := mockService.NewMockIParticipants(c)
			test.mockBehavior(mockIParticipants, test.channel)

			services := &service.Service{IParticipants: mockIParticipants}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Get("/participants/{channel}", func(w http.ResponseWriter, r *http.Request) {
				ch := chi.URLParam(r, "channel")
				handler.GetParticipantsByChannel(w, r, ch)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/participants/" + test.channel
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestRoute_PostParticipantsByChannel(t *testing.T) {
	//Init Test Table
	type mockBehavior func(r *mockService.MockIParticipants, channel string, user models.PostParticipant)

	username := "test"
	email := "test@mail.ru"
	fullname := "test test"

	user := models.PostParticipant{
		Username: &username,
		Email:    &email,
		Fullname: &fullname,
	}

	userWithoutFullname := models.PostParticipant{
		Username: &username,
		Email:    &email,
	}
	userWithoutEmail := models.PostParticipant{
		Username: &username,
		Fullname: &fullname,
	}
	userName := models.PostParticipant{
		Email:    &email,
		Fullname: &fullname,
	}

	jsonUser, _ := json.Marshal(user)
	jsonUserWithoutFullname, _ := json.Marshal(userWithoutFullname)
	jsonUserWithoutEmail, _ := json.Marshal(userWithoutEmail)
	jsonWithoutUsername, _ := json.Marshal(userName)

	tests := []struct {
		name                 string
		channel              string
		inputBody            string
		inputUser            models.PostParticipant
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			channel:   "test",
			inputBody: string(jsonUser),
			inputUser: user,
			mockBehavior: func(r *mockService.MockIParticipants, channel string, user models.PostParticipant) {
				r.EXPECT().CreateParticipant(channel, user).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:      "Service failure",
			channel:   "test",
			inputBody: string(jsonUser),
			inputUser: user,
			mockBehavior: func(r *mockService.MockIParticipants, channel string, user models.PostParticipant) {
				r.EXPECT().CreateParticipant(channel, user).Return(errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceCreateParticipant + `"}` + "\n",
		},
		{
			name:                 "Fullname field is empty",
			channel:              "test",
			inputBody:            string(jsonUserWithoutFullname),
			inputUser:            userWithoutFullname,
			mockBehavior:         func(r *mockService.MockIParticipants, channel string, user models.PostParticipant) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgFullnameEmpty + `"}` + "\n",
		},
		{
			name:                 "Email field is empty",
			channel:              "test",
			inputBody:            string(jsonUserWithoutEmail),
			inputUser:            userWithoutEmail,
			mockBehavior:         func(r *mockService.MockIParticipants, channel string, user models.PostParticipant) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgEmailEmpty + `"}` + "\n",
		},
		{
			name:                 "Username field is empty",
			channel:              "test",
			inputBody:            string(jsonWithoutUsername),
			inputUser:            userName,
			mockBehavior:         func(r *mockService.MockIParticipants, channel string, user models.PostParticipant) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIParticipants := mockService.NewMockIParticipants(c)
			test.mockBehavior(mockIParticipants, test.channel, test.inputUser)

			services := &service.Service{IParticipants: mockIParticipants}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Post("/participants/{channel}", func(w http.ResponseWriter, r *http.Request) {
				ch := chi.URLParam(r, "channel")
				handler.PostParticipantsByChannel(w, r, ch)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/participants/" + test.channel
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}
