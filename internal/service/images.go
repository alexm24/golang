package service

import (
	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/transport"
	"github.com/deepmap/oapi-codegen/pkg/types"
)

type ImagesService struct {
	imagesPostgres transport.IImagesPostgres
}

func NewImagesService(imagesPostgres transport.IImagesPostgres) *ImagesService {
	return &ImagesService{imagesPostgres}
}

func (i *ImagesService) GetImageById(id types.UUID) ([]uint8, error) {
	return i.imagesPostgres.GetImageById(id)
}

func (i *ImagesService) DelImageById(id types.UUID) (api.SIdentifier, error) {
	return i.imagesPostgres.DelImageById(id)
}

func (i *ImagesService) CreateImage(id string, file []byte) (models.ResImage, error) {
	return i.imagesPostgres.CreateImage(id, file)
}
