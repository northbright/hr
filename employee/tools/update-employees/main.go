package main

import (
	"flag"
	"fmt"

	"github.com/northbright/csvhelper"
	"github.com/northbright/hr/employee"
	"github.com/northbright/redishelper"
)

const (
	// incorrectArgsMsg is the message while arguments error occurs.
	incorrectArgsMsg string = "Incorrect arguments, please see usage:\n"
	// usage is the message of usage.
	usage       string = "usage:\nupdate-employees -a=<Redis server address> -p=<Redis password> -f=<csv file>\nEx: update-employees -a='127.0.0.1:6379' -p='' -f='employees.csv'"
	maxActive   int    = 10
	maxIdle     int    = 10
	idleTimeout        = 60
	wait        bool   = true
)

func main() {
	var (
		redisServerAddr string
		redisPassword   string
		inputFile       string
	)

	flag.StringVar(&redisServerAddr, "a", "", "Redis server address. Ex: -a='127.0.0.1:6379'")
	flag.StringVar(&redisPassword, "p", "", "Redis password. Ex: -p='my_password'")
	flag.StringVar(&inputFile, "f", "", "path of CSV file which contains employee data. Ex: -f='employees.csv'")
	flag.Parse()

	if inputFile == "" || redisServerAddr == "" {
		fmt.Printf("%s\n", incorrectArgsMsg)
		flag.PrintDefaults()
		fmt.Printf("%s\n", usage)
		return
	}

	records, err := csvhelper.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("ReadFile() error: %v", err)
		return
	}

	// Create a Redis pool.
	pool := redishelper.NewRedisPool(
		redisServerAddr,
		redisPassword,
		maxActive,
		maxIdle,
		idleTimeout,
		wait,
	)
	defer pool.Close()

	for i, record := range records {
		// Skip header.
		if i == 0 {
			continue
		}

		e := employee.Employee{
			Name:           record[0],
			Sex:            record[1],
			IDCardNo:       record[2],
			MobilePhoneNum: record[3],
		}

		exists, ID, err := employee.Exists(pool, &e)
		if err != nil {
			fmt.Printf("Exists() error: %v", err)
			return
		}

		if exists {
			fmt.Printf("employee exists, updating...\n")
			err = employee.Set(pool, ID, &e)
			if err != nil {
				fmt.Printf("Set() error: %v", err)
				return
			}
		} else {
			fmt.Printf("new employee, adding...\n")
			ID, err = employee.Add(pool, &e)
			if err != nil {
				fmt.Printf("Add() error: %v", err)
				return
			}
		}
	}
	fmt.Printf("Update done")
}
