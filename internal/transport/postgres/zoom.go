package postgres

import (
	"database/sql"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/models"
)

const (
	zoomTable = "zoom"
)

type ZoomPostgres struct {
	db *sqlx.DB
}

func NewZoomPostgres(db *sqlx.DB) *ZoomPostgres {
	return &ZoomPostgres{db}
}

func (z *ZoomPostgres) SaveZoom(item models.Zoom) (models.Zoom, error) {
	var i models.Zoom
	q := fmt.Sprintf(
		`INSERT INTO %s (id, start_time, email, topic, json, recording_count) values (uuid_generate_v4(), $1, $2, $3, $4, $5)
				RETURNING id, topic, email;`,
		zoomTable)
	row := z.db.QueryRowx(q, item.StartTime, item.Email, item.Topic, item.Json, item.RecordingCount)
	if err := row.StructScan(&i); err != nil {
		return i, err
	}
	return i, nil
}

func (z *ZoomPostgres) GetZoomByEmail(email string) ([]models.Zoom, error) {
	item := make([]models.Zoom, 0)
	q := fmt.Sprintf(
		"SELECT id, start_time, topic, recording_count FROM %s WHERE LOWER(email)=LOWER($1) ORDER BY start_time DESC",
		zoomTable)
	if err := z.db.Select(&item, q, email); err != nil {
		return item, err
	}
	return item, nil
}

func (z *ZoomPostgres) GetZoomById(id types.UUID) (models.Zoom, error) {
	var item models.Zoom
	q := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 ", zoomTable)
	err := z.db.Get(&item, q, id)
	if err != nil {
		if sql.ErrNoRows == err {
			return item, nil
		}
		return item, err
	}
	return item, nil
}
