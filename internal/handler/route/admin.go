package route

import (
	"encoding/json"
	"net/http"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

func (c *Route) CheckAdmin(w http.ResponseWriter, r *http.Request) {
	var username api.SUsername
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}
	if username.Username == nil {
		newErrorResponse(w, http.StatusBadRequest, models.MsgUsernameEmpty, models.MsgUsernameEmpty)
		return
	}

	item, err := c.service.IAdmin.CheckAdminUser(username)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceCheckAdminUser)
		return
	}

	_ = json.NewEncoder(w).Encode(item)
	w.WriteHeader(http.StatusOK)
}

func (c *Route) PostUserGetToken(w http.ResponseWriter, r *http.Request) {
	var username api.SUsername
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}
	if username.Username == nil {
		newErrorResponse(w, http.StatusBadRequest, models.MsgUsernameEmpty, models.MsgUsernameEmpty)
		return
	}

	item, err := c.service.IAdmin.GetToken(username)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetToken)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}
