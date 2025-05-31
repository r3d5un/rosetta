package validator

import (
	"slices"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func PermittedValuesCaseInsensitive(values []string, permitted []string) (string, bool) {
	for i, str := range values {
		values[i] = strings.ToLower(str)
	}
	for i, str := range permitted {
		permitted[i] = strings.ToLower(str)
	}

	return PermittedValues(values, permitted)
}

func PermittedValues(values []string, permitted []string) (string, bool) {
	for _, value := range values {
		if !slices.Contains(permitted, value) {
			return value, false
		}
	}

	return "", true
}
