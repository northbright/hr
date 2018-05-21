package hr

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID             int64  `db:id`
	Name           string `db:"name"`
	Sex            string `db:"sex"`
	IDCardNo       string `db:"id_card_no"`
	MobilePhoneNum string `db:"mobile_phone_num"`
}

var (
	ErrInvalidSex = fmt.Errorf("invalid sex")
)

func CreateEmployee(db *sqlx.DB, name, sex, IDCardNo, mobilePhoneNum string) (int64, error) {
	stat := `
	INSERT INTO employee (name, sex, id_card_no, mobile_phone_num) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var ID int64

	if sex == "男" {
		sex = "m"
	} else if sex == "女" {
		sex = "f"
	} else if sex == "" {
		sex = "n"
	}

	if sex != "m" && sex != "f" && sex != "n" {
		fmt.Printf("invalid sex: %v", sex)
		return 0, ErrInvalidSex

	}

	err := db.QueryRow(stat, name, sex, IDCardNo, mobilePhoneNum).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}
