package transport

import (
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport/centrifugo"
	"github.com/alexm24/golang/internal/transport/mail"
	"github.com/alexm24/golang/internal/transport/postgres"
	redisPool "github.com/alexm24/golang/internal/transport/redis"
)

type IBroadcastsPostgres interface {
	CreateBroadcast(item models.PostBroadcast) (models.Broadcasts, error)
	ChangeBroadcast(item models.PutBroadcast) (models.Broadcasts, error)
	GetBroadcasts() ([]models.Broadcasts, error)
	GetBroadcastById(id types.UUID) (models.Broadcasts, error)
	DeleteBroadcast(id types.UUID) (api.SIdentifier, error)
	CheckAdminUser(username api.SUsername) (bool, error)
	GetPastBroadcastAll() ([]models.Broadcasts, error)
	GetPastBroadcastOwner(user api.SUsername) ([]models.Broadcasts, error)
}

type IParticipantsPostgres interface {
	GetParticipants(channel string) ([]models.Participant, error)
}

type IParticipantsRedis interface {
	CreateParticipant(channel string, user models.PostParticipant) error
}

type IMessagesPostgres interface {
	GetMessageByChannel(channel string) ([]models.Messages, error)
	CreateMsg(channel string, msg models.PostMessage) (models.Messages, error)
	DeleteMessages(channel string) error
	AddReaction(item models.PostReactionMsg) (models.Messages, error)
	DeleteReaction(item models.PatchReactionMsg) (models.Messages, error)
}

type IStreamPostgres interface {
	CreateStream(username api.SUsername) (models.Stream, error)
	GetStream(username string) (models.Stream, error)
	ChangeDescByUsername(stream models.PutStream) (models.Stream, error)
}

type ICentrifugo interface {
	GetToken(username api.SUsername) (api.SToken, error)
	Publish(channel string, msg interface{}) error
}

type ILivePostgres interface {
	GetLive() ([]models.Live, error)
	GetLiveById(id types.UUID) (models.Live, error)
}

type IImagesPostgres interface {
	CreateImage(id string, file []byte) (models.ResImage, error)
	GetImageById(id types.UUID) ([]uint8, error)
	DelImageById(id types.UUID) (sid api.SIdentifier, err error)
}

type IZoomPostgres interface {
	SaveZoom(item models.Zoom) (models.Zoom, error)
	GetZoomByEmail(email string) ([]models.Zoom, error)
	GetZoomById(id types.UUID) (models.Zoom, error)
}

type IMail interface {
	SendMail(item models.Zoom) error
}

type Transport struct {
	IParticipantsRedis
	IBroadcastsPostgres
	IParticipantsPostgres
	IMessagesPostgres
	IStreamPostgres
	ILivePostgres
	ICentrifugo
	IImagesPostgres
	IZoomPostgres
	IMail
}

func NewTransport(db *sqlx.DB, rp *redis.Pool, c models.CentrifugoConfig) *Transport {
	return &Transport{
		IParticipantsRedis:    redisPool.NewParticipantsRedis(rp),
		IBroadcastsPostgres:   postgres.NewBroadcastsPostgres(db),
		IParticipantsPostgres: postgres.NewParticipantsPostgres(db),
		IMessagesPostgres:     postgres.NewMessagesPostgres(db),
		IStreamPostgres:       postgres.NewStreamPostgres(db),
		ILivePostgres:         postgres.NewLivePostgres(db),
		ICentrifugo:           centrifugo.NewCentrifugo(c),
		IImagesPostgres:       postgres.NewImagesPostgres(db),
		IZoomPostgres:         postgres.NewZoomPostgres(db),
		IMail:                 mail.NewMail(),
	}
}
