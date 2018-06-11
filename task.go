package hr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/northbright/uuid"
)

type TaskComment struct {
	author  string `json:"author"`
	content string `json:"content"`
}

type TaskData struct {
	Assigner  string        `json:"assigner"`
	Assignees []string      `json:"assignees"`
	Priority  int           `json:"priority"`
	Closed    bool          `json:"closed"`
	Tags      []string      `json:"tags"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Comments  []TaskComment `json:"comments"`
}

type Task struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	*TaskData
}

var (
	ErrTaskAssignerNotFound = fmt.Errorf("task assigner not found")
	ErrTaskAssigneeNotFound = fmt.Errorf("at least one of the task assignees not found")
)

func (t *TaskData) Valid(db *sqlx.DB) error {
	var n int64

	// Check assigner employee ID.
	stat := `SELECT COUNT(*) FROM employee
WHERE ID = $1`

	if err := db.Get(&n, stat, t.Assigner); err != nil {
		return err
	}

	if n <= 0 {
		return ErrTaskAssignerNotFound
	}

	// Check assignees employee IDs.
	for _, assignee := range t.Assignees {
		if err := db.Get(&n, stat, assignee); err != nil {
			return err
		}

		if n <= 0 {
			return ErrTaskAssigneeNotFound
		}
	}
	return nil
}

func CreateTask(db *sqlx.DB, t *TaskData) (string, error) {
	var (
		err error
		ID  string
	)

	stat := `INSERT INTO task (id, data) VALUES ($1, $2)`

	if err = t.Valid(db); err != nil {
		return "", err
	}

	if ID, err = uuid.New(); err != nil {
		return "", err
	}

	newTask := Task{ID, time.Now().Unix(), t}
	jsonData, err := json.Marshal(newTask)
	if err != nil {
		return "", err
	}

	if _, err = db.Exec(stat, ID, string(jsonData)); err != nil {
		return "", err
	}
	return ID, nil
}

func GetTask(db *sqlx.DB, ID string) ([]byte, error) {
	stat := `SELECT data FROM task WHERE id = $1`
	return GetJSONData(db, stat, ID)
}

func RemoveAllTasks(db *sqlx.DB) error {
	stat := `DELETE FROM task`

	_, err := db.Exec(stat)
	return err
}

func GetTasksByAssigner(db *sqlx.DB, assigner string) ([][]byte, error) {
	stat := `SELECT data FROM task
WHERE data @> jsonb_build_object('assigner',$1::text)`
	return SelectJSONData(db, stat, assigner)
}
