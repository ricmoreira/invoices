package util

import "time"

// Parse_YYYYMMDD_Date returns a time.Time object for string dates formated as YYYY-MM-DD
func Parse_YYYYMMDD_Date(date string) (time.Time, error) {
	form := "2006-01-02"
	return time.Parse(form, date)
}
