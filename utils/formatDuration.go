package utils

import "fmt"

func FormatDuration(seconds int64) string {
	mins := seconds / 60
	minutes := addZero(mins)
	hours := addZero(mins / 60)
  remainingSeconds := addZero(seconds % 60)

	if hours != "00" {
		return hours + ":" + minutes + ":" + remainingSeconds
	} else {
		return minutes + ":" + remainingSeconds
	}
}

func addZero(num int64) string {
	formatted := ""
	if num <= 9 {
		formatted = "0" + fmt.Sprint(num)
	} else {
		formatted = fmt.Sprint(num)
	}

	return formatted
}