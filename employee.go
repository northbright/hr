package hr

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/northbright/validate"
)

type Employee struct {
	ID             int64  `db:"id"`
	Name           string `db:"name"`
	Sex            string `db:"sex"`
	IDCardNo       string `db:"id_card_no"`
	MobilePhoneNum string `db:"mobile_phone_num"`
}

var (
	ErrInvalidIDCardNo       = fmt.Errorf("invalid ID card number")
	ErrInvalidMobilePhoneNum = fmt.Errorf("invalid mobile phone number")
)

func CreateEmployee(db *sqlx.DB, name, sex, IDCardNo, mobilePhoneNum string) (int64, error) {
	stat := `
	INSERT INTO employee (name, sex, id_card_no, mobile_phone_num) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var ID int64

	sex = UpdateSex(sex)
	if !validate.ValidIDCardNo(IDCardNo) {
		return 0, ErrInvalidIDCardNo
	}

	if !validate.ValidMobilePhoneNum(mobilePhoneNum) {
		return 0, ErrInvalidMobilePhoneNum
	}

	err := db.QueryRow(stat, name, sex, IDCardNo, mobilePhoneNum).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func GetEmployee(db *sqlx.DB, ID int64) (*Employee, error) {
	stat := `SELECT * FROM employee WHERE id = $1`
	e := &Employee{}

	err := db.Get(e, stat, ID)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func RemoveAllEmployees(db *sqlx.DB) error {
	stat := `DELETE FROM employee`

	_, err := db.Exec(stat)
	return err
}

func QueryEmployeeByIDCardNo(db *sqlx.DB, IDCardNo string) (*Employee, error) {
	stat := `SELECT * FROM employee WHERE id_card_no = $1`
	e := &Employee{}

	err := db.Get(e, stat, IDCardNo)
	if err != nil {
		return nil, err
	}
	return e, nil
}
