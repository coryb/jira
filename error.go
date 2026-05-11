package jira

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ErrorCollection mirrors the Jira REST API error response body.
// It can be removed once all raw HTTP calls are migrated to go-atlassian.
type ErrorCollection struct {
	ErrorMessages []string          `json:"errorMessages,omitempty" yaml:"errorMessages,omitempty"`
	Errors        map[string]string `json:"errors,omitempty" yaml:"errors,omitempty"`
	Status        int               `json:"status,omitempty" yaml:"status,omitempty"`
}

func (e ErrorCollection) Error() string {
	if len(e.ErrorMessages) > 0 {
		return strings.Join(e.ErrorMessages, ". ")
	}
	parts := make([]string, 0, len(e.Errors))
	for k, v := range e.Errors {
		parts = append(parts, k+": "+v)
	}
	return strings.Join(parts, ". ")
}

func responseError(resp *http.Response) error {
	results := &ErrorCollection{}
	if err := json.NewDecoder(resp.Body).Decode(results); err != nil {
		results.Status = resp.StatusCode
		results.ErrorMessages = append(results.ErrorMessages, err.Error())
	}
	if len(results.ErrorMessages) == 0 && len(results.Errors) == 0 {
		results.Status = resp.StatusCode
		results.ErrorMessages = append(results.ErrorMessages, resp.Status)
	}
	return results
}
