package route

import (
	"encoding/json"
	"net/http"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

func (c *Route) GetStreamByUsername(w http.ResponseWriter, _ *http.Request, username string) {
	stream, err := c.service.IStream.GetStream(username)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetStream)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stream)
}

func (c *Route) PostStream(w http.ResponseWriter, r *http.Request) {
	var user api.SUsername
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if user.Username == nil {
		newErrorResponse(w, http.StatusBadRequest, models.MsgUsernameEmpty, models.MsgUsernameEmpty)
		return
	}

	stream, err := c.service.IStream.CreateStream(user)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceCreateStream)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(stream)
}

func (c *Route) PutStream(w http.ResponseWriter, r *http.Request) {
	var item models.PutStream
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := item.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	res, err := c.service.IStream.ChangeDescByUsername(item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceChangeDescByUsername)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (c *Route) DeleteStreamChat(w http.ResponseWriter, _ *http.Request, channel string) {
	err := c.service.IStream.ClearChat(channel)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceClearChat)
		return
	}
	w.WriteHeader(http.StatusOK)
}
