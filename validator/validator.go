package validator

import (
	"fmt"
	"time"
)

func ValidatePost(title, content, date string) error {
	if title == "" || len(title) > 200 {
		return fmt.Errorf(`{"error":"Title is required and must be less than 200 characters"}`)
	}
	if content == "" {
		return fmt.Errorf(`{"error":"Content is required"}`)
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return fmt.Errorf(`{"error":"Date must be in YYYY-MM-DD format"}`)
	}
	return nil
}
