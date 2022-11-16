package route

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexm24/golang/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/service"
	mockService "github.com/alexm24/golang/internal/service/mocks"
)

func TestRoute_GetImageById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIImages, id uuid.UUID)

	data := []uint8{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 1, 0, 96, 0, 96, 0, 0, 255, 219, 0, 67, 0, 2, 1, 1, 2, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 3, 5, 3, 3, 3, 3, 3, 6, 4, 4, 3, 5, 7, 6, 7, 7, 7, 6, 7, 7, 8, 9, 11, 9, 8, 8, 10, 8, 7, 7, 10, 13, 10, 10, 11, 12, 12, 12, 12, 7, 9, 14, 15, 13, 12, 14, 11, 12, 12, 12, 255, 219, 0, 67, 1, 2, 2, 2, 3, 3, 3, 6, 3, 3, 6, 12, 8, 7, 8, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 255, 192, 0, 17, 8, 0, 1, 0, 2, 3, 1, 34, 0, 2, 17, 1, 3, 17, 1, 255, 196, 0, 31, 0, 0, 1, 5, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 255, 196, 0, 181, 16, 0, 2, 1, 3, 3, 2, 4, 3, 5, 5, 4, 4, 0, 0, 1, 125, 1, 2, 3, 0, 4, 17, 5, 18, 33, 49, 65, 6, 19, 81, 97, 7, 34, 113, 20, 50, 129, 145, 161, 8, 35, 66, 177, 193, 21, 82, 209, 240, 36, 51, 98, 114, 130, 9, 10, 22, 23, 24, 25, 26, 37, 38, 39, 40, 41, 42, 52, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73, 74, 83, 84, 85, 86, 87, 88, 89, 90, 99, 100, 101, 102, 103, 104, 105, 106, 115, 116, 117, 118, 119, 120, 121, 122, 131, 132, 133, 134, 135, 136, 137, 138, 146, 147, 148, 149, 150, 151, 152, 153, 154, 162, 163, 164, 165, 166, 167, 168, 169, 170, 178, 179, 180, 181, 182, 183, 184, 185, 186, 194, 195, 196, 197, 198, 199, 200, 201, 202, 210, 211, 212, 213, 214, 215, 216, 217, 218, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 255, 196, 0, 31, 1, 0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 255, 196, 0, 181, 17, 0, 2, 1, 2, 4, 4, 3, 4, 7, 5, 4, 4, 0, 1, 2, 119, 0, 1, 2, 3, 17, 4, 5, 33, 49, 6, 18, 65, 81, 7, 97, 113, 19, 34, 50, 129, 8, 20, 66, 145, 161, 177, 193, 9, 35, 51, 82, 240, 21, 98, 114, 209, 10, 22, 36, 52, 225, 37, 241, 23, 24, 25, 26, 38, 39, 40, 41, 42, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73, 74, 83, 84, 85, 86, 87, 88, 89, 90, 99, 100, 101, 102, 103, 104, 105, 106, 115, 116, 117, 118, 119, 120, 121, 122, 130, 131, 132, 133, 134, 135, 136, 137, 138, 146, 147, 148, 149, 150, 151, 152, 153, 154, 162, 163, 164, 165, 166, 167, 168, 169, 170, 178, 179, 180, 181, 182, 183, 184, 185, 186, 194, 195, 196, 197, 198, 199, 200, 201, 202, 210, 211, 212, 213, 214, 215, 216, 217, 218, 226, 227, 228, 229, 230, 231, 232, 233, 234, 242, 243, 244, 245, 246, 247, 248, 249, 250, 255, 218, 0, 12, 3, 1, 0, 2, 17, 3, 17, 0, 63, 0, 254, 127, 232, 162, 138, 0, 255, 217}
	dataErr := []uint8{255, 13, 11}
	id := uuid.New()

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mockService.MockIImages, id uuid.UUID) {
				r.EXPECT().GetImageById(id).Return(data, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(data),
		},
		{
			name: "Invalid file type",
			mockBehavior: func(r *mockService.MockIImages, id uuid.UUID) {
				r.EXPECT().GetImageById(id).Return(dataErr, nil)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.MsgInvalidFileType + `"}` + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIImages, id uuid.UUID) {
				r.EXPECT().GetImageById(id).Return(data, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetImageById + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIImages := mockService.NewMockIImages(c)
			test.mockBehavior(mockIImages, id)

			services := &service.Service{IImages: mockIImages}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Get("/images/{id}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				uuID, _ := uuid.Parse(iD)
				handler.GetImageById(w, r, uuID)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/images/" + id.String()
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}

func TestRoute_PutImageById(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIImages, params uuid.UUID)

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
			mockBehavior: func(r *mockService.MockIImages, i uuid.UUID) {
				r.EXPECT().DelImageById(i).Return(id, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonId) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIImages, p uuid.UUID) {
				r.EXPECT().DelImageById(p).Return(id, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceDelImageById + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIImages := mockService.NewMockIImages(c)
			test.mockBehavior(mockIImages, params)

			services := &service.Service{IImages: mockIImages}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Put("/images/{id}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				uuID, _ := uuid.Parse(iD)
				handler.PutImageById(w, r, uuID)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/images/" + params.String()
			req := httptest.NewRequest(http.MethodPut, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}

func TestRoute_PostImage(t *testing.T) {

}
