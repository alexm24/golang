package route

import (
	"encoding/json"
	"net/http"

	"github.com/alexm24/golang/internal/models"
)

func (c *Route) GetMsgByChannel(w http.ResponseWriter, _ *http.Request, channel string) {
	msg, err := c.service.IMessages.GetMessageByChannel(channel)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetMessageByChannel)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(msg)
}

func (c *Route) PostMsgByChannel(w http.ResponseWriter, r *http.Request, channel string) {
	var msgBody models.PostMessage
	if err := json.NewDecoder(r.Body).Decode(&msgBody); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := msgBody.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	msg, err := c.service.IMessages.CreateMsg(channel, msgBody)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceErrCreateMsg)
		return
	}

	_ = json.NewEncoder(w).Encode(msg)
	w.WriteHeader(http.StatusOK)
}

func (c *Route) PostReactionMsg(w http.ResponseWriter, r *http.Request, channel string) {
	var item models.PostReactionMsg
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := item.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	err := c.service.IMessages.CreateReaction(channel, item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceCreateReaction)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Route) PatchReactionMsg(w http.ResponseWriter, r *http.Request, channel string) {
	var item models.PatchReactionMsg
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), models.MsgInvalidJson)
		return
	}

	if err := item.Validate(); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	err := c.service.IMessages.DeleteReaction(channel, item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceDeleteReaction)
		return
	}

	w.WriteHeader(http.StatusOK)
}
