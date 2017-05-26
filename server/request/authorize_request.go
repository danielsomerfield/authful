package request

import (
	"net/url"
	"strings"
	"sort"
)

type AuthorizeRequest struct {
	ResponseType string
	ClientId     string
	RedirectURI  string
	Scope        string
	State        string
}

type ParseError struct {
	MissingFields []string
}

func ParseAuthorizeRequest(values url.Values) (*AuthorizeRequest, *ParseError) {
	authorizeRequest := AuthorizeRequest{}

	fields := map[string]*mapping{
		"response_type": required(&authorizeRequest.ResponseType),
		"client_id":    required(&authorizeRequest.ClientId),
		"redirect_uri": optional(&authorizeRequest.RedirectURI),
		"scope":        optional(&authorizeRequest.Scope),
		"state":        optional(&authorizeRequest.State),
	}

	parseError := ParseRequest(values, fields)

	if parseError != nil {
		return nil, parseError
	} else {
		return &authorizeRequest, nil
	}
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

type mapping struct {
	field    *string
	required bool
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
