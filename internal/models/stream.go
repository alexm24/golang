package models

import (
	"errors"

	"github.com/alexm24/golang/internal/handler/api"
)

type Stream struct {
	api.SIdentifier
	api.SUsername
	api.SDescription
}

type PutStream api.PutStreamJSONBody

func (p *PutStream) Validate() error {
	if p.Username == nil {
		return errors.New(MsgUsernameEmpty)
	}
	if p.Description == nil {
		return errors.New(MsgDescriptionEmpty)
	}
	return nil
}
