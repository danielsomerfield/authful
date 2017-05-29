package request

import (
	"net/url"
	"sort"
	"strings"
	"fmt"
	"net/http"
)

type mapping struct {
	field    *string
	required bool
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

func ParseRequest(httpRequest http.Request, fields map[string]*mapping) error {
	missingFields := []string{}

	httpRequest.ParseForm()

	for key, mapping := range fields {
		setValueOnRequest(key, httpRequest.Form, &missingFields, mapping)
	}

	if len(missingFields) > 0 {
		sort.Strings(missingFields)
		return fmt.Errorf("The following fields are required: %s", missingFields)
	} else {
		return nil
	}
}

func setValueOnRequest(fieldName string, values url.Values, missingFields *[]string, field *mapping) {
	value := values.Get(fieldName)
	if strings.TrimSpace(value) == "" && field.required {
		*missingFields = append(*missingFields, fieldName)
	}
	*field.field = value
}