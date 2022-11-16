package service

import (
	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

type StreamService struct {
	streamPostgres   transport.IStreamPostgres
	messagesPostgres transport.IMessagesPostgres
	centrifugo       transport.ICentrifugo
}

func NewStreamService(
	streamPostgres transport.IStreamPostgres,
	messagesPostgres transport.IMessagesPostgres,
	centrifugo transport.ICentrifugo) *StreamService {
	return &StreamService{streamPostgres, messagesPostgres, centrifugo}
}

func (s *StreamService) CreateStream(username api.SUsername) (models.Stream, error) {
	return s.streamPostgres.CreateStream(username)
}

func (s *StreamService) GetStream(username string) (models.Stream, error) {
	return s.streamPostgres.GetStream(username)
}

func (s *StreamService) ChangeDescByUsername(item models.PutStream) (models.Stream, error) {
	return s.streamPostgres.ChangeDescByUsername(item)
}

func (s *StreamService) ClearChat(channel string) error {
	err := s.messagesPostgres.DeleteMessages(channel)
	if err != nil {
		return err
	}
	msg := models.ActionCentrifugo{Type: models.ActionChatClear, Payload: ""}
	err = s.centrifugo.Publish(channel, msg)
	if err != nil {
		return err
	}
	return nil
}
