package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/northbright/csvhelper"
	"github.com/northbright/hr"
	"github.com/northbright/pathhelper"
)

type PostgreSQLConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type CSVConfig struct {
	File      string `json:"file"`
	HasHeader bool   `json:"has_header"`
}

type Config struct {
	PostgreSQL PostgreSQLConfig `json:"postgresql"`
	CSV        CSVConfig        `json:"csv"`
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
	//          "dbname":"hr",
	//      },
	//      "csv": {
	//          "file":"/home/xx/hr-employees.csv",
	//          "has_header":true
	//      }
	//}

	// Get Absolute path of "./config.json"
	configFile, _ := pathhelper.GetAbsPath("config.json")
	if err = loadConfig(configFile, &config); err != nil {
		log.Printf("loadConfig() error: %v", err)
		return
	}

	records, err := csvhelper.ReadFile(config.CSV.File)
	if err != nil {
		log.Printf("ReadFile() error: %v", err)
		return
	}

	//fmt.Printf("records: %v\n", records)

	info := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		config.PostgreSQL.Host,
		config.PostgreSQL.Port,
		config.PostgreSQL.User,
		config.PostgreSQL.Password,
		config.PostgreSQL.DBName,
	)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		log.Printf("sqlx.Connect() error: %v", err)
		return
	}

	defer db.Close()

	if err = hr.InitDB(db); err != nil {
		log.Printf("hr.InitDB() error: %v", err)
		return
	}

	if config.CSV.HasHeader {
		records = records[1:]
	}

	for _, record := range records {
		name := record[0]
		sex := record[1]
		IDCardNo := record[2]
		mobilePhoneNum := ""

		phoneNums := strings.Split(record[3], "/")
		if len(phoneNums) == 2 {
			mobilePhoneNum = phoneNums[0]
		} else {
			mobilePhoneNum = record[3]
		}

		_, err := hr.CreateEmployee(db, name, sex, IDCardNo, mobilePhoneNum)
		if err != nil {
			log.Printf("hr.CreateEmployee() error: %v, name: %v, sex: %v, ID Card No: %v, mobile phone num: %v\n",
				err, name, sex, IDCardNo, mobilePhoneNum)
			continue
		}
		//log.Printf("CreateEmployee() OK. ID: %v, name: %v, sex: %v, id_card_no: %v, mobile_phone_num: %v",
		//	ID,
		//	name,
		//	sex,
		//	IDCardNo,
		//	mobilePhoneNum,
		//)
	}

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