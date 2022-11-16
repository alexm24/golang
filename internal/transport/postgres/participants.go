package postgres

import (
	"fmt"

	"github.com/alexm24/golang/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	participantsTable = "participants"
)

type ParticipantsPostgres struct {
	db *sqlx.DB
}

func NewParticipantsPostgres(db *sqlx.DB) *ParticipantsPostgres {
	return &ParticipantsPostgres{db}
}

func (p *ParticipantsPostgres) GetParticipants(channel string) ([]models.Participant, error) {
	var items = make([]models.Participant, 0)

	query := fmt.Sprintf(
		`SELECT username, fullname, email FROM %s WHERE channel ='%s';`,
		participantsTable, channel)

	if err := p.db.Select(&items, query); err != nil {
		return items, err
	}

	return items, nil
}
