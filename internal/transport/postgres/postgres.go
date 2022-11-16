package postgres

import "github.com/jmoiron/sqlx"

func NewPostgresDB(strConn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", strConn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
