package route

import (
	"encoding/json"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/alexm24/golang/internal/models"
)

func (c *Route) GetLive(w http.ResponseWriter, _ *http.Request) {
	items, err := c.service.ILive.GetLive()
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetLive)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}

func (c *Route) GetLiveById(w http.ResponseWriter, _ *http.Request, id types.UUID) {
	items, err := c.service.ILive.GetLiveById(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetLiveById)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(items)
}
