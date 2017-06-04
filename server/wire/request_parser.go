package wire

import (
	"net/url"
	"sort"
	"strings"
	"fmt"
	"net/http"
)

type Mapping struct {
	field    *string
	required bool
}

func Required(field *string) *Mapping {
	return &Mapping{
		field:    field,
		required: true,
	}
}

func Optional(field *string) *Mapping {
	return &Mapping{
		field:    field,
		required: false,
	}
}

func ParseRequest(httpRequest http.Request, fields map[string]*Mapping) error {
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

func setValueOnRequest(fieldName string, values url.Values, missingFields *[]string, field *Mapping) {
	value := values.Get(fieldName)
	if strings.TrimSpace(value) == "" && field.required {
		*missingFields = append(*missingFields, fieldName)
	}
	*field.field = value
}