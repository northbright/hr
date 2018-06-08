package hr

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS employee (
    id TEXT PRIMARY KEY,
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

func GetJSONData(db *sqlx.DB, sqlStat string, args ...interface{}) ([]byte, error) {
	var jsonData []byte

	err := db.Get(&jsonData, sqlStat, args...)
	switch err {
	case sql.ErrNoRows:
		return []byte(`{}`), nil
	case nil:
		return jsonData, nil
	default:
		return []byte(`{}`), err
	}
}
