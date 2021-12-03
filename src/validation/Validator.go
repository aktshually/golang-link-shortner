package validation

import (
	"fmt"
	"reflect"
	"strings"
)

func Validate(s interface{}) (errors []string) {
	errors = []string{}
	fields := reflect.VisibleFields(reflect.TypeOf(s))

	for _, field := range fields {
		tags := field.Tag
		value := reflect.ValueOf(s).FieldByName(field.Name)

		if tags.Get("required") == "true" && strings.Trim(value.String(), " ") == "" {
			errors = append(errors, fmt.Sprintf("'%s' can not be empty", tags.Get("json")))
		}

		// TODO: other validations
	}
	return
}
