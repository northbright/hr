package employee_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/northbright/hr/employee"
	"github.com/northbright/redishelper"
)

// RedisConfig represents the configuration of Redis.
type RedisConfig struct {
	// ServerAddr is the address of Redis server. e.g. ":6379".
	ServerAddr string `json:"server_addr"`
	// Password is the password of Redis server.
	Password string `json:"password"`
	// MaxActive is the maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int `json:"max_active"`
	// MaxIdle is maximum number of idle connections in the pool.
	MaxIdle int `json:"max_idle"`
	// IdleTimeout is the timeout: Close connections after remaining idle for this duration.
	// If the value is zero, then idle connections are not closed.
	// Applications should set the timeout to a value less than the server's timeout.
	IdleTimeout int `json:"idle_time_out"`
	// Wait is the flag:
	// If Wait is true and the pool is at the MaxActive limit,
	// then Get() waits for a connection to be returned to the pool before returning.
	Wait bool `json:"wait"`
}

type Config struct {
	Redis RedisConfig `json:"redis"`
}

func Example() {
	var (
		err    error
		config Config
	)

	// Load config from file("config.json").
	// It looks like this:
	//{
	//    "redis":{
	//        "server_addr":":6379",
	//        "password":"",
	//        "max_active":10,
	//        "max_idle":10,
	//        "idle_time_out":60,
	//        "wait":true
	//    }
	//}
	if err = loadConfig("config.json", &config); err != nil {
		log.Printf("loadConfig() error: %v", err)
		return
	}

	// New a redis pool.
	pool := redishelper.NewRedisPool(
		config.Redis.ServerAddr,
		config.Redis.Password,
		config.Redis.MaxActive,
		config.Redis.MaxIdle,
		config.Redis.IdleTimeout,
		config.Redis.Wait,
	)
	defer pool.Close()

	// Add new employees.
	employees := []employee.Employee{
		employee.Employee{
			Name:           "Bob",
			IDCardNo:       "310104222222222222",
			Sex:            "female",
			MobilePhoneNum: "13777777777",
		},
		employee.Employee{
			Name:           "John",
			IDCardNo:       "310104111111111111",
			Sex:            "male",
			MobilePhoneNum: "13333333333",
		},
	}

	var IDs []string
	for _, e := range employees {
		ID, err := employee.Add(pool, &e)
		if err != nil {
			log.Printf("Add() error: %v", err)
			return
		}
		IDs = append(IDs, ID)
	}
	log.Printf("Add() OK")

	// Get employee data.
	ID := IDs[0]
	e, err := employee.Get(pool, ID)
	if err != nil {
		log.Printf("Get() error: %v", err)
		return
	}
	log.Printf("Get() OK.\nID: %v\nEmployee: %v", ID, e)

	// Set employee data.
	e.Name = "Bob.Smith"
	e.Sex = "male"
	e.IDCardNo = "310104000000000000"
	e.MobilePhoneNum = "18000000000"
	err = employee.Set(pool, ID, e)
	if err != nil {
		log.Printf("Set() error: %v", err)
		return
	}

	// Get employee data again.
	log.Printf("Set() OK.\nID: %v\nEmployee: %v", ID, e)

	// Get employee by ID card number.
	IDCardNo := "310104000000000000"
	ID, err = employee.GetIDByIDCardNo(pool, IDCardNo)
	if err != nil {
		log.Printf("GetIDByIDCardNo() error: %v", err)
		return
	}
	log.Printf("GetIDByIDCardNo() OK.\nID Card No: %v\nID: %v", IDCardNo, ID)

	// Get employee by mobile phone number.
	mobilePhoneNum := "18000000000"
	ID, err = employee.GetIDByMobilePhoneNum(pool, mobilePhoneNum)
	if err != nil {
		log.Printf("GetIDBMobilePhoneNum() error: %v", err)
		return
	}
	log.Printf("GetIDBMobilePhoneNum() OK.\nmobile phone num: %v\nID: %v", mobilePhoneNum, ID)

	// Delete employee.
	for _, ID := range IDs {
		if err = employee.Del(pool, ID); err != nil {
			log.Printf("Del() error: %v, ID: %v", err, ID)
			return
		}
	}
	log.Printf("Del() OK")

	// Output:
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
