package validation

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	reUpper = regexp.MustCompile(`[A-Z]`)
	reLower = regexp.MustCompile(`[a-z]`)
	reDigit = regexp.MustCompile(`[0-9]`)
	reSym   = regexp.MustCompile(`[^A-Za-z0-9]`)
)

func Password(fl validator.FieldLevel) bool {
	pw := fl.Field().String()
	pw = strings.TrimSpace(pw)

	if len(pw) < 8 {
		return false
	}
	if !reUpper.MatchString(pw) {
		return false
	}
	if !reLower.MatchString(pw) {
		return false
	}
	if !reDigit.MatchString(pw) {
		return false
	}
	if !reSym.MatchString(pw) {
		return false
	}
	if strings.Contains(pw, " ") {
		return false
	}

	parent := fl.Parent()
	if parent.IsValid() && parent.Kind() == reflect.Struct {
		emailField := parent.FieldByName("Email")
		if emailField.IsValid() && emailField.Kind() == reflect.String {
			if email := emailField.String(); email != "" {
				if i := strings.IndexByte(email, '@'); i > 0 {
					local := strings.ToLower(email[:i])
					if local != "" && strings.Contains(strings.ToLower(pw), local) {
						return false
					}
				}
			}
		}
	}
	return true
}

func RequiredStrict(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return field.String() != ""
	case reflect.Ptr:
		return !field.IsNil()
	}

	return field.IsValid() && !field.IsZero()
}

func OmitemptyStrict(fl validator.FieldLevel) bool {
	field := fl.Field()

	if !field.IsValid() || field.IsZero() {
		return true
	}

	if field.Kind() == reflect.String {
		return field.String() != ""
	}

	return true
}
