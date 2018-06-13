package hr_test

import (
	"log"

	"github.com/northbright/hr"
)

func ExampleCreateTask() {
	var (
		IDs     []string
		taskIDs []string
	)

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

	// Create employees.
	for _, e := range employees {
		ID, err := hr.CreateEmployee(db, &e)
		if err != nil {
			log.Printf("CreateEmployee() error: %v", err)
			return
		}
		IDs = append(IDs, ID)
		log.Printf("CreateEmployee() OK. ID = %v, employee: %v", ID, e)

	}

	// Create task.
	tasks := []hr.TaskData{
		hr.TaskData{
			Assigner:  IDs[0],
			Assignees: []string{IDs[1], IDs[2]},
			Priority:  1,
			Closed:    false,
			Tags:      []string{"IT", "PC"},
			Title:     "Purchase PC",
			Content:   "Purchase 10 PCs",
		},
		hr.TaskData{
			Assigner:  IDs[0],
			Assignees: []string{IDs[1]},
			Priority:  2,
			Closed:    false,
			Tags:      []string{"IT", "Printer"},
			Title:     "Fix the Printer",
			Content:   "Laser Printer does not work",
		},
	}

	// Remove all tasks in the table.
	if err := hr.RemoveAllTasks(db); err != nil {
		log.Printf("RemoveAllTasks() error: %v", err)
		return
	}

	// Create tasks.
	for _, t := range tasks {
		ID, err := hr.CreateTask(db, &t)
		if err != nil {
			log.Printf("CreateTask() error: %v", err)
			return
		}
		taskIDs = append(taskIDs, ID)
		log.Printf("CreateTask() OK. ID = %v, task: %v", ID, t)
	}

	// Get task data by ID.
	for _, ID := range taskIDs {
		jsonData, err := hr.GetTask(db, ID)
		if err != nil {
			log.Printf("GetTask(%v) error: %v", ID, err)
			return
		}
		log.Printf("GetTask(%v) OK. JSON: %s", ID, jsonData)
	}

	// Get tasks by assigner.
	n, err := hr.GetTaskCountByAssigner(db, IDs[0])
	if err != nil {
		log.Printf("GetTaskCountByAssigner() error: %v", err)
		return
	}

	limit := int64(2)
	offset := int64(0)
	for offset = 0; offset < n; offset += limit {
		jsonDataArr, err := hr.GetTasksByAssigner(db, IDs[0], limit, offset)
		if err != nil {
			log.Printf("GetTasksByAssigner() error: %v", err)
			return
		}
		log.Printf("GetTasksByAssigner() OK. LIMIT: %v, OFFSET: %v", limit, offset)
		for _, data := range jsonDataArr {
			log.Printf("JSON: %s\n", data)
		}
	}

	// Output:
}
