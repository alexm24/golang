package service

import (
	"errors"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/tidwall/gjson"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

type ZoomService struct {
	zoomPostgres transport.IZoomPostgres
	mail         transport.IMail
}

func NewZoomService(zoomPostgres transport.IZoomPostgres, mail transport.IMail) *ZoomService {
	return &ZoomService{zoomPostgres, mail}
}

func (z *ZoomService) GetObjectZoom(json string) (models.Zoom, error) {
	var item models.Zoom

	recCount := gjson.Get(json, "payload.object.recording_count").Int()
	email := gjson.Get(json, "payload.object.host_email").String()
	if len(email) == 0 {
		return item, errors.New(msgZoomEmailEmpty)
	}
	topic := gjson.Get(json, "payload.object.topic").String()
	if len(topic) == 0 {
		return item, errors.New(msgZoomTopicEmpty)
	}
	startTime := gjson.Get(json, "payload.object.start_time").String()
	date, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return item, err
	}

	item = models.Zoom{
		SStartTime: api.SStartTime{StartTime: &date},
		SEMail:     api.SEMail{Email: &email},
		SZoom:      api.SZoom{RecordingCount: &recCount, Topic: &topic},
		SJson:      api.SJson{Json: &json},
	}
	return item, nil
}

func (z *ZoomService) SaveZoom(item models.Zoom) (models.Zoom, error) {
	return z.zoomPostgres.SaveZoom(item)
}

func (z *ZoomService) GetZoomByEmail(email string) ([]models.Zoom, error) {
	return z.zoomPostgres.GetZoomByEmail(email)
}

func (z *ZoomService) GetZoomById(id types.UUID) (models.Zoom, error) {
	return z.zoomPostgres.GetZoomById(id)
}

func (z *ZoomService) SendMail(item models.Zoom) error {
	return z.mail.SendMail(item)
}
