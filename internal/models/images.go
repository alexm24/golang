package models

import "github.com/alexm24/golang/internal/handler/api"

type ResImage struct {
	api.SIdentifier
	api.SPreviewUrl
}
