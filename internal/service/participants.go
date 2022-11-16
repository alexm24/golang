package service

import (
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

type ParticipantsService struct {
	participantsPostgres transport.IParticipantsPostgres
	participantsRedis    transport.IParticipantsRedis
}

func NewParticipantsService(
	participantsPostgres transport.IParticipantsPostgres,
	participantsRedis transport.IParticipantsRedis) *ParticipantsService {
	return &ParticipantsService{participantsPostgres, participantsRedis}
}

func (p *ParticipantsService) GetParticipants(channel string) ([]models.Participant, error) {
	return p.participantsPostgres.GetParticipants(channel)
}

func (p *ParticipantsService) CreateParticipant(channel string, user models.PostParticipant) error {
	return p.participantsRedis.CreateParticipant(channel, user)
}
