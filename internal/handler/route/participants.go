package route

import (
	"encoding/json"
	"net/http"

	"github.com/alexm24/golang/internal/models"
)

func (c *Route) PostParticipantsByChannel(w http.ResponseWriter, r *http.Request, channel string) {
	var item models.PostParticipant
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := item.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	err := c.service.IParticipants.CreateParticipant(channel, item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceCreateParticipant)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Route) GetParticipantsByChannel(w http.ResponseWriter, _ *http.Request, channel string) {
	items, err := c.service.IParticipants.GetParticipants(channel)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetParticipants)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}
