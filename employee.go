package hr

import (
	//"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID             int64  `db:id`
	Name           string `db:"name"`
	Sex            string `db:"sex"`
	IDCardNo       string `db:"id_card_no"`
	MobilePhoneNum string `db:"mobile_phone_num"`
}

var ()

func updateSex(sex string) string {
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

func CreateEmployee(db *sqlx.DB, name, sex, IDCardNo, mobilePhoneNum string) (int64, error) {
	stat := `
	INSERT INTO employee (name, sex, id_card_no, mobile_phone_num) 
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var ID int64

	sex = updateSex(sex)

	err := db.QueryRow(stat, name, sex, IDCardNo, mobilePhoneNum).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}
