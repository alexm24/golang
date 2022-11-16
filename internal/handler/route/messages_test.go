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

func TestRoute_GetMsgByChannel(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIMessages, channel string)

	id := uuid.New()
	channel := "channel"
	user := "test"
	avatar := "e3b0c44298fc1c1494c8996fb92427ae41e4649b934ca495991b7852b855"
	fullname := "test test"
	isQuestion := false
	text := "messages"
	date := time.Now()
	isAnon := true
	reactions := "{}"

	arrayMsg := []models.Messages{
		{
			SIdentifier: api.SIdentifier{Id: &id},
			SUsername:   api.SUsername{Username: &user},
			SFullname:   api.SFullname{Fullname: &fullname},
			SMessage: api.SMessage{
				Avatar:     &avatar,
				IsQuestion: &isQuestion,
				Text:       &text,
				Time:       &date,
				IsAnon:     &isAnon,
				Reactions:  &reactions,
			},
		},
	}

	jsonArrayMsg, _ := json.Marshal(arrayMsg)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mockService.MockIMessages, channel string) {
				r.EXPECT().GetMessageByChannel(channel).Return(arrayMsg, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonArrayMsg) + "\n",
		},
		{
			name: "Service failure",
			mockBehavior: func(r *mockService.MockIMessages, channel string) {
				r.EXPECT().GetMessageByChannel(channel).Return(arrayMsg, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceGetMessageByChannel + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIMsg := mockService.NewMockIMessages(c)
			test.mockBehavior(mockIMsg, channel)

			services := &service.Service{IMessages: mockIMsg}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Get("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				handler.GetMsgByChannel(w, r, iD)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/messages/" + channel
			req := httptest.NewRequest(http.MethodGet, path, nil)

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}

func TestRoute_PostMsgByChannel(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIMessages, channel string, msg models.PostMessage)

	id := uuid.New()
	channel := "channel"
	user := "test"
	avatar := "e3b0c44298fc1c1494c8996fb92427ae41e4649b934ca495991b7852b855"
	fullname := "test test"
	isQuestion := false
	text := "messages"
	date := time.Date(2022, time.July, 1, 8, 0, 0, 0, time.UTC)
	isAnon := true

	postMsg := models.PostMessage{
		Avatar:     &avatar,
		Fullname:   &fullname,
		IsAnon:     &isAnon,
		IsQuestion: &isQuestion,
		Text:       &text,
		Time:       &date,
		Username:   &user,
	}

	postMsgWithoutUsername := models.PostMessage{
		Avatar:     &avatar,
		Fullname:   &fullname,
		IsAnon:     &isAnon,
		IsQuestion: &isQuestion,
		Text:       &text,
		Time:       &date,
	}
	postMsgWithoutFullname := models.PostMessage{
		Avatar:     &avatar,
		IsAnon:     &isAnon,
		IsQuestion: &isQuestion,
		Text:       &text,
		Time:       &date,
		Username:   &user,
	}
	postMsgWithoutText := models.PostMessage{
		Avatar:     &avatar,
		Fullname:   &fullname,
		IsAnon:     &isAnon,
		IsQuestion: &isQuestion,
		Time:       &date,
		Username:   &user,
	}
	postMsgWithoutAvatar := models.PostMessage{
		Fullname:   &fullname,
		IsAnon:     &isAnon,
		IsQuestion: &isQuestion,
		Text:       &text,
		Time:       &date,
		Username:   &user,
	}
	postMsgWithoutTime := models.PostMessage{
		Avatar:     &avatar,
		Fullname:   &fullname,
		IsAnon:     &isAnon,
		IsQuestion: &isQuestion,
		Text:       &text,
		Username:   &user,
	}
	postMsgWithoutIsAnon := models.PostMessage{
		Avatar:     &avatar,
		Fullname:   &fullname,
		IsQuestion: &isQuestion,
		Text:       &text,
		Time:       &date,
		Username:   &user,
	}
	postMsgWithoutIsQuestion := models.PostMessage{
		Avatar:   &avatar,
		Fullname: &fullname,
		IsAnon:   &isAnon,
		Text:     &text,
		Time:     &date,
		Username: &user,
	}

	resMsg := models.Messages{
		SIdentifier: api.SIdentifier{Id: &id},
		SUsername:   api.SUsername{Username: &user},
		SFullname:   api.SFullname{Fullname: &fullname},
		SMessage: api.SMessage{
			Avatar:     &avatar,
			IsQuestion: &isQuestion,
			Text:       &text,
			Time:       &date,
			IsAnon:     &isAnon,
		},
	}

	jsonPostMsg, _ := json.Marshal(postMsg)
	jsonPostMsgWithoutUsername, _ := json.Marshal(postMsgWithoutUsername)
	jsonPostMsgWithoutFullname, _ := json.Marshal(postMsgWithoutFullname)
	jsonPostMsgWithoutText, _ := json.Marshal(postMsgWithoutText)
	jsonPostMsgWithoutAvatar, _ := json.Marshal(postMsgWithoutAvatar)
	jsonPostMsgWithoutTime, _ := json.Marshal(postMsgWithoutTime)
	jsonPostMsgWithoutIsAnon, _ := json.Marshal(postMsgWithoutIsAnon)
	jsonPostMsgWithoutIsQuestion, _ := json.Marshal(postMsgWithoutIsQuestion)

	jsonResMsg, _ := json.Marshal(resMsg)

	tests := []struct {
		name                 string
		inputBody            string
		inputMsg             models.PostMessage
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: string(jsonPostMsg),
			inputMsg:  postMsg,
			mockBehavior: func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {
				r.EXPECT().CreateMsg(channel, postMsg).Return(resMsg, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(jsonResMsg) + "\n",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonPostMsg),
			inputMsg:  postMsg,
			mockBehavior: func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {
				r.EXPECT().CreateMsg(channel, postMsg).Return(resMsg, errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceErrCreateMsg + `"}` + "\n",
		},
		{
			name:                 "username field is empty",
			inputBody:            string(jsonPostMsgWithoutUsername),
			inputMsg:             postMsgWithoutUsername,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:                 "fullname field is empty",
			inputBody:            string(jsonPostMsgWithoutFullname),
			inputMsg:             postMsgWithoutFullname,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgFullnameEmpty + `"}` + "\n",
		},
		{
			name:                 "text field is empty",
			inputBody:            string(jsonPostMsgWithoutText),
			inputMsg:             postMsgWithoutText,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgTextEmpty + `"}` + "\n",
		},
		{
			name:                 "avatar field is empty",
			inputBody:            string(jsonPostMsgWithoutAvatar),
			inputMsg:             postMsgWithoutAvatar,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgAvatarEmpty + `"}` + "\n",
		},
		{
			name:                 "time field is empty",
			inputBody:            string(jsonPostMsgWithoutTime),
			inputMsg:             postMsgWithoutTime,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgTimeEmpty + `"}` + "\n",
		},
		{
			name:                 "is_anon field is empty",
			inputBody:            string(jsonPostMsgWithoutIsAnon),
			inputMsg:             postMsgWithoutIsAnon,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgIsAnonEmpty + `"}` + "\n",
		},
		{
			name:                 "is_question field is empty",
			inputBody:            string(jsonPostMsgWithoutIsQuestion),
			inputMsg:             postMsgWithoutIsQuestion,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, msg models.PostMessage) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgIsQuestionEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIMsg := mockService.NewMockIMessages(c)
			test.mockBehavior(mockIMsg, channel, postMsg)

			services := &service.Service{IMessages: mockIMsg}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Post("/messages/{id}", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				handler.PostMsgByChannel(w, r, iD)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/messages/" + channel
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}

func TestRoute_PostReactionMsg(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIMessages, channel string, item models.PostReactionMsg)

	id := uuid.New()
	username := "test"
	typ := "like"

	msgReaction := models.PostReactionMsg{
		Id:       &id,
		Username: &username,
		Type:     &typ,
	}
	msgReactionWithoutId := models.PostReactionMsg{
		Username: &username,
		Type:     &typ,
	}
	msgReactionWithoutUsername := models.PostReactionMsg{
		Id:   &id,
		Type: &typ,
	}
	msgReactionWithoutType := models.PostReactionMsg{
		Id:       &id,
		Username: &username,
	}

	jsonMsgReaction, _ := json.Marshal(msgReaction)
	jsonMsgReactionWithoutId, _ := json.Marshal(msgReactionWithoutId)
	jsonMsgReactionWithoutUsername, _ := json.Marshal(msgReactionWithoutUsername)
	jsonMsgReactionWithoutType, _ := json.Marshal(msgReactionWithoutType)

	tests := []struct {
		name                 string
		channel              string
		inputBody            string
		input                models.PostReactionMsg
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			channel:   "channel",
			inputBody: string(jsonMsgReaction),
			input:     msgReaction,
			mockBehavior: func(r *mockService.MockIMessages, channel string, item models.PostReactionMsg) {
				r.EXPECT().CreateReaction(channel, item).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonMsgReaction),
			input:     msgReaction,
			mockBehavior: func(r *mockService.MockIMessages, channel string, item models.PostReactionMsg) {
				r.EXPECT().CreateReaction(channel, item).Return(errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceCreateReaction + `"}` + "\n",
		},
		{
			name:                 models.MsgUsernameEmpty,
			inputBody:            string(jsonMsgReactionWithoutUsername),
			input:                msgReactionWithoutUsername,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, item models.PostReactionMsg) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:                 models.MsgIdEmpty,
			inputBody:            string(jsonMsgReactionWithoutId),
			input:                msgReactionWithoutId,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, item models.PostReactionMsg) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgIdEmpty + `"}` + "\n",
		},
		{
			name:                 models.MsgTypeEmpty,
			inputBody:            string(jsonMsgReactionWithoutType),
			input:                msgReactionWithoutType,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, item models.PostReactionMsg) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgTypeEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIMsg := mockService.NewMockIMessages(c)
			test.mockBehavior(mockIMsg, test.channel, test.input)

			services := &service.Service{IMessages: mockIMsg}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Post("/messages/{id}/reaction", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				handler.PostReactionMsg(w, r, iD)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/messages/" + test.channel + "/reaction"
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}

func TestRoute_PatchReactionMsg(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mockService.MockIMessages, channel string, item models.PatchReactionMsg)

	id := uuid.New()
	username := "test"

	msgReaction := models.PatchReactionMsg{
		Id:       &id,
		Username: &username,
	}
	msgReactionWithoutId := models.PatchReactionMsg{
		Username: &username,
	}
	msgReactionWithoutUsername := models.PatchReactionMsg{
		Id: &id,
	}

	jsonMsgReaction, _ := json.Marshal(msgReaction)
	jsonMsgReactionWithoutId, _ := json.Marshal(msgReactionWithoutId)
	jsonMsgReactionWithoutUsername, _ := json.Marshal(msgReactionWithoutUsername)

	tests := []struct {
		name                 string
		channel              string
		inputBody            string
		input                models.PatchReactionMsg
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			channel:   "channel",
			inputBody: string(jsonMsgReaction),
			input:     msgReaction,
			mockBehavior: func(r *mockService.MockIMessages, channel string, item models.PatchReactionMsg) {
				r.EXPECT().DeleteReaction(channel, item).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "",
		},
		{
			name:      "Service failure",
			inputBody: string(jsonMsgReaction),
			input:     msgReaction,
			mockBehavior: func(r *mockService.MockIMessages, channel string, item models.PatchReactionMsg) {
				r.EXPECT().DeleteReaction(channel, item).Return(errors.New("error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"code":` + "500" + `,"message":"` + models.ErrServiceDeleteReaction + `"}` + "\n",
		},
		{
			name:                 models.MsgUsernameEmpty,
			inputBody:            string(jsonMsgReactionWithoutUsername),
			input:                msgReactionWithoutUsername,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, item models.PatchReactionMsg) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgUsernameEmpty + `"}` + "\n",
		},
		{
			name:                 models.MsgIdEmpty,
			inputBody:            string(jsonMsgReactionWithoutId),
			input:                msgReactionWithoutId,
			mockBehavior:         func(r *mockService.MockIMessages, channel string, item models.PatchReactionMsg) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"code":` + "400" + `,"message":"` + models.MsgIdEmpty + `"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			mockIMsg := mockService.NewMockIMessages(c)
			test.mockBehavior(mockIMsg, test.channel, test.input)

			services := &service.Service{IMessages: mockIMsg}
			handler := Route{services}

			// Init Endpoint
			route := chi.NewRouter()
			route.Patch("/messages/{id}/reaction", func(w http.ResponseWriter, r *http.Request) {
				iD := chi.URLParam(r, "id")
				handler.PatchReactionMsg(w, r, iD)
			})

			// Create Request
			w := httptest.NewRecorder()
			path := "/messages/" + test.channel + "/reaction"
			req := httptest.NewRequest(http.MethodPatch, path, bytes.NewBufferString(test.inputBody))

			// Make Request
			route.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}

}
