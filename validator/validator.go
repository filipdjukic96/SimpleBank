package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)

	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}

func ValidateUsername(username string) error {
	err := ValidateString(username, 3, 100)
	if err != nil {
		return err
	}

	if !isValidUsername(username) {
		return fmt.Errorf("must contain only lowercase letters, digits or underscore")
	}

	return nil
}

func ValidateFullName(fullname string) error {
	err := ValidateString(fullname, 3, 100)
	if err != nil {
		return err
	}

	if !isValidFullName(fullname) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	err := ValidateString(email, 3, 200)
	if err != nil {
		return err
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("is not a valid email address")
	}

	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("value must be a positive integer")
	}

	return nil
}

func ValidateSecretCode(code string) error {
	return ValidateString(code, 32, 128)
}
