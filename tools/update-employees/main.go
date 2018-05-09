package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/northbright/csvhelper"
	"github.com/northbright/pathhelper"
	//"github.com/northbright/hr"
)

type PostgreSQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSL      bool   `json:"ssl"`
}

type Config struct {
	PostgreSQL PostgreSQLConfig
	CSVFile    string `json:"csv_file"`
}

func main() {
	var (
		err    error
		config Config
	)

	// Load config from './config.json'.
	// You may rename "config.example.json" to "config.json" and modify it.
	// It looks like:
	//{
	//      "postgresql": {
	//          "host":"localhost",
	//          "port":"5432",
	//          "user":"postgres",
	//          "password":"",
	//          "db_name":"hr",
	//          "ssl":false
	//      },
	//      "csv_file":"/home/xx/hr-employees.csv"
	//}

	// Get Absolute path of "./config.json"
	configFile, _ := pathhelper.GetAbsPath("config.json")
	if err = loadConfig(configFile, &config); err != nil {
		log.Printf("loadConfig() error: %v", err)
		return
	}

	records, err := csvhelper.ReadFile(config.CSVFile)
	if err != nil {
		fmt.Printf("ReadFile() error: %v", err)
		return
	}

	fmt.Printf("records: %v\n", records)
}

// loadConfig loads app config.
func loadConfig(configFile string, config *Config) error {
	// Load Conifg
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("load config file error: %v", err)

	}

	if err = json.Unmarshal(buf, config); err != nil {
		return fmt.Errorf("parse config err: %v", err)
	}

	return nil
}