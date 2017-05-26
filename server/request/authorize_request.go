package request

import (
	"net/url"
	"strings"
	"sort"
)

type AuthorizeRequest struct {
	RequestType string
	ClientId    string
}

type ParseError struct {
	MissingFields    []string
	UnexpectedFields [] string
}

func ParseAuthorizeRequest(values url.Values) (*AuthorizeRequest, *ParseError) {
	authorizeRequest := AuthorizeRequest{}

	fields := map[string]*string{
		"request_type": &authorizeRequest.RequestType,
		"client_id":    &authorizeRequest.ClientId,
	}

	parseError := ParseRequest(values, authorizeRequest, fields)

	if parseError != nil {
		return nil, parseError
	} else {
		return &authorizeRequest, nil
	}
}

func ParseRequest(values url.Values, request interface{}, fields map[string]*string) *ParseError {
	parseError := ParseError{
		MissingFields:    []string{},
		UnexpectedFields: []string{},
	}

	for key, value := range fields {
		setValueOnRequest(key, values, &parseError, value)
	}

	if len(parseError.MissingFields) > 0 {
		sort.Strings(parseError.MissingFields)
		return &parseError
	} else {
		return nil;
	}
}

func setValueOnRequest(fieldName string, values url.Values, parseError *ParseError, field *string) {
	value := values.Get(fieldName)
	if strings.TrimSpace(value) == "" {
		parseError.MissingFields = append(parseError.MissingFields, fieldName)
	}
	*field = value
}
