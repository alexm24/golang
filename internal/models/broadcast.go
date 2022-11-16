package models

import (
	"errors"

	"github.com/alexm24/golang/internal/handler/api"
)

type LifeCycleBroadcast int

func (l LifeCycleBroadcast) String() string {
	return [...]string{"created", "past"}[l]
}

const (
	Created LifeCycleBroadcast = iota
	Past
)

type Broadcasts struct {
	api.SIdentifier
	api.SPreviewUrl
	api.SBroadcast
	api.SLifeCycle
	api.SStartTime
}

type PutBroadcast api.PutBroadcastJSONBody

func (p *PutBroadcast) Validate() error {
	if p.Id == nil {
		return errors.New(MsgIdEmpty)
	}
	if p.Name == nil {
		return errors.New(MsgNameEmpty)
	}
	if p.Description == nil {
		return errors.New(MsgDescriptionEmpty)
	}
	if p.Owner == nil {
		return errors.New(MsgOwnerEmpty)
	}
	if p.StreamKey == nil {
		return errors.New(MsgStreamKeyEmpty)
	}
	if p.StartTime == nil {
		return errors.New(MsgStartTimeEmpty)
	}
	return nil
}

type PostBroadcast api.PostBroadcastsJSONBody

func (p *PostBroadcast) Validate() error {
	if p.Name == nil {
		return errors.New(MsgNameEmpty)
	}
	if p.Description == nil {
		return errors.New(MsgDescriptionEmpty)
	}
	if p.Owner == nil {
		return errors.New(MsgOwnerEmpty)
	}
	if p.StreamKey == nil {
		return errors.New(MsgStreamKeyEmpty)
	}
	if p.StartTime == nil {
		return errors.New(MsgStartTimeEmpty)
	}
	return nil
}
