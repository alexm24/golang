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

func TestRoute_PostStream(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIStream, username api.SUsername)

	user := "test"
	username := api.SUsername{Username: &user}

	id := uuid.New()

	stream := models.Stream{
		SIdentifier:  api.SIdentifier{Id: &id},
		SUsername:    username,
		SDescription: api.SDescription{Description: &user},
	}

	jsonStream, _ := json.Marshal(stream)
	jsonUsername, _ := json.Marshal(username)

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
			mockBehavior: func(r *mockService.MockIStream, user api.SUsername) {
				r.EXPECT().CreateStream(user).Return(stream, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonStream) + "\n",
		},
		{
			name:                 "Username field is empty",
			inputBody:            `{}`,
			inputUser:            api.SUsername{},
			mockBehavior:         func(r *mockService.MockIStream, username api.SUsername) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonUsername),
			inputUser: username,
			mockBehavior: func(r *mockService.MockIStream, user api.SUsername) {
				r.EXPECT().CreateStream(user).Return(stream, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceCreateStream + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIStream := mockService.NewMockIStream(c)
			test.mockBehavior(mockIStream, test.inputUser)

			services := &service.Service{IStream: mockIStream}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/stream"
			r.Post(path, handler.PostStream)

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

func TestRoute_GetStreamByUsername(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIStream, username string)

	username := "test"
	id := uuid.New()

	stream := models.Stream{
		SIdentifier:  api.SIdentifier{Id: &id},
		SUsername:    api.SUsername{Username: &username},
		SDescription: api.SDescription{Description: &username},
	}

	itemJson, _ := json.Marshal(stream)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mockService.MockIStream, username string) {
				r.EXPECT().GetStream(username).Return(stream, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(itemJson) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIStream, username string) {
				r.EXPECT().GetStream(username).Return(stream, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetStream + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIStream := mockService.NewMockIStream(c)
			test.mockBehavior(mockIStream, username)

			services := &service.Service{IStream: mockIStream}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Get("/stream/{username}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "username")
				handler.GetStreamByUsername(w, r, iD)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/stream/" + username
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestRoute_PutStream(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIStream, item models.PutStream)

	username := "test"
	id := uuid.New()

	stream := models.Stream{
		SIdentifier:  api.SIdentifier{Id: &id},
		SUsername:    api.SUsername{Username: &username},
		SDescription: api.SDescription{Description: &username},
	}

	putStream := models.PutStream{
		Description: &username,
		Username:    &username,
	}

	putStreamWithoutUsername := models.PutStream{
		Description: &username,
	}

	putStreamWithoutDesc := models.PutStream{
		Username: &username,
	}

	jsonStream, _ := json.Marshal(stream)
	jsonPutStream, _ := json.Marshal(putStream)
	jsonPutStreamWithoutUsername, _ := json.Marshal(putStreamWithoutUsername)
	jsonPutStreamWithoutDesc, _ := json.Marshal(putStreamWithoutDesc)

	tests := []struct {
		name                 string
		inputBody            string
		inputPutStream       models.PutStream
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "Ok",
			inputBody:      string(jsonPutStream),
			inputPutStream: putStream,
			mockBehavior: func(r *mockService.MockIStream, item models.PutStream) {
				r.EXPECT().ChangeDescByUsername(item).Return(stream, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonStream) + "\n",
		},
		{
			name:           "Service failure",
			inputBody:      string(jsonPutStream),
			inputPutStream: putStream,
			mockBehavior: func(r *mockService.MockIStream, item models.PutStream) {
				r.EXPECT().ChangeDescByUsername(item).Return(stream, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceChangeDescByUsername + `"}` + "\n",
		},
		{
			name:                 "Username field is empty",
			inputBody:            string(jsonPutStreamWithoutUsername),
			inputPutStream:       putStreamWithoutUsername,
			mockBehavior:         func(r *mockService.MockIStream, item models.PutStream) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:                 "Description field is empty",
			inputBody:            string(jsonPutStreamWithoutDesc),
			inputPutStream:       putStreamWithoutDesc,
			mockBehavior:         func(r *mockService.MockIStream, item models.PutStream) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgDescriptionEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIStream := mockService.NewMockIStream(c)
			test.mockBehavior(mockIStream, test.inputPutStream)

			services := &service.Service{IStream: mockIStream}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			path := "/stream"
			route.Put(path, handler.PutStream)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}

}

func TestRoute_DeleteStreamChat(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIStream, channel string)

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
			mockBehavior: func(r *mockService.MockIStream, channel string) {
				r.EXPECT().ClearChat(channel).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:    "Service failure",
			channel: "test",
			mockBehavior: func(r *mockService.MockIStream, channel string) {
				r.EXPECT().ClearChat(channel).Return(errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceClearChat + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIStream := mockService.NewMockIStream(c)
			test.mockBehavior(mockIStream, test.channel)

			services := &service.Service{IStream: mockIStream}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Delete("/stream/chat/{channel}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "channel")
				handler.DeleteStreamChat(w, r, iD)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/stream/chat/" + test.channel
			req := httptest.NewRequest(http.MethodDelete, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}

}
