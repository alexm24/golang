package route

import (
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

func TestRoute_GetLive(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockILive)

	id := uuid.New()
	desc := "test"
	place := "Rostov"
	url := "https://vp.ru"

	live := []models.Live{
		{
			SIdentifier:  api.SIdentifier{Id: &id},
			SDescription: api.SDescription{Description: &desc},
			SPlace:       api.SPlace{Place: &place},
			SStreamUrl:   api.SStreamUrl{StreamUrl: &url},
		},
	}

	jsonLive, _ := json.Marshal(live)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mockService.MockILive) {
				r.EXPECT().GetLive().Return(live, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonLive) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockILive) {
				r.EXPECT().GetLive().Return(live, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetLive + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockILive := mockService.NewMockILive(c)
			test.mockBehavior(mockILive)

			services := &service.Service{ILive: mockILive}
			handler := Route{services}

			// Init Endpoint
			r := chi.NewRouter()
			path := "/live"
			r.Get(path, handler.GetLive)

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

func TestRoute_GetLiveById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockILive, id uuid.UUID)

	id := uuid.New()
	desc := "test"
	place := "Rostov"
	url := "https://vp.ru"

	live := models.Live{
		SIdentifier:  api.SIdentifier{Id: &id},
		SDescription: api.SDescription{Description: &desc},
		SPlace:       api.SPlace{Place: &place},
		SStreamUrl:   api.SStreamUrl{StreamUrl: &url},
	}

	jsonItem, _ := json.Marshal(live)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mockService.MockILive, id uuid.UUID) {
				r.EXPECT().GetLiveById(id).Return(live, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonItem) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockILive, id uuid.UUID) {
				r.EXPECT().GetLiveById(id).Return(live, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetLiveById + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockILive := mockService.NewMockILive(c)
			test.mockBehavior(mockILive, id)

			services := &service.Service{ILive: mockILive}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Get("/live/{id}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				uuID, _ := uuid.Parse(iD)
				handler.GetLiveById(w, r, uuID)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/live/" + id.String()
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}
