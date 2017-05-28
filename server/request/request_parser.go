package request

import (
	"net/url"
	"sort"
	"strings"
)

type mapping struct {
	field    *string
	required bool
}

type ParseError struct {
	MissingFields []string
}

func required(field *string) *mapping {
	return &mapping{
		field:    field,
		required: true,
	}
}

func optional(field *string) *mapping {
	return &mapping{
		field:    field,
		required: false,
	}
}

func ParseRequest(values url.Values, fields map[string]*mapping) *ParseError {
	parseError := ParseError{
		MissingFields: []string{},
	}

	for key, mapping := range fields {
		setValueOnRequest(key, values, &parseError, mapping)
	}

	if len(parseError.MissingFields) > 0 {
		sort.Strings(parseError.MissingFields)
		return &parseError
	} else {
		return nil;
	}
}

func setValueOnRequest(fieldName string, values url.Values, parseError *ParseError, field *mapping) {
	value := values.Get(fieldName)
	if strings.TrimSpace(value) == "" && field.required {
		parseError.MissingFields = append(parseError.MissingFields, fieldName)
	}
	*field.field = value
}