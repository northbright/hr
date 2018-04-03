package leave

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/northbright/uuid"
)

var (
	// Types represents the leave types.
	Types = []string{
		"annual leave",
		"sick leave",
		"maternity leave",
		"paternity leave",
		"personal leave",
		"others",
	}

	Statuses = []string{
		"processing",
		"rejected",
		"approved",
	}

	// ErrInvalidType is the error of invalid leave type.
	ErrInvalidType = fmt.Errorf("invalid leave type")
	// ErrInvalidStatus is the error of invalid leave application status.
	ErrInvalidStatus = fmt.Errorf("invalid leave status")
)

type LeaveTime struct {
	Begin time.Time
	End   time.Time
}

// Application represents the leave application.
type Application struct {
	// EmployeeID is the employee ID of applicant.
	EmployeeID string `redis:"employee_id" json:"employee_id"`
	// TimeStr contains leave begin / end time for one or more days.
	// Format: YYYY-MM-DD/HH:MM-HH:MM
	// Use "," for multiple days.
	// e.g. "2018-04-02/08:00-17:00,2018-04-03/08:00-12:00"
	TimeStr string `redis:"leave_time" json:"leave_time"`
	// Type is the leave type. See Types for available types.
	Type string `redis:"type" json:"type"`
	// Status is the status of leave application.
	Status string `redis:"status" json:"status"`
}

func ValidType(leaveType string) bool {
	for _, t := range Types {
		if leaveType == t {
			return true
		}
	}
	return false
}

func ValidStatus(status string) bool {
	for _, s := range Statuses {
		if status == s {
			return true
		}
	}
	return false
}

func ParseLeaveTimes(timeStr string) (map[string]LeaveTime, error) {
	str := strings.Trim(timeStr, "")
	arr := strings.Split(str, ",")
}

func (a *Application) Valid() (bool, error) {
	if !ValidLeaveType(a.Type) {
		return false, ErrInvalidType
	}

	if !ValidStatus(a.Status) {
		return false, ErrInvalidStatus
	}

}
