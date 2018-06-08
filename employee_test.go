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
		if jsonData, err := hr.GetEmployee(db, ID); err != nil {
			log.Printf("GetEmployee(%v) error: %v", ID, err)
			return
		} else {
			log.Printf("GetEmployee(%v) OK. JSON: %s", ID, jsonData)
		}
	}

	// Get employee by ID card number.
	for _, e := range employees {
		if jsonData, err := hr.GetEmployeeByIDCardNo(db, e.IDCardNo); err != nil {
			log.Printf("GetEmployeeByIDCardNo(%v) error: %v", e.IDCardNo, err)
			return
		} else {
			log.Printf("GetEmployeeByIDCardNo(%v) OK. JSON: %s", e.IDCardNo, jsonData)
		}
	}

	// Test of invalid ID card number.
	jsonData, err := hr.GetEmployeeByIDCardNo(db, "0000")
	log.Printf("\nTest of invalid ID card number: JSON: %s, err: %v", jsonData, err)

	// Get employee by mobile phone number.
	for _, e := range employees {
		if jsonData, err := hr.GetEmployeeByMobilePhoneNum(db, e.MobilePhoneNum); err != nil {
			log.Printf("GetEmployeeByMobilePhoneNum(%v) error: %v", e.MobilePhoneNum, err)
			return
		} else {
			log.Printf("GetEmployeeByMobilePhoneNum(%v) OK. JSON: %s", e.MobilePhoneNum, jsonData)
		}
	}

	// Test of invalid mobile phone number.
	jsonData, err = hr.GetEmployeeByMobilePhoneNum(db, "0000")
	log.Printf("\nTest of invalid mobile phone number: JSON: %s, err: %v", jsonData, err)

	// Output:
}
