package models

import "github.com/alexm24/golang/internal/handler/api"

type Zoom struct {
	api.SIdentifier
	api.SStartTime
	api.SEMail
	api.SZoom
	api.SJson
}
