package jira

import (
	"bytes"
	"encoding/json"
)

type AuthParams struct {
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
}

type SessionInfo struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

type LoginInfo struct {
	FailedLoginCount    int    `json:"failedLoginCount,omitempty" yaml:"failedLoginCount,omitempty"`
	LastFailedLoginTime string `json:"lastFailedLoginTime,omitempty" yaml:"lastFailedLoginTime,omitempty"`
	LoginCount          int    `json:"loginCount,omitempty" yaml:"loginCount,omitempty"`
	PreviousLoginTime   string `json:"previousLoginTime,omitempty" yaml:"previousLoginTime,omitempty"`
}

type AuthSuccess struct {
	LoginInfo *LoginInfo   `json:"loginInfo,omitempty" yaml:"loginInfo,omitempty"`
	Session   *SessionInfo `json:"session,omitempty" yaml:"session,omitempty"`
}

type CurrentUser struct {
	LoginInfo *LoginInfo `json:"loginInfo,omitempty" yaml:"loginInfo,omitempty"`
	Name      string     `json:"name,omitempty" yaml:"name,omitempty"`
	Self      string     `json:"self,omitempty" yaml:"self,omitempty"`
}

type AuthOptions struct {
	Username string
	Password string
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login
func (j *Jira) NewSession(ap *AuthOptions) (*AuthSuccess, error) {
	return NewSession(j.UA, j.Endpoint, ap)
}

func NewSession(ua HttpClient, endpoint string, ap *AuthOptions) (*AuthSuccess, error) {
	req := &AuthParams{Username: ap.Username, Password: ap.Password}
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := URLJoin(endpoint, "rest/auth/1/session")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &AuthSuccess{}
		return results, json.NewDecoder(resp.Body).Decode(results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-currentUser
func (j *Jira) GetSession() (*CurrentUser, error) {
	return GetSession(j.UA, j.Endpoint)
}

func GetSession(ua HttpClient, endpoint string) (*CurrentUser, error) {
	uri := URLJoin(endpoint, "rest/auth/1/session")
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &CurrentUser{}
		return results, json.NewDecoder(resp.Body).Decode(results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-logout
func (j *Jira) DeleteSession() error {
	return DeleteSession(j.UA, j.Endpoint)
}

func DeleteSession(ua HttpClient, endpoint string) error {
	uri := URLJoin(endpoint, "rest/auth/1/session")
	resp, err := ua.Delete(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}
