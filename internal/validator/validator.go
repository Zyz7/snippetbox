package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Returns true if the fieldErrors map does not contain any entries
func (v *Validator) Valid() bool {

	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// Adding error messages to the new NonFieldErrors slice
func (v *Validator) AddNonFieldError(message string) {

	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {

	// We need to initialize the map first if it is not already initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exits := v.FieldErrors[key]; !exits {
		v.FieldErrors[key] = message
	}
}

// Adds an error  message to the fieldErrors map only if a validation check is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Returns true if a value is a not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Returns true if a values is in a list of specif permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
