package employee

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/gomodule/redigo/redis"
	"github.com/northbright/uuid"
)

type Employee struct {
	Name           string `redis:"name" json:"name"`
	Sex            string `redis:"sex" json:"sex"`
	IDCardNo       string `redis:"id_card_no" json:"id_card_no"`
	MobilePhoneNum string `redis:"mobile_phone_num" json:"mobile_phone_num"`
}

var (
	ErrInvalidName           = fmt.Errorf("invalid name")
	ErrInvalidSex            = fmt.Errorf("invalid sex")
	ErrInvalidIDCardNo       = fmt.Errorf("invalid ID card number")
	ErrInvalidMobilePhoneNum = fmt.Errorf("invalid mobile phone number")
	ErrNotExist              = fmt.Errorf("employee does not exist")
	ErrAlreadyExists         = fmt.Errorf("employee already exists")
)

func ValidName(name string) bool {
	l := utf8.RuneCount([]byte(name))
	if l <= 0 || l > 60 {
		return false
	}
	return true
}

func ValidSex(sex string) bool {
	if sex != "male" && sex != "female" {
		return false
	}
	return true
}

func (e *Employee) Valid() error {
	if !ValidName(e.Name) {
		return ErrInvalidName
	}

	if !ValidSex(e.Sex) {
		return ErrInvalidSex
	}

	if e.IDCardNo == "" {
		return ErrInvalidIDCardNo
	}

	if e.MobilePhoneNum == "" {
		return ErrInvalidMobilePhoneNum
	}

	return nil
}

func GetIDByIDCardNo(pool *redis.Pool, IDCardNo string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	k := "hr:employees:index:id_card_no_to_id"
	ID, err := redis.String(conn.Do("HGET", k, IDCardNo))
	switch err {
	case nil:
		return ID, nil
	case redis.ErrNil:
		return "", nil
	default:
		return "", err
	}
}

func GetIDByMobilePhoneNum(pool *redis.Pool, mobilePhoneNum string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	k := "hr:employees:index:mobile_phone_num_to_id"
	ID, err := redis.String(conn.Do("HGET", k, mobilePhoneNum))
	switch err {
	case nil:
		return ID, nil
	case redis.ErrNil:
		return "", nil
	default:
		return "", err
	}
}

func GetIDsByName(pool *redis.Pool, name string) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	k := fmt.Sprintf("hr:employees:index:name_to_ids:%v", name)
	IDs, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	switch err {
	case nil:
		return IDs, nil
	case redis.ErrNil:
		return []string{}, nil
	default:
		return []string{}, err
	}
}

func Exists(pool *redis.Pool, e *Employee) (bool, string, error) {
	ID, err := GetIDByIDCardNo(pool, e.IDCardNo)
	if err != nil {
		return false, "", err
	}
	if ID != "" {
		return true, ID, nil
	}

	ID, err = GetIDByMobilePhoneNum(pool, e.MobilePhoneNum)
	if err != nil {
		return false, "", err
	}
	if ID != "" {
		return true, ID, nil
	}

	return false, "", nil
}

func Add(pool *redis.Pool, e *Employee) (string, error) {
	err := e.Valid()
	if err != nil {
		return "", err
	}

	exists, _, err := Exists(pool, e)
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrAlreadyExists
	}

	ID, err := uuid.New()
	if err != nil {
		return "", err
	}

	t := time.Now().Unix()

	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")

	k := fmt.Sprintf("hr:employee:%v", ID)
	conn.Send("HMSET", redis.Args{}.Add(k).AddFlat(e)...)

	k = "hr:employees:index:id_card_no_to_id"
	conn.Send("HSET", k, e.IDCardNo, ID)

	k = "hr:employees:index:mobile_phone_num_to_id"
	conn.Send("HSET", k, e.MobilePhoneNum, ID)

	k = fmt.Sprintf("hr:employees:index:name_to_ids:%v", e.Name)
	conn.Send("ZADD", k, t, ID)

	k = "hr:employees"
	conn.Send("ZADD", k, t, ID)

	if _, err = conn.Do("EXEC"); err != nil {
		return "", err
	}

	return ID, nil
}

func Del(pool *redis.Pool, ID string) error {
	e, err := Get(pool, ID)
	if err != nil {
		return err
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")

	k := fmt.Sprintf("hr:employee:%v", ID)
	conn.Send("DEL", k)

	k = "hr:employees:index:id_card_no_to_id"
	conn.Send("HDEL", k, e.IDCardNo)

	k = "hr:employees:index:mobile_phone_num_to_id"
	conn.Send("HDEL", k, e.MobilePhoneNum)

	k = fmt.Sprintf("hr:employees:index:name_to_ids:%v", e.Name)
	conn.Send("ZREM", k, ID)

	k = "hr:employees"
	conn.Send("ZREM", k, ID)

	if _, err = conn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

func GetAllIDs(pool *redis.Pool) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	k := "hr:employees"
	IDs, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []string{}, err
	}

	return IDs, nil
}

func IDExists(pool *redis.Pool, ID string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	k := fmt.Sprintf("hr:employee:%v", ID)
	exists, err := redis.Bool(conn.Do("EXISTS", k))
	if err != nil {
		return false, err
	}
	return exists, nil
}

func Get(pool *redis.Pool, ID string) (*Employee, error) {
	exists, err := IDExists(pool, ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotExist
	}

	conn := pool.Get()
	defer conn.Close()

	k := fmt.Sprintf("hr:employee:%v", ID)
	v, err := redis.Values(conn.Do("HGETALL", k))
	if err != nil {
		return nil, err
	}

	e := Employee{}
	if err = redis.ScanStruct(v, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func Set(pool *redis.Pool, ID string, e *Employee) error {
	var (
		nameChanged           bool
		IDCardNoChanged       bool
		mobilePhoneNumChanged bool
	)

	err := e.Valid()
	if err != nil {
		return err
	}

	// Get old employee data by ID.
	oldEmplyee, err := Get(pool, ID)
	if err != nil {
		return err
	}

	if oldEmplyee.Name != e.Name {
		nameChanged = true
	}

	if oldEmplyee.IDCardNo != e.IDCardNo {
		IDCardNoChanged = true
	}

	if oldEmplyee.MobilePhoneNum != e.MobilePhoneNum {
		mobilePhoneNumChanged = true
	}

	conn := pool.Get()
	pipedConn := pool.Get()
	defer conn.Close()
	defer pipedConn.Close()

	pipedConn.Send("MULTI")

	k := fmt.Sprintf("hr:employee:%v", ID)
	pipedConn.Send("HMSET", redis.Args{}.Add(k).AddFlat(e)...)

	if IDCardNoChanged {
		k = "hr:employees:index:id_card_no_to_id"
		pipedConn.Send("HDEL", k, oldEmplyee.IDCardNo)
		pipedConn.Send("HSET", k, e.IDCardNo, ID)
	}

	if mobilePhoneNumChanged {
		k = "hr:employees:index:mobile_phone_num_to_id"
		pipedConn.Send("HDEL", k, oldEmplyee.MobilePhoneNum)
		pipedConn.Send("HSET", k, e.MobilePhoneNum, ID)
	}

	if nameChanged {
		k = fmt.Sprintf("hr:employees:index:name_to_ids:%v", oldEmplyee.Name)
		t, err := redis.String(conn.Do("ZSCORE", k, ID))
		if err != nil {
			return err
		}

		pipedConn.Send("ZREM", k, ID)

		k = fmt.Sprintf("hr:employees:index:name_to_ids:%v", e.Name)
		pipedConn.Send("ZADD", k, t, ID)
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}
