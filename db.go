package hr

import (
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS employee (
    id SERIAL PRIMARY KEY,
    data JSONB
);

CREATE INDEX on employee USING GIN (data);
`

func InitDB(db *sqlx.DB) error {
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	return nil
}
