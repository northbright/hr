package hr_test

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/northbright/hr"
)

var (
	err error
	db  *sqlx.DB
)

func init() {
	host := "localhost"
	port := "5432"
	user := "postgres"
	dbname := "test"

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=disable",
		host, port, user, dbname)

	if db, err = sqlx.Open("postgres", psqlInfo); err != nil {
		log.Printf("sqlx.Open() error: %v", err)
		return
	}

	if err = hr.InitDB(db); err != nil {
		log.Printf("hr.InitDB() error: %v", err)
		return
	}
}
