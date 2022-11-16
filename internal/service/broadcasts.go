package service

import (
	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

type BroadcastsService struct {
	broadcastsPostgres transport.IBroadcastsPostgres
	messagesPostgres   transport.IMessagesPostgres
}

func NewBroadcastsService(broadcastsPostgres transport.IBroadcastsPostgres,
	messagesPostgres transport.IMessagesPostgres) *BroadcastsService {
	return &BroadcastsService{broadcastsPostgres, messagesPostgres}
}

func (b *BroadcastsService) GetBroadcasts() ([]models.Broadcasts, error) {
	return b.broadcastsPostgres.GetBroadcasts()
}

func (b *BroadcastsService) GetBroadcastById(id types.UUID) (models.Broadcasts, error) {
	return b.broadcastsPostgres.GetBroadcastById(id)
}

func (b *BroadcastsService) GetArchBroadcasts(username api.SUsername) ([]models.Broadcasts, error) {
	isAdmin, err := b.broadcastsPostgres.CheckAdminUser(username)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		return b.broadcastsPostgres.GetPastBroadcastAll()
	}

	return b.broadcastsPostgres.GetPastBroadcastOwner(username)
}

func (b *BroadcastsService) CreateBroadcast(item models.PostBroadcast) (models.Broadcasts, error) {
	return b.broadcastsPostgres.CreateBroadcast(item)
}

func (b *BroadcastsService) DeleteBroadcast(id types.UUID) (api.SIdentifier, error) {
	item, err := b.broadcastsPostgres.DeleteBroadcast(id)
	if err != nil {
		return item, err
	}

	channel := id.String()
	err = b.messagesPostgres.DeleteMessages(channel)
	if err != nil {
		return item, err
	}

	return item, nil
}

func (b *BroadcastsService) ChangeBroadcast(item models.PutBroadcast) (models.Broadcasts, error) {
	return b.broadcastsPostgres.ChangeBroadcast(item)
}
