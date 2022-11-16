package models

import (
	"errors"

	"github.com/alexm24/golang/internal/handler/api"
)

type Messages struct {
	api.SIdentifier
	api.SFullname
	api.SUsername
	api.SMessage
	Channel string `json:"-" db:"channel"`
}

type PostMessage api.PostMsgByChannelJSONBody

func (p *PostMessage) Validate() error {
	if p.Username == nil {
		return errors.New(MsgUsernameEmpty)
	}
	if p.Fullname == nil {
		return errors.New(MsgFullnameEmpty)
	}
	if p.Text == nil {
		return errors.New(MsgTextEmpty)
	}
	if p.Avatar == nil {
		return errors.New(MsgAvatarEmpty)
	}
	if p.Time == nil {
		return errors.New(MsgTimeEmpty)
	}
	if p.IsAnon == nil {
		return errors.New(MsgIsAnonEmpty)
	}
	if p.IsQuestion == nil {
		return errors.New(MsgIsQuestionEmpty)
	}
	return nil
}

type PostReactionMsg api.PostReactionMsgJSONBody

func (p *PostReactionMsg) Validate() error {
	if p.Id == nil {
		return errors.New(MsgIdEmpty)
	}
	if p.Username == nil {
		return errors.New(MsgUsernameEmpty)
	}
	if p.Type == nil {
		return errors.New(MsgTypeEmpty)
	}
	return nil
}

type PatchReactionMsg api.PatchReactionMsgJSONBody

func (p *PatchReactionMsg) Validate() error {
	if p.Id == nil {
		return errors.New(MsgIdEmpty)
	}
	if p.Username == nil {
		return errors.New(MsgUsernameEmpty)
	}
	return nil
}
