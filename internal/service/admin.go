package service

import (
	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/transport"
)

type AdminService struct {
	broadcastsPostgres transport.IBroadcastsPostgres
	centrifugo         transport.ICentrifugo
}

func NewAdminService(broadcastsPostgres transport.IBroadcastsPostgres, centrifugo transport.ICentrifugo) *AdminService {
	return &AdminService{broadcastsPostgres, centrifugo}
}

func (a *AdminService) CheckAdminUser(username api.SUsername) (api.SAdmin, error) {
	isAdmin, err := a.broadcastsPostgres.CheckAdminUser(username)
	return api.SAdmin{IsAdmin: &isAdmin}, err
}

func (a *AdminService) GetToken(username api.SUsername) (api.SToken, error) {
	return a.centrifugo.GetToken(username)
}
