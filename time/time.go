package time

import (
	"time"
)

// DateTime format
func DateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}

	r := t.Format("2006-01-02 15:04:05")
	return &r
}

// Date format
func Date(t *time.Time) *string {
	if t == nil {
		return nil
	}

	r := t.Format("2006-01-02")
	return &r
}
