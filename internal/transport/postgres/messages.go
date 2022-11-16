package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/alexm24/golang/internal/models"
)

const (
	messagesTable = "messages"
)

type MessagesPostgres struct {
	db *sqlx.DB
}

func NewMessagesPostgres(db *sqlx.DB) *MessagesPostgres {
	return &MessagesPostgres{db}
}

func (m *MessagesPostgres) GetMessageByChannel(channel string) ([]models.Messages, error) {
	var msg = make([]models.Messages, 0)

	query := fmt.Sprintf(
		`SELECT id, fullname, text, time, username, avatar, is_question, is_anon, reactions
		FROM %s WHERE channel = $1 ORDER by time ASC;`,
		messagesTable)

	if err := m.db.Select(&msg, query, channel); err != nil {
		return msg, err
	}

	return msg, nil
}

func (m *MessagesPostgres) CreateMsg(channel string, msg models.PostMessage) (models.Messages, error) {
	var resMsg models.Messages

	q := fmt.Sprintf(
		`INSERT INTO %s 
		(id, channel, username, fullname, text, avatar, time, is_question, is_anon) 
		VALUES 
		(uuid_generate_v4(), $1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;`,
		messagesTable)
	row := m.db.QueryRowx(q, channel, *msg.Username, *msg.Fullname, *msg.Text, *msg.Avatar, *msg.Time, *msg.IsQuestion, *msg.IsAnon)

	err := row.StructScan(&resMsg)
	if err != nil {
		return resMsg, err
	}

	return resMsg, nil
}

func (m *MessagesPostgres) DeleteMessages(channel string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE channel=$1", messagesTable)
	_, err := m.db.Exec(query, channel)
	if err != nil {
		return err
	}
	return nil
}

func (m *MessagesPostgres) AddReaction(item models.PostReactionMsg) (models.Messages, error) {
	var msg models.Messages
	q := fmt.Sprintf(
		`UPDATE %s SET reactions = reactions || '{ "%s": "%s" }'::jsonb WHERE id=$1 RETURNING *;`,
		messagesTable, *item.Username, *item.Type)
	row := m.db.QueryRowx(q, *item.Id)
	err := row.StructScan(&msg)
	return msg, err
}

func (m *MessagesPostgres) DeleteReaction(item models.PatchReactionMsg) (models.Messages, error) {
	var msg models.Messages
	q := fmt.Sprintf(`UPDATE %s SET reactions = reactions - $1 WHERE id=$2 RETURNING *;`, messagesTable)
	row := m.db.QueryRowx(q, *item.Username, *item.Id)
	err := row.StructScan(&msg)
	return msg, err
}
