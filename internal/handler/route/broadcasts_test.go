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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/service"
	mockService "github.com/alexm24/golang/internal/service/mocks"
)

func TestRoute_GetBroadcasts(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIBroadcasts)

	id := uuid.New()
	desc := "test"
	previewUrl := "https://avatar/image/qwerty"
	life := "created"
	name := "test"
	owner := "test"
	date := time.Now()
	key := "test"

	broadcasts := []models.Broadcasts{
		{
			SIdentifier: api.SIdentifier{Id: &id},
			SPreviewUrl: api.SPreviewUrl{PreviewUrl: &previewUrl},
			SLifeCycle:  api.SLifeCycle{Life: &life},
			SStartTime:  api.SStartTime{StartTime: &date},
			SBroadcast: api.SBroadcast{
				Description: &desc,
				Name:        &name,
				Owner:       &owner,
				StreamKey:   &key,
			},
		},
	}

	jsonBroadcasts, _ := json.Marshal(broadcasts)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mockService.MockIBroadcasts) {
				r.EXPECT().GetBroadcasts().Return(broadcasts, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonBroadcasts) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIBroadcasts) {
				r.EXPECT().GetBroadcasts().Return(broadcasts, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetBroadcasts + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIBroadcasts := mockService.NewMockIBroadcasts(c)
			test.mockBehavior(mockIBroadcasts)

			services := &service.Service{IBroadcasts: mockIBroadcasts}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/broadcasts"
			r.Get(path, handler.GetBroadcasts)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}

}

func TestRoute_PostUserGetBroadcastArch(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIBroadcasts, username api.SUsername)

	user := "test"
	username := api.SUsername{Username: &user}
	jsonUsername, _ := json.Marshal(username)

	id := uuid.New()
	desc := "test"
	previewUrl := "https://avatar/image/qwerty"
	life := "created"
	name := "test"
	owner := "test"
	date := time.Now()
	key := "test"

	broadcasts := []models.Broadcasts{
		{
			SIdentifier: api.SIdentifier{Id: &id},
			SPreviewUrl: api.SPreviewUrl{PreviewUrl: &previewUrl},
			SLifeCycle:  api.SLifeCycle{Life: &life},
			SStartTime:  api.SStartTime{StartTime: &date},
			SBroadcast: api.SBroadcast{
				Description: &desc,
				Name:        &name,
				Owner:       &owner,
				StreamKey:   &key,
			},
		},
	}

	jsonBroadcasts, _ := json.Marshal(broadcasts)

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
			mockBehavior: func(r *mockService.MockIBroadcasts, u api.SUsername) {
				r.EXPECT().GetArchBroadcasts(u).Return(broadcasts, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: string(jsonBroadcasts) + "\n",
		},
		{
			name:                 "Username field is empty",
			inputBody:            `{}`,
			inputUser:            api.SUsername{},
			mockBehavior:         func(r *mockService.MockIBroadcasts, username api.SUsername) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonUsername),
			inputUser: username,
			mockBehavior: func(r *mockService.MockIBroadcasts, u api.SUsername) {
				r.EXPECT().GetArchBroadcasts(u).Return(broadcasts, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetArchBroadcasts + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIBroadcasts := mockService.NewMockIBroadcasts(c)
			test.mockBehavior(mockIBroadcasts, test.inputUser)

			services := &service.Service{IBroadcasts: mockIBroadcasts}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/broadcasts/arch"
			r.Post(path, handler.PostUserGetBroadcastArch)

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

func TestRoute_PostBroadcasts(t *testing.T) {
	//Init Test Table
	type mockBehavior func(r *mockService.MockIBroadcasts, i models.PostBroadcast)

	id := uuid.New()
	life := "created"
	url := "https://vp.ru/avatar/test"

	name := "test"
	desc := "test broadcast"
	owner := "test@vp.ru"
	streamKey := "stream"
	sb := time.Date(2022, time.July, 1, 8, 0, 0, 0, time.UTC)

	broadcast := models.PostBroadcast{
		Description: &desc,
		Name:        &name,
		Owner:       &owner,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}

	resBroadcast := models.Broadcasts{
		SIdentifier: api.SIdentifier{Id: &id},
		SPreviewUrl: api.SPreviewUrl{PreviewUrl: &url},
		SBroadcast: api.SBroadcast{
			Description: &desc,
			Name:        &name,
			Owner:       &owner,
			StreamKey:   &streamKey,
		},
		SLifeCycle: api.SLifeCycle{Life: &life},
		SStartTime: api.SStartTime{StartTime: &sb},
	}

	broadcastWithoutName := models.PostBroadcast{
		Description: &desc,
		Owner:       &owner,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}

	broadcastWithoutDesc := models.PostBroadcast{
		Name:      &name,
		Owner:     &owner,
		StartTime: &sb,
		StreamKey: &streamKey,
	}

	broadcastWithoutOwner := models.PostBroadcast{
		Description: &desc,
		Name:        &name,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}

	broadcastWithoutStreamKey := models.PostBroadcast{
		Description: &desc,
		Name:        &name,
		Owner:       &owner,
		StartTime:   &sb,
	}

	broadcastWithoutStartTime := models.PostBroadcast{
		Description: &desc,
		Name:        &name,
		Owner:       &owner,
		StreamKey:   &streamKey,
	}

	jsonBroadcastWithoutName, _ := json.Marshal(broadcastWithoutName)
	jsonBroadcastWithoutDesc, _ := json.Marshal(broadcastWithoutDesc)
	jsonBroadcastWithoutOwner, _ := json.Marshal(broadcastWithoutOwner)
	jsonBroadcastWithoutStreamKey, _ := json.Marshal(broadcastWithoutStreamKey)
	jsonBroadcastWithoutStartTime, _ := json.Marshal(broadcastWithoutStartTime)

	jsonBroadcast, _ := json.Marshal(broadcast)
	jsonResBroadcast, _ := json.Marshal(resBroadcast)

	tests := []struct {
		name                 string
		inputBody            string
		inputBroadcast       models.PostBroadcast
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "Ok",
			inputBody:      string(jsonBroadcast),
			inputBroadcast: broadcast,
			mockBehavior: func(r *mockService.MockIBroadcasts, i models.PostBroadcast) {
				r.EXPECT().CreateBroadcast(i).Return(resBroadcast, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonResBroadcast) + "\n",
		},
		{
			name:           "Service failure",
			inputBody:      string(jsonBroadcast),
			inputBroadcast: broadcast,
			mockBehavior: func(r *mockService.MockIBroadcasts, b models.PostBroadcast) {
				r.EXPECT().CreateBroadcast(b).Return(resBroadcast, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceCreateBroadcast + `"}` + "\n",
		},
		{
			name:                 "Name field is empty",
			inputBody:            string(jsonBroadcastWithoutName),
			inputBroadcast:       broadcastWithoutName,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PostBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgNameEmpty + `"}` + "\n",
		},
		{
			name:                 "Description field is empty",
			inputBody:            string(jsonBroadcastWithoutDesc),
			inputBroadcast:       broadcastWithoutDesc,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PostBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgDescriptionEmpty + `"}` + "\n",
		},
		{
			name:                 "Owner field is empty",
			inputBody:            string(jsonBroadcastWithoutOwner),
			inputBroadcast:       broadcastWithoutOwner,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PostBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgOwnerEmpty + `"}` + "\n",
		},
		{
			name:                 "StreamKey field is empty",
			inputBody:            string(jsonBroadcastWithoutStreamKey),
			inputBroadcast:       broadcastWithoutStreamKey,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PostBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgStreamKeyEmpty + `"}` + "\n",
		},
		{
			name:                 "StartBroadcast field is empty",
			inputBody:            string(jsonBroadcastWithoutStartTime),
			inputBroadcast:       broadcastWithoutStartTime,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PostBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgStartTimeEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIBroadcasts := mockService.NewMockIBroadcasts(c)
			test.mockBehavior(mockIBroadcasts, test.inputBroadcast)

			services := &service.Service{IBroadcasts: mockIBroadcasts}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/broadcasts"
			r.Post(path, handler.PostBroadcasts)

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

func TestRoute_DeleteBroadcast(t *testing.T) {
	//Init Test Table
	type mockBehavior func(r *mockService.MockIBroadcasts, id uuid.UUID)

	params := uuid.New()
	id := api.SIdentifier{Id: &params}

	jsonId, _ := json.Marshal(id)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mockService.MockIBroadcasts, i uuid.UUID) {
				r.EXPECT().DeleteBroadcast(i).Return(id, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonId) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIBroadcasts, i uuid.UUID) {
				r.EXPECT().DeleteBroadcast(i).Return(id, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceDeleteBroadcast + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIBroadcasts := mockService.NewMockIBroadcasts(c)
			test.mockBehavior(mockIBroadcasts, params)

			services := &service.Service{IBroadcasts: mockIBroadcasts}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/broadcasts/" + params.String()
			r.Delete("/broadcasts/{id}", func(w http.ResponseWriter, r *http.Request) {
				ID := chi.URLParam(r, "id")
				uuID, _ := uuid.Parse(ID)
				handler.DeleteBroadcast(w, r, uuID)
			})

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, path, nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}

func TestRoute_PutBroadcast(t *testing.T) {
	//Init Test Table
	type mockBehavior func(r *mockService.MockIBroadcasts, b models.PutBroadcast)

	id := uuid.New()
	life := "created"
	url := "https://vp.ru/avatar/test"

	name := "test"
	desc := "test broadcast"
	owner := "test@vp.ru"
	streamKey := "stream"
	sb := time.Date(2022, time.July, 1, 8, 0, 0, 0, time.UTC)

	putBroadcast := models.PutBroadcast{
		Description: &desc,
		Id:          &id,
		Name:        &name,
		Owner:       &owner,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}
	putBroadcastWithoutId := models.PutBroadcast{
		Description: &desc,
		Name:        &name,
		Owner:       &owner,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}
	putBroadcastWithoutName := models.PutBroadcast{
		Description: &desc,
		Id:          &id,
		Owner:       &owner,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}
	putBroadcastWithoutDesc := models.PutBroadcast{
		Id:        &id,
		Name:      &name,
		Owner:     &owner,
		StartTime: &sb,
		StreamKey: &streamKey,
	}
	putBroadcastWithoutOwner := models.PutBroadcast{
		Description: &desc,
		Id:          &id,
		Name:        &name,
		StartTime:   &sb,
		StreamKey:   &streamKey,
	}
	putBroadcastWithoutStreamKey := models.PutBroadcast{
		Description: &desc,
		Id:          &id,
		Name:        &name,
		Owner:       &owner,
		StartTime:   &sb,
	}
	putBroadcastWithoutStartTime := models.PutBroadcast{
		Description: &desc,
		Id:          &id,
		Name:        &name,
		Owner:       &owner,
		StreamKey:   &streamKey,
	}

	broadcast := models.Broadcasts{
		SIdentifier: api.SIdentifier{Id: &id},
		SPreviewUrl: api.SPreviewUrl{PreviewUrl: &url},
		SStartTime:  api.SStartTime{StartTime: &sb},
		SBroadcast: api.SBroadcast{
			Description: &desc,
			Name:        &name,
			Owner:       &owner,
			StreamKey:   &streamKey,
		},
		SLifeCycle: api.SLifeCycle{Life: &life},
	}

	jsonPutBroadcast, _ := json.Marshal(putBroadcast)
	jsonBroadcast, _ := json.Marshal(broadcast)

	jsonPutBroadcastWithoutId, _ := json.Marshal(putBroadcastWithoutId)
	jsonPutBroadcastWithoutName, _ := json.Marshal(putBroadcastWithoutName)
	jsonPutBroadcastWithoutDesc, _ := json.Marshal(putBroadcastWithoutDesc)
	jsonPutBroadcastWithoutOwner, _ := json.Marshal(putBroadcastWithoutOwner)
	jsonPutBroadcastWithoutStreamKey, _ := json.Marshal(putBroadcastWithoutStreamKey)
	jsonPutBroadcastWithoutStartTime, _ := json.Marshal(putBroadcastWithoutStartTime)

	tests := []struct {
		name                 string
		inputBody            string
		inputBroadcast       models.PutBroadcast
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:           "Ok",
			inputBody:      string(jsonPutBroadcast),
			inputBroadcast: putBroadcast,
			mockBehavior: func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {
				r.EXPECT().ChangeBroadcast(b).Return(broadcast, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonBroadcast) + "\n",
		},
		{
			name:           "Service failure",
			inputBody:      string(jsonPutBroadcast),
			inputBroadcast: putBroadcast,
			mockBehavior: func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {
				r.EXPECT().ChangeBroadcast(b).Return(broadcast, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceChangeBroadcast + `"}` + "\n",
		},
		{
			name:                 "Id field is empty",
			inputBody:            string(jsonPutBroadcastWithoutId),
			inputBroadcast:       putBroadcastWithoutId,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgIdEmpty + `"}` + "\n",
		},
		{
			name:                 "Name field is empty",
			inputBody:            string(jsonPutBroadcastWithoutName),
			inputBroadcast:       putBroadcastWithoutName,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgNameEmpty + `"}` + "\n",
		},
		{
			name:                 "Description field is empty",
			inputBody:            string(jsonPutBroadcastWithoutDesc),
			inputBroadcast:       putBroadcastWithoutDesc,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgDescriptionEmpty + `"}` + "\n",
		},
		{
			name:                 "Owner field is empty",
			inputBody:            string(jsonPutBroadcastWithoutOwner),
			inputBroadcast:       putBroadcastWithoutOwner,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgOwnerEmpty + `"}` + "\n",
		},
		{
			name:                 "StreamKey field is empty",
			inputBody:            string(jsonPutBroadcastWithoutStreamKey),
			inputBroadcast:       putBroadcastWithoutStreamKey,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgStreamKeyEmpty + `"}` + "\n",
		},
		{
			name:                 "StartBroadcast field is empty",
			inputBody:            string(jsonPutBroadcastWithoutStartTime),
			inputBroadcast:       putBroadcastWithoutStartTime,
			mockBehavior:         func(r *mockService.MockIBroadcasts, b models.PutBroadcast) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgStartTimeEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIBroadcasts := mockService.NewMockIBroadcasts(c)
			test.mockBehavior(mockIBroadcasts, test.inputBroadcast)

			services := &service.Service{IBroadcasts: mockIBroadcasts}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/broadcasts"
			r.Post(path, handler.PutBroadcast)

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

func TestRoute_GetBroadcastById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIBroadcasts, id uuid.UUID)

	id := uuid.New()
	desc := "test"
	previewUrl := "https://avatar/image/qwerty"
	life := "created"
	name := "test"
	owner := "test"
	date := time.Now()
	key := "test"

	broadcast := models.Broadcasts{
		SIdentifier: api.SIdentifier{Id: &id},
		SPreviewUrl: api.SPreviewUrl{PreviewUrl: &previewUrl},
		SLifeCycle:  api.SLifeCycle{Life: &life},
		SStartTime:  api.SStartTime{StartTime: &date},
		SBroadcast: api.SBroadcast{
			Description: &desc,
			Name:        &name,
			Owner:       &owner,
			StreamKey:   &key,
		},
	}

	jsonBroadcast, _ := json.Marshal(broadcast)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mockService.MockIBroadcasts, id uuid.UUID) {
				r.EXPECT().GetBroadcastById(id).Return(broadcast, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonBroadcast) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIBroadcasts, id uuid.UUID) {
				r.EXPECT().GetBroadcastById(id).Return(broadcast, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetBroadcastById + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIBroadcasts := mockService.NewMockIBroadcasts(c)
			test.mockBehavior(mockIBroadcasts, id)

			services := &service.Service{IBroadcasts: mockIBroadcasts}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Get("/broadcasts/{id}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				uuId, _ := uuid.Parse(iD)
				handler.GetBroadcastById(w, r, uuId)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/broadcasts/" + id.String()
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}
