package hr

import (
	"strings"
)

func UpdateSex(sex string) string {
	sex = strings.ToLower(sex)

	switch sex {
	case "m", "f", "n":
		return sex
	case "男":
		return "m"
	case "女":
		return "f"
	default:
		return "n"
	}
}
