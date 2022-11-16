package route

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/alexm24/golang/internal/models"
)

func (c *Route) GetImageById(w http.ResponseWriter, _ *http.Request, id types.UUID) {
	item, err := c.service.IImages.GetImageById(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceGetImageById)
		return
	}

	if len(item) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	detectedFileType := http.DetectContentType(item)
	switch detectedFileType {
	case "image/jpeg", "image/jpg":
	case "image/png":
	default:
		newErrorResponse(w, http.StatusInternalServerError, models.MsgInvalidFileType, models.MsgInvalidFileType)
		return
	}

	buf := bytes.NewBuffer(item)

	w.Header().Set("Content-Type", detectedFileType)
	_, _ = io.Copy(w, buf)
}

func (c *Route) PutImageById(w http.ResponseWriter, _ *http.Request, id types.UUID) {
	item, err := c.service.IImages.DelImageById(id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceDelImageById)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}

func (c *Route) PostImage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), "ParseMultipartForm")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.MsgNoSuchFile)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), "ReadAll")
		return
	}

	id := r.FormValue("id")
	if len(id) == 0 {
		newErrorResponse(w, http.StatusBadRequest, models.MsgIdEmpty, models.MsgIdEmpty)
		return
	}

	item, err := c.service.IImages.CreateImage(id, fileBytes)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error(), models.ErrServiceCreateImage)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(item)
}
