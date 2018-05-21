package hr

import (
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    sex TEXT NOT NULL,
    id_card_no TEXT NOT NULL UNIQUE,
    mobile_phone_num TEXT NOT NULL UNIQUE
);
`

func InitDB(db *sqlx.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}
