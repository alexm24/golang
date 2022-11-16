package models

import "github.com/alexm24/golang/internal/handler/api"

type Live struct {
	api.SIdentifier
	api.SDescription
	api.SPlace
	api.SStreamUrl
}
