package models

import (
	"errors"

	"github.com/alexm24/golang/internal/handler/api"
)

type Participant struct {
	api.SIdentifier
	api.SUsername
	api.SFullname
	api.SEMail
}

type PostParticipant api.PostParticipantsByChannelJSONBody

func (u *PostParticipant) Validate() error {
	if u.Fullname == nil {
		return errors.New(MsgFullnameEmpty)
	}
	if u.Username == nil {
		return errors.New(MsgUsernameEmpty)
	}
	if u.Email == nil {
		return errors.New(MsgEmailEmpty)
	}
	return nil
}
