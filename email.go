package std

import (
	"regexp"
	"strings"
)

// ---------------------------------------------------------------------------------------------------------------------
// Const
// ---------------------------------------------------------------------------------------------------------------------

const EmailLengthMin = 8 // Min 2 chars for each part + "@" and "."
const EmailLengthMax = 255

var EmailNil = Email{}

// ---------------------------------------------------------------------------------------------------------------------
// Struct
// ---------------------------------------------------------------------------------------------------------------------

type Email struct {
	value string
}

// ---------------------------------------------------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------------------------------------------------

var emailRegexpCompiled = regexp.MustCompile(func() string {
	pattern := ""
	pattern += "[a-zA-Z0-9]{1}"
	pattern += "((\\.|-|_)?[a-zA-Z0-9]+(-|_)?)*"
	pattern += "("
	pattern += "\\+"
	pattern += "([a-zA-Z0-9]+(\\.|-|_)?)*"
	pattern += "[a-zA-Z0-9\\-_]{1}"
	pattern += ")?"
	pattern += "@"
	pattern += "([a-zA-Z0-9]+((-|_)?[a-zA-Z0-9]+)*)+"
	pattern += "(\\.[a-zA-Z0-9]+((-|_)?[a-zA-Z0-9]+)*)+"
	return "(^" + pattern + "$)"
}())

// EmailFromString
//
// Errors:
//   - ErrorValidation
func EmailFromString(value string) (Email, error) {
	if value == "" {
		return EmailNil, NewErrorValidationFf("Email value cannot be empty")
	}

	if !(EmailLengthMin <= len(value) && len(value) <= EmailLengthMax) {
		return EmailNil, NewErrorValidationFf("Email length must be between %d and %d", EmailLengthMin, EmailLengthMax)
	}

	if !emailRegexpCompiled.MatchString(value) {
		return EmailNil, NewErrorValidationFf("Email value should match regexp %q", emailRegexpCompiled.String())
	}

	return Email{value: value}, nil
}

func EmailFromStringMust(value string) Email {
	email, err := EmailFromString(value)
	if err != nil {
		panic(err)
	}
	return email
}

// ---------------------------------------------------------------------------------------------------------------------
// State
// ---------------------------------------------------------------------------------------------------------------------

func (e Email) IsNil() bool {
	return e == EmailNil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Name() string {
	name, _, _ := strings.Cut(e.value, "@")
	return name
}

func (e Email) Domain() string {
	_, domain, _ := strings.Cut(e.value, "@")
	return domain
}

// ---------------------------------------------------------------------------------------------------------------------
