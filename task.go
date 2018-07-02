package hr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/northbright/uuid"
)

type TaskCommentData struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

type TaskComment struct {
	Created int64 `json:"created"`
	*TaskCommentData
}

type TaskData struct {
	Assigner  string   `json:"assigner"`
	Assignees []string `json:"assignees,omitempty"`
	Priority  int      `json:"priority"`
	Closed    bool     `json:"closed"`
	Tags      []string `json:"tags,omitempty"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
}

type Task struct {
	ID       string        `json:"id"`
	Created  int64         `json:"created"`
	Comments []TaskComment `json:"comments,omitempty"`
	*TaskData
}

var (
	ErrTaskAssignerNotFound      = fmt.Errorf("task assigner not found")
	ErrTaskAssigneeNotFound      = fmt.Errorf("at least one of the task assignees not found")
	ErrEmptyTaskComment          = fmt.Errorf("empty task comment")
	ErrTaskCommentAuthorNotFound = fmt.Errorf("task comment author not found")
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

	stat := `INSERT INTO task (id, created, data) VALUES ($1, $2, $3)`

	if err = t.Valid(db); err != nil {
		return "", err
	}

	if ID, err = uuid.New(); err != nil {
		return "", err
	}

	nanoSeconds := time.Now().UnixNano()

	newTask := Task{ID: ID, Created: nanoSeconds, TaskData: t}
	jsonData, err := json.Marshal(newTask)
	if err != nil {
		return "", err
	}

	if _, err = db.Exec(stat, ID, nanoSeconds, string(jsonData)); err != nil {
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

func GetTaskCountByAssigner(db *sqlx.DB, assigner string) (int64, error) {
	var n int64
	stat := `SELECT COUNT(*) FROM task
WHERE data @> jsonb_build_object('assigner',$1::text)`
	if err := db.Get(&n, stat, assigner); err != nil {
		return 0, err
	}
	return n, nil
}

func GetTasksByAssigner(db *sqlx.DB, assigner string, limit, offset int64) ([][]byte, error) {
	stat := `SELECT data FROM task
WHERE data @> jsonb_build_object('assigner',$1::text)
ORDER BY created DESC LIMIT $2 OFFSET $3`
	return SelectJSONData(db, stat, assigner, limit, offset)
}

func GetTaskCountByAssignee(db *sqlx.DB, assignee string) (int64, error) {
	var n int64
	stat := `SELECT COUNT(*) FROM task
WHERE (data->'assignees')::jsonb ? $1::text`
	if err := db.Get(&n, stat, assignee); err != nil {
		return 0, err
	}
	return n, nil
}

func GetTasksByAssignee(db *sqlx.DB, assignee string, limit, offset int64) ([][]byte, error) {
	stat := `SELECT data FROM task
WHERE (data->'assignees')::jsonb ? $1::text
ORDER BY created DESC LIMIT $2 OFFSET $3`
	return SelectJSONData(db, stat, assignee, limit, offset)
}

func (c *TaskCommentData) Valid(db *sqlx.DB) error {
	var n int64

	if c.Content == "" {
		return ErrEmptyTaskComment
	}

	// Check assigner employee ID.
	stat := `SELECT COUNT(*) FROM employee
WHERE ID = $1`

	if err := db.Get(&n, stat, c.Author); err != nil {
		return err
	}

	if n <= 0 {
		return ErrTaskCommentAuthorNotFound
	}
	return nil
}

func CreateTaskComment(db *sqlx.DB, taskID string, c *TaskCommentData) error {
	stat := `UPDATE task SET data = jsonb_set(
data, '{comments}', COALESCE(data->'comments', '[]'::jsonb) || jsonb_build_array($1::jsonb), true)
WHERE id = $2`

	if err := c.Valid(db); err != nil {
		return err
	}

	nanoSeconds := time.Now().UnixNano()

	comment := TaskComment{nanoSeconds, c}
	jsonData, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	if _, err = db.Exec(stat, string(jsonData), taskID); err != nil {
		return err
	}
	return nil
}
