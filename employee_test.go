package hr_test

import (
	"log"

	"github.com/northbright/hr"
)

func ExampleCreateEmployee() {
	employeeData := [][]string{
		{"Frank", "m", "310104198101010000", "13100000000"},
		{"Bob", "m", "310104198201010000", "13300000000"},
		{"Alice", "f", "310104198302020000", "13500000000"},
	}

	// Remove all employees in the table.
	if err := hr.RemoveAllEmployees(db); err != nil {
		log.Printf("RemoveAllEmployees() error: %v", err)
		return
	}

	var IDs []int64

	// Create employees.
	for _, data := range employeeData {
		if ID, err := hr.CreateEmployee(db, data[0], data[1], data[2], data[3]); err != nil {
			log.Printf("CreateEmployee() error: %v", err)
			return
		} else {
			IDs = append(IDs, ID)
			log.Printf("CreateEmployee() OK. ID = %v, name: %v, sex: %v, ID card num: %v, mobile phone num: %v", ID, data[0], data[1], data[2], data[3])
		}
	}

	// Get employee data by ID.
	for _, ID := range IDs {
		if e, err := hr.GetEmployee(db, ID); err != nil {
			log.Printf("GetEmployee(%v) error: %v", ID, err)
			return
		} else {
			log.Printf("GetEmployee(%v) OK. employee: %v", ID, e)
		}
	}

	// Get employee by ID card number.
	for _, data := range employeeData {
		if e, err := hr.QueryEmployeeByIDCardNo(db, data[2]); err != nil {
			log.Printf("QueryEmployeeByIDCardNo(%v) error: %v", data[2], err)
			return
		} else {
			log.Printf("QueryEmployeeByIDCardNo(%v) OK. employee: %v", data[2], e)
		}
	}

	// Output:
}
