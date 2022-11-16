package postgres

import (
	"database/sql"
	"fmt"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/models"
)

const (
	broadcastTable = "broadcasts"
	usersTable     = "users"
	imagesTable    = "images"
)

type BroadcastsPostgres struct {
	db *sqlx.DB
}

func NewBroadcastsPostgres(db *sqlx.DB) *BroadcastsPostgres {
	return &BroadcastsPostgres{db}
}

func (b *BroadcastsPostgres) GetBroadcasts() ([]models.Broadcasts, error) {
	var items []models.Broadcasts
	query := fmt.Sprintf(
		`SELECT id, name, owner, description, previewurl, streamkey, start_time, life
				FROM %s WHERE life='%s' ORDER by start_time ASC;`,
		broadcastTable, models.Created.String())

	if err := b.db.Select(&items, query); err != nil {
		return items, err
	}

	return items, nil
}

func (b *BroadcastsPostgres) GetBroadcastById(id types.UUID) (models.Broadcasts, error) {
	var item models.Broadcasts
	q := fmt.Sprintf(
		"SELECT id, name, owner, description, previewurl, streamkey, start_time, life FROM %s WHERE id = $1",
		broadcastTable)
	if err := b.db.Get(&item, q, id); err != nil {
		if err == sql.ErrNoRows {
			return item, nil
		}
		return item, err
	}
	return item, nil
}

func (b *BroadcastsPostgres) DeleteBroadcast(id types.UUID) (api.SIdentifier, error) {
	var item api.SIdentifier

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id;", broadcastTable)
	if err := b.db.QueryRowx(query, id).StructScan(&item); err != nil {
		if err == sql.ErrNoRows {
			return item, nil
		}
		return item, err
	}
	return item, nil
}

func (b *BroadcastsPostgres) CheckAdminUser(user api.SUsername) (bool, error) {
	var item api.SUsername
	query := fmt.Sprintf(`SELECT username FROM %s WHERE username = $1`, usersTable)
	if err := b.db.Get(&item, query, *user.Username); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (b *BroadcastsPostgres) GetPastBroadcastAll() ([]models.Broadcasts, error) {
	var items = make([]models.Broadcasts, 0)
	query := fmt.Sprintf(
		`SELECT id, name, owner, description, previewurl, streamkey, start_time, life
				FROM %s WHERE life='%s' ORDER by start_time DESC;`,
		broadcastTable, models.Past.String())
	if err := b.db.Select(&items, query); err != nil {
		return nil, err
	}
	return items, nil
}

func (b *BroadcastsPostgres) GetPastBroadcastOwner(user api.SUsername) ([]models.Broadcasts, error) {
	var items = make([]models.Broadcasts, 0)
	query := fmt.Sprintf(
		`SELECT id, name, owner, description, previewurl, streamkey, start_time, life
				FROM %s WHERE life='%s' AND owner = '%s' ORDER by start_time DESC;`,
		broadcastTable, models.Past.String(), *user.Username)
	if err := b.db.Select(&items, query); err != nil {
		return nil, err
	}
	return items, nil
}

func (b *BroadcastsPostgres) CreateBroadcast(item models.PostBroadcast) (models.Broadcasts, error) {
	var broadcast models.Broadcasts

	tx := b.db.MustBegin()

	query := fmt.Sprintf(`INSERT INTO %s
	(id, life, name, owner, description, streamkey , start_time) 
	values (uuid_generate_v4(), '%s', $1, $2, $3, $4, $5) RETURNING id;`, broadcastTable, models.Created)

	row := tx.QueryRowx(query, *item.Name, *item.Owner, *item.Description, *item.StreamKey, *item.StartTime)
	if err := row.Scan(&broadcast.Id); err != nil {
		err = tx.Rollback()
		if err != nil {
			return broadcast, err
		}
		return broadcast, err
	}

	queryEmpty := fmt.Sprintf("INSERT INTO %s (id) VALUES ('%s');", imagesTable, *broadcast.Id)

	if _, err := tx.Exec(queryEmpty); err != nil {
		err = tx.Rollback()
		if err != nil {
			return broadcast, err
		}
		return broadcast, err
	}

	err := tx.Commit()
	if err != nil {
		return broadcast, err
	}

	return broadcast, nil
}

func (b *BroadcastsPostgres) ChangeBroadcast(i models.PutBroadcast) (models.Broadcasts, error) {
	var item models.Broadcasts

	q := fmt.Sprintf(`UPDATE %s
		SET name = $1, description = $2, streamkey = $3 , start_time = $4 WHERE id = $5 RETURNING *;`,
		broadcastTable)

	if err := b.db.QueryRowx(q, *i.Name, *i.Description, *i.StreamKey, *i.StartTime, *i.Id).StructScan(&item); err != nil {
		if err == sql.ErrNoRows {
			return item, nil
		}
		return item, err
	}

	return item, nil
}
