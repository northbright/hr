package hr

type Employee struct {
	ID             int    `db:id`
	Name           string `db:"name"`
	Sex            string `db:"sex"`
	IDCardNo       string `db:"id_card_no"`
	MobilePhoneNum string `db:"mobile_phone_num"`
}
