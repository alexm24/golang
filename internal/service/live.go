package service

import (
	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

type LiveService struct {
	livePostgres transport.ILivePostgres
}

func NewLiveService(livePostgres transport.ILivePostgres) *LiveService {
	return &LiveService{livePostgres}
}

func (l *LiveService) GetLive() ([]models.Live, error) {
	return l.livePostgres.GetLive()
}

func (l *LiveService) GetLiveById(id types.UUID) (models.Live, error) {
	return l.livePostgres.GetLiveById(id)
}
