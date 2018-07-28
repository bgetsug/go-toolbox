package validation

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

var e164PhoneRegex = regexp.MustCompile(e164PhoneRegexString)

// Validates that the field is formatted like an E.164 phone number.
// This does NOT check that it is actually a valid phone number.
func IsE164Phone(fl validator.FieldLevel) bool {
	return e164PhoneRegex.MatchString(fl.Field().String())
}
