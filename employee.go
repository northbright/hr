package hr

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/northbright/validate"
)

type Employee struct {
	Name           string `json:"name"`
	Sex            string `json:"sex"`
	IDCardNo       string `json:"id_card_no"`
	MobilePhoneNum string `json:"mobile_phone_num"`
}

var (
	ErrInvalidName           = fmt.Errorf("invalid name")
	ErrInvalidSex            = fmt.Errorf("invalid sex")
	ErrInvalidIDCardNo       = fmt.Errorf("invalid ID card number")
	ErrInvalidMobilePhoneNum = fmt.Errorf("invalid mobile phone number")
	ErrNotUnique             = fmt.Errorf("at least one unique item is not unique")
)

func (e *Employee) Valid() error {
	if e.Name == "" {
		return ErrInvalidName
	}

	if e.Sex != "m" && e.Sex != "f" && e.Sex != "n" {
		return ErrInvalidSex
	}

	if !validate.ValidIDCardNo(e.IDCardNo) {
		return ErrInvalidIDCardNo
	}

	if !validate.ValidMobilePhoneNum(e.MobilePhoneNum) {
		return ErrInvalidMobilePhoneNum
	}
	return nil
}

func (e *Employee) CheckUniqueItems(db *sqlx.DB) error {
	var n int64
	stat := `SELECT COUNT(*) FROM employee
WHERE data @> jsonb_build_object('id_card_no',$1::text)
OR data @> jsonb_build_object('mobile_phone_num',$2::text)

`
	if err := db.Get(&n, stat, e.IDCardNo, e.MobilePhoneNum); err != nil {
		return err
	}

	if n > 0 {
		return ErrNotUnique
	}
	return nil
}

func CreateEmployee(db *sqlx.DB, e *Employee) (int64, error) {
	stat := `
	INSERT INTO employee (data) 
	VALUES ($1)
	RETURNING id`
	var ID int64

	e.Sex = UpdateSex(e.Sex)
	if err := e.Valid(); err != nil {
		return 0, err
	}

	if err := e.CheckUniqueItems(db); err != nil {
		return 0, err
	}

	jsonData, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}

	if err = db.QueryRow(stat, string(jsonData)).Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}

func GetEmployee(db *sqlx.DB, ID int64) (string, error) {
	var jsonStr string
	stat := `SELECT data FROM employee WHERE id = $1`

	err := db.Get(&jsonStr, stat, ID)
	switch err {
	case sql.ErrNoRows:
		return `{}`, nil
	case nil:
		return jsonStr, nil
	default:
		return `{}`, err
	}
}

func RemoveAllEmployees(db *sqlx.DB) error {
	stat := `DELETE FROM employee`

	_, err := db.Exec(stat)
	return err
}

func QueryEmployeeByIDCardNo(db *sqlx.DB, IDCardNo string) (string, error) {
	var jsonStr string
	stat := `SELECT data FROM employee WHERE data @> jsonb_build_object('id_card_no',$1::text)`

	err := db.Get(&jsonStr, stat, IDCardNo)
	switch err {
	case sql.ErrNoRows:
		return `{}`, nil
	case nil:
		return jsonStr, nil
	default:
		return `{}`, err
	}
}

func QueryEmployeeByMobilePhoneNum(db *sqlx.DB, mobilePhoneNum string) (string, error) {
	var jsonStr string
	stat := `SELECT data FROM employee WHERE data @> jsonb_build_object('mobile_phone_num',$1::text)`

	err := db.Get(&jsonStr, stat, mobilePhoneNum)
	switch err {
	case sql.ErrNoRows:
		return `{}`, nil
	case nil:
		return jsonStr, nil
	default:
		return `{}`, err
	}
}
