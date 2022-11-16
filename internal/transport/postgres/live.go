package postgres

import (
	"database/sql"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/models"
)

const (
	liveTable = "live"
)

type LivePostgres struct {
	db *sqlx.DB
}

func NewLivePostgres(db *sqlx.DB) *LivePostgres {
	return &LivePostgres{db}
}

func (l *LivePostgres) GetLive() ([]models.Live, error) {
	var live []models.Live
	query := fmt.Sprintf("SELECT * FROM %s", liveTable)
	if err := l.db.Select(&live, query); err != nil {
		return live, err
	}
	return live, nil
}

func (l *LivePostgres) GetLiveById(id types.UUID) (models.Live, error) {
	var live models.Live
	query := fmt.Sprintf("SELECT * FROM %s WHERE id::text = $1", liveTable)
	if err := l.db.Get(&live, query, id); err != nil {
		if err == sql.ErrNoRows {
			return live, nil
		}
		return live, err
	}
	return live, nil
}
