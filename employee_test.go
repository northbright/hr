package hr_test

import (
	"log"

	"github.com/northbright/hr"
)

func ExampleCreateEmployee() {
	employees := []hr.EmployeeData{
		hr.EmployeeData{"Frank", "m", "310104198101010000", "13100000000"},
		hr.EmployeeData{"Bob", "m", "310104198201010000", "13300000000"},
		hr.EmployeeData{"Alice", "f", "310104198302020000", "13500000000"},
	}

	// Remove all employees in the table.
	if err := hr.RemoveAllEmployees(db); err != nil {
		log.Printf("RemoveAllEmployees() error: %v", err)
		return
	}

	var IDs []string

	// Create employees.
	for _, e := range employees {
		if ID, err := hr.CreateEmployee(db, &e); err != nil {
			log.Printf("CreateEmployee() error: %v", err)
			return
		} else {
			IDs = append(IDs, ID)
			log.Printf("CreateEmployee() OK. ID = %v, employee: %v", ID, e)
		}
	}

	// Get employee data by ID.
	for _, ID := range IDs {
		if jsonStr, err := hr.GetEmployee(db, ID); err != nil {
			log.Printf("GetEmployee(%v) error: %v", ID, err)
			return
		} else {
			log.Printf("GetEmployee(%v) OK. JSON: %v", ID, jsonStr)
		}
	}

	// Get employee by ID card number.
	for _, e := range employees {
		if jsonStr, err := hr.QueryEmployeeByIDCardNo(db, e.IDCardNo); err != nil {
			log.Printf("QueryEmployeeByIDCardNo(%v) error: %v", e.IDCardNo, err)
			return
		} else {
			log.Printf("QueryEmployeeByIDCardNo(%v) OK. JSON: %v", e.IDCardNo, jsonStr)
		}
	}

	// Test of invalid ID card number.
	jsonStr, err := hr.QueryEmployeeByIDCardNo(db, "0000")
	log.Printf("\nTest of invalid ID card number: JSON: %v, err: %v", jsonStr, err)

	// Get employee by mobile phone number.
	for _, e := range employees {
		if jsonStr, err := hr.QueryEmployeeByMobilePhoneNum(db, e.MobilePhoneNum); err != nil {
			log.Printf("QueryEmployeeByMobilePhoneNum(%v) error: %v", e.MobilePhoneNum, err)
			return
		} else {
			log.Printf("QueryEmployeeByMobilePhoneNum(%v) OK. JSON: %v", e.MobilePhoneNum, jsonStr)
		}
	}

	// Test of invalid mobile phone number.
	jsonStr, err = hr.QueryEmployeeByMobilePhoneNum(db, "0000")
	log.Printf("\nTest of invalid mobile phone number: JSON: %v, err: %v", jsonStr, err)

	// Output:
}
