package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// ValidateRequired checks a map of field names to values and returns a slice of missing field names
func ValidateRequired(fields map[string]string) []string {
	var missing []string
	for k, v := range fields {
		if strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}
	return missing
}

// ValidateEmail performs a basic email validation
func ValidateEmail(email string) bool {
	if strings.TrimSpace(email) == "" {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ValidatePhone performs a basic phone number validation (digits, spaces, +, -, parentheses)
func ValidatePhone(phone string) bool {
	if strings.TrimSpace(phone) == "" {
		return false
	}
	re := regexp.MustCompile(`^[0-9()+\-\s]{6,20}$`)
	return re.MatchString(phone)
}

// ValidateDate checks whether a date string matches YYYY-MM-DD
func ValidateDate(dateStr string) bool {
	if strings.TrimSpace(dateStr) == "" {
		return false
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// ValidateDateRange verifies that the given date is not before today and within 'days' days ahead
func ValidateDateRange(dateStr string, days int) (bool, error) {
	if !ValidateDate(dateStr) {
		return false, errors.New("invalid date format")
	}
	d, _ := time.Parse("2006-01-02", dateStr)
	today := time.Now().Truncate(24 * time.Hour)
	end := today.AddDate(0, 0, days)

	if d.Before(today) {
		return false, nil
	}
	if d.After(end) {
		return false, nil
	}
	return true, nil
}
