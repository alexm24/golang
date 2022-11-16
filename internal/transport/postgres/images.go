package postgres

import (
	"database/sql"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

const imageApi = "api/images/"

type ImagesPostgres struct {
	db *sqlx.DB
}

func NewImagesPostgres(db *sqlx.DB) *ImagesPostgres {
	return &ImagesPostgres{db}
}

func (i *ImagesPostgres) GetImageById(id types.UUID) ([]uint8, error) {
	var fileBytes []uint8
	query := fmt.Sprintf("SELECT file FROM %s WHERE id = $1", imagesTable)
	if err := i.db.Get(&fileBytes, query, id); err != nil {
		if sql.ErrNoRows == err {
			return fileBytes, nil
		}
		return nil, err
	}
	return fileBytes, nil
}

func (i *ImagesPostgres) DelImageById(id types.UUID) (api.SIdentifier, error) {
	var item api.SIdentifier

	tx := i.db.MustBegin()

	query := fmt.Sprintf(`UPDATE %s SET file = NULL WHERE id = $1`, imagesTable)
	if _, err := tx.Exec(query, id); err != nil {
		if e := tx.Rollback(); e != nil {
			return item, e
		}
		return item, err
	}

	query = fmt.Sprintf(`UPDATE %s SET previewurl = '' WHERE id = $1`, broadcastTable)
	if _, err := tx.Exec(query, id); err != nil {
		if e := tx.Rollback(); e != nil {
			return item, e
		}
		return item, err
	}

	if e := tx.Commit(); e != nil {
		return item, e
	}

	item.Id = &id

	return item, nil
}

func (i *ImagesPostgres) CreateImage(id string, file []byte) (models.ResImage, error) {
	var item models.ResImage

	previewUrl := fmt.Sprintf("%s%s", imageApi, id)

	tx := i.db.MustBegin()

	query := fmt.Sprintf(`UPDATE %s SET file = $1 WHERE id = $2`, imagesTable)
	if _, err := tx.Exec(query, file, id); err != nil {
		if e := tx.Rollback(); e != nil {
			return item, e
		}
		return item, err
	}

	query = fmt.Sprintf(`UPDATE %s SET previewurl = $1 WHERE id = $2`, broadcastTable)
	if _, err := tx.Exec(query, previewUrl, id); err != nil {
		if e := tx.Rollback(); e != nil {
			return item, e
		}
		return item, err
	}

	if e := tx.Commit(); e != nil {
		return item, e
	}

	//item.Id = &id
	item.PreviewUrl = &previewUrl

	return item, nil
}
