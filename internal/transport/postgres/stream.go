package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

const (
	streamTable = "stream"
)

type StreamPostgres struct {
	db *sqlx.DB
}

func NewStreamPostgres(db *sqlx.DB) *StreamPostgres {
	return &StreamPostgres{db}
}

func (s *StreamPostgres) CreateStream(user api.SUsername) (models.Stream, error) {
	var stream models.Stream

	qSelect := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1`, streamTable)
	qInsert := fmt.Sprintf(`INSERT INTO %s (id, username, description) values (uuid_generate_v4(), $1, $1) RETURNING *;`, streamTable)

	err := s.db.Get(&stream, qSelect, *user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			row := s.db.QueryRowx(qInsert, *user.Username)
			if err = row.Scan(&stream.Id, &stream.Username, &stream.Description); err != nil {
				return stream, err
			}
			return stream, nil
		}
		return stream, err
	}
	return stream, nil
}

func (s *StreamPostgres) GetStream(username string) (models.Stream, error) {
	var stream models.Stream

	query := fmt.Sprintf(`SELECT * FROM %s WHERE username=$1`, streamTable)

	err := s.db.Get(&stream, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return stream, nil
		}
		return stream, err
	}
	return stream, err
}

func (s *StreamPostgres) ChangeDescByUsername(stream models.PutStream) (models.Stream, error) {
	var item models.Stream

	query := fmt.Sprintf(`UPDATE %s SET description = $1 WHERE username = $2 RETURNING *;`, streamTable)

	row := s.db.QueryRowx(query, *stream.Description, *stream.Username)
	if err := row.Scan(&item.Id, &item.Username, &item.Description); err != nil {
		return item, err
	}
	return item, nil
}
