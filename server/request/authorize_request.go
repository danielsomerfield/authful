package request

import (
	"net/url"
	"strings"
)

type AuthorizeRequest struct {
	requestType string
	clientId    string
}

type ParseError struct {
	missingFields []string
}

func ParseAuthorizeRequest(values url.Values) (*AuthorizeRequest, *ParseError) {
	missingFields := []string{}
	authorizeRequest := AuthorizeRequest{}

	authorizeRequest.requestType, missingFields = getValue("request_type", values, missingFields)
	authorizeRequest.clientId, missingFields = getValue("client_id", values, missingFields)

	if len(missingFields) > 0 {
		return nil, &ParseError{
			missingFields,
		}
	} else {
		return &authorizeRequest, nil
	}

}
func getValue(fieldName string, values url.Values, missingFields []string) (string, []string) {
	value := values.Get(fieldName)
	if strings.TrimSpace(value) == "" {
		missingFields = append(missingFields, fieldName)
	}
	return value, missingFields
}
