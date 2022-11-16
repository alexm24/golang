package route

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alexm24/golang/internal/models"
	"github.com/deepmap/oapi-codegen/pkg/types"
)

func (c *Route) PostZoom(w http.ResponseWriter, r *http.Request) {
	//TODO check json
	b, _ := io.ReadAll(r.Body)
	str := string(b)

	zoom, err := c.service.IZoom.GetObjectZoom(str)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	item, err := c.service.IZoom.SaveZoom(zoom)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	err = c.service.IZoom.SendMail(item)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Route) GetZoomByEmail(w http.ResponseWriter, _ *http.Request, email string) {
	item, err := c.service.IZoom.GetZoomByEmail(email)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetZoomByEmail)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (c *Route) GetZoomById(w http.ResponseWriter, _ *http.Request, id types.UUID) {
	item, err := c.service.IZoom.GetZoomById(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetZoomById)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}
