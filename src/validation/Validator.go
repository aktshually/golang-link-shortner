package validation

import (
	"fmt"
	"reflect"
	"strconv"
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

		if maxValue, _ := strconv.Atoi(tags.Get("max")); tags.Get("max") != "" && len(value.String()) > maxValue {
			errors = append(errors, fmt.Sprintf("'%s' must have a max length of %d characters", tags.Get("json"), maxValue))
		}

		if minValue, _ := strconv.Atoi(tags.Get("min")); tags.Get("min") != "" && len(value.String()) > minValue {
			errors = append(errors, fmt.Sprintf("'%s' must have a min length of %d characters", tags.Get("json"), minValue))
		}

		// TODO: other validations
	}
	return
}
