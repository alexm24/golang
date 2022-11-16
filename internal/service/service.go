package service

import (
	"github.com/deepmap/oapi-codegen/pkg/types"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type IAdmin interface {
	CheckAdminUser(username api.SUsername) (api.SAdmin, error)
	GetToken(user api.SUsername) (api.SToken, error)
}

type IBroadcasts interface {
	CreateBroadcast(item models.PostBroadcast) (models.Broadcasts, error)
	ChangeBroadcast(item models.PutBroadcast) (models.Broadcasts, error)
	GetBroadcasts() ([]models.Broadcasts, error)
	GetBroadcastById(id types.UUID) (models.Broadcasts, error)
	DeleteBroadcast(id types.UUID) (api.SIdentifier, error)
	GetArchBroadcasts(username api.SUsername) ([]models.Broadcasts, error)
}

type IParticipants interface {
	CreateParticipant(channel string, user models.PostParticipant) error
	GetParticipants(channel string) ([]models.Participant, error)
}

type IMessages interface {
	GetMessageByChannel(channel string) ([]models.Messages, error)
	CreateMsg(channel string, msg models.PostMessage) (models.Messages, error)
	CreateReaction(channel string, item models.PostReactionMsg) error
	DeleteReaction(channel string, item models.PatchReactionMsg) error
}

type IStream interface {
	CreateStream(username api.SUsername) (models.Stream, error)
	GetStream(username string) (models.Stream, error)
	ChangeDescByUsername(stream models.PutStream) (models.Stream, error)
	ClearChat(channel string) error
}

type ILive interface {
	GetLive() ([]models.Live, error)
	GetLiveById(id types.UUID) (models.Live, error)
}

type IImages interface {
	CreateImage(id string, file []byte) (models.ResImage, error)
	GetImageById(id types.UUID) ([]uint8, error)
	DelImageById(id types.UUID) (api.SIdentifier, error)
}

type IZoom interface {
	GetObjectZoom(json string) (models.Zoom, error)
	SaveZoom(item models.Zoom) (models.Zoom, error)
	GetZoomByEmail(email string) ([]models.Zoom, error)
	GetZoomById(id types.UUID) (models.Zoom, error)
	SendMail(item models.Zoom) error
}

type Service struct {
	IAdmin
	IBroadcasts
	IParticipants
	IMessages
	IStream
	ILive
	IImages
	IZoom
}

func NewService(t *transport.Transport) *Service {
	return &Service{
		IAdmin:        NewAdminService(t.IBroadcastsPostgres, t.ICentrifugo),
		IBroadcasts:   NewBroadcastsService(t.IBroadcastsPostgres, t.IMessagesPostgres),
		IParticipants: NewParticipantsService(t.IParticipantsPostgres, t.IParticipantsRedis),
		IMessages:     NewMessagesService(t.IMessagesPostgres, t.ICentrifugo),
		IStream:       NewStreamService(t.IStreamPostgres, t.IMessagesPostgres, t.ICentrifugo),
		ILive:         NewLiveService(t.ILivePostgres),
		IImages:       NewImagesService(t.IImagesPostgres),
		IZoom:         NewZoomService(t.IZoomPostgres, t.IMail),
	}
}
