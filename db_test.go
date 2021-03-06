package hr_test

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/northbright/hr"
)

func ExampleInitDB() {
	host := "localhost"
	port := "5432"
	user := "postgres"
	dbname := "test"

	log.Printf(`Make sure to run the command before test: psql -U postgres -c 'CREATE DATABASE test;'`)
	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v dbname=%v sslmode=disable",
		host, port, user, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("sqlx.Open() error: %v", err)
		return
	}
	defer db.Close()

	if err = hr.InitDB(db); err != nil {
		log.Printf("hr.InitDB() error: %v", err)
		return
	}
	log.Printf("hr.InitDB() OK")
	// Output:
}
