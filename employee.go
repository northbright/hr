package hr

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/northbright/uuid"
	"github.com/northbright/validate"
)

type EmployeeData struct {
	Name           string `json:"name"`
	Sex            string `json:"sex"`
	IDCardNo       string `json:"id_card_no"`
	MobilePhoneNum string `json:"mobile_phone_num"`
}

// Employee contains ID and full employee data.
type Employee struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	*EmployeeData
}

var (
	ErrInvalidName           = fmt.Errorf("invalid name")
	ErrInvalidSex            = fmt.Errorf("invalid sex")
	ErrInvalidIDCardNo       = fmt.Errorf("invalid ID card number")
	ErrInvalidMobilePhoneNum = fmt.Errorf("invalid mobile phone number")
	ErrNotUnique             = fmt.Errorf("at least one unique item is not unique")
)

func UpdateSex(sex string) string {
	sex = strings.ToLower(sex)

	switch sex {
	case "m", "f", "n":
		return sex
	case "男":
		return "m"
	case "女":
		return "f"
	default:
		return "n"
	}
}

func (e *EmployeeData) Valid() error {
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

func (e *EmployeeData) CheckUniqueItems(db *sqlx.DB) error {
	var n int64
	stat := `SELECT COUNT(*) FROM employee
WHERE data @> jsonb_build_object('id_card_no',$1::text)
OR data @> jsonb_build_object('mobile_phone_num',$2::text)`

	if err := db.Get(&n, stat, e.IDCardNo, e.MobilePhoneNum); err != nil {
		return err
	}

	if n > 0 {
		return ErrNotUnique
	}
	return nil
}

func CreateEmployee(db *sqlx.DB, e *EmployeeData) (string, error) {
	var (
		err error
		ID  string
	)

	stat := `INSERT INTO employee (id, created, data) VALUES ($1, $2, $3)`

	e.Sex = UpdateSex(e.Sex)
	if err = e.Valid(); err != nil {
		return "", err
	}

	if err = e.CheckUniqueItems(db); err != nil {
		return "", err
	}

	if ID, err = uuid.New(); err != nil {
		return "", err
	}

	nanoSeconds := time.Now().UnixNano()

	newEmployee := Employee{ID, nanoSeconds, e}
	jsonData, err := json.Marshal(newEmployee)
	if err != nil {
		return "", err
	}

	if _, err = db.Exec(stat, ID, nanoSeconds, string(jsonData)); err != nil {
		return "", err
	}
	return ID, nil
}

func GetEmployee(db *sqlx.DB, ID string) ([]byte, error) {
	stat := `SELECT data FROM employee WHERE id = $1`
	return GetJSONData(db, stat, ID)
}

func RemoveAllEmployees(db *sqlx.DB) error {
	stat := `DELETE FROM employee`

	_, err := db.Exec(stat)
	return err
}

func GetEmployeeByIDCardNo(db *sqlx.DB, IDCardNo string) ([]byte, error) {
	stat := `SELECT data FROM employee
WHERE data @> jsonb_build_object('id_card_no',$1::text)`
	return GetJSONData(db, stat, IDCardNo)
}

func GetEmployeeByMobilePhoneNum(db *sqlx.DB, mobilePhoneNum string) ([]byte, error) {
	stat := `SELECT data FROM employee
WHERE data @> jsonb_build_object('mobile_phone_num',$1::text)`
	return GetJSONData(db, stat, mobilePhoneNum)
}
