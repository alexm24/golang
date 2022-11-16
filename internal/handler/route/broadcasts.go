package route

import (
	"encoding/json"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

func (c *Route) GetBroadcasts(w http.ResponseWriter, _ *http.Request) {
	items, err := c.service.IBroadcasts.GetBroadcasts()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetBroadcasts)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (c *Route) GetBroadcastById(w http.ResponseWriter, _ *http.Request, id types.UUID) {
	item, err := c.service.IBroadcasts.GetBroadcastById(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetBroadcastById)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (c *Route) PostBroadcasts(w http.ResponseWriter, r *http.Request) {
	var item models.PostBroadcast
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := item.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	broadcast, err := c.service.IBroadcasts.CreateBroadcast(item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceCreateBroadcast)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(broadcast)
}

func (c *Route) PostUserGetBroadcastArch(w http.ResponseWriter, r *http.Request) {
	var user api.SUsername
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}
	if user.Username == nil {
		newErrorResponse(w, http.StatusBadRequest, "username field is empty", models.MsgUsernameEmpty)
		return
	}

	items, err := c.service.IBroadcasts.GetArchBroadcasts(user)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetArchBroadcasts)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(items)
}

func (c *Route) DeleteBroadcast(w http.ResponseWriter, _ *http.Request, id types.UUID) {
	item, err := c.service.IBroadcasts.DeleteBroadcast(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceDeleteBroadcast)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (c *Route) PutBroadcast(w http.ResponseWriter, r *http.Request) {
	var item models.PutBroadcast
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := item.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	broadcast, err := c.service.IBroadcasts.ChangeBroadcast(item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceChangeBroadcast)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(broadcast)
}
