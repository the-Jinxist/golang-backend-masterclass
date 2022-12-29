package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

//We're using regex to make sure our username only contains alphanum characters
var (
	isValidUserName = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)

	if n < minLength {
		return fmt.Errorf("%s is too short", value)
	}

	if n > maxLength {
		return fmt.Errorf("%s is too long", value)
	}

	return nil
}

func ValidateUserName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	//We're using regex to make sure our username only contains alphanum characters
	matchString := isValidUserName(value)
	if !matchString {
		return fmt.Errorf("must contain only lowercase alphabets or num characters")
	}

	return nil

}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 200)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 6, 100); err != nil {
		return err
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("input string %s is not a valid email address", value)
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	//We're using regex to make sure our username only contains alphanum characters
	matchString := isValidFullName(value)
	if !matchString {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil

}
