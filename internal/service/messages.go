package service

import (
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

type MessagesService struct {
	messagesPostgres   transport.IMessagesPostgres
	messagesCentrifugo transport.ICentrifugo
}

func NewMessagesService(messagesPostgres transport.IMessagesPostgres, messagesCentrifugo transport.ICentrifugo) *MessagesService {
	return &MessagesService{messagesPostgres, messagesCentrifugo}
}

func (m *MessagesService) GetMessageByChannel(channel string) ([]models.Messages, error) {
	return m.messagesPostgres.GetMessageByChannel(channel)
}

func (m *MessagesService) CreateMsg(channel string, msg models.PostMessage) (models.Messages, error) {
	message, err := m.messagesPostgres.CreateMsg(channel, msg)
	if err != nil {
		return message, err
	}

	err = m.messagesCentrifugo.Publish(channel, message)

	return message, err
}

func (m *MessagesService) CreateReaction(channel string, item models.PostReactionMsg) error {
	message, err := m.messagesPostgres.AddReaction(item)
	if err != nil {
		return err
	}

	msg := models.ActionCentrifugo{Type: models.ActionChatReactions, Payload: message}

	err = m.messagesCentrifugo.Publish(channel, msg)

	return err
}

func (m *MessagesService) DeleteReaction(channel string, item models.PatchReactionMsg) error {
	message, err := m.messagesPostgres.DeleteReaction(item)
	if err != nil {
		return err
	}

	msg := models.ActionCentrifugo{Type: models.ActionChatReactions, Payload: message}

	err = m.messagesCentrifugo.Publish(channel, msg)
	return err
}
