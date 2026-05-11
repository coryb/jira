package jira

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	jv2 "github.com/ctreminiom/go-atlassian/v2/jira/v2"
	jv3 "github.com/ctreminiom/go-atlassian/v2/jira/v3"
	common "github.com/ctreminiom/go-atlassian/v2/service/common"
)

type HttpClient interface {
	Delete(url string) (*http.Response, error)
	Do(*http.Request) (*http.Response, error)
	GetJSON(url string) (*http.Response, error)
	Post(url, bodyType string, body io.Reader) (*http.Response, error)
	Put(url, bodyType string, body io.Reader) (*http.Response, error)
}

// pathPrefixDoer prepends a URL path prefix to all requests. go-atlassian
// builds absolute paths (e.g. /rest/api/3/...) via url.ResolveReference,
// which discards any base path in the endpoint URL.
type pathPrefixDoer struct {
	inner  HttpClient
	prefix string
}

func (w *pathPrefixDoer) Do(req *http.Request) (*http.Response, error) {
	req.URL.Path = w.prefix + req.URL.Path
	return w.inner.Do(req)
}

func atlassianBaseURL(endpoint string) (baseURL, prefix string, err error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}
	return u.Scheme + "://" + u.Host, strings.TrimRight(u.Path, "/"), nil
}

func atlassianDoer(ua HttpClient, prefix string) common.HTTPClient {
	if prefix != "" {
		return &pathPrefixDoer{inner: ua, prefix: prefix}
	}
	return ua
}

// newAtlassianClient creates a go-atlassian v3 client (ADF description format).
func newAtlassianClient(ua HttpClient, endpoint string) (*jv3.Client, error) {
	baseURL, prefix, err := atlassianBaseURL(endpoint)
	if err != nil {
		return nil, err
	}
	return jv3.New(atlassianDoer(ua, prefix), baseURL)
}

// newAtlassianClientV2 creates a go-atlassian v2 client (plain-text description format).
func newAtlassianClientV2(ua HttpClient, endpoint string) (*jv2.Client, error) {
	baseURL, prefix, err := atlassianBaseURL(endpoint)
	if err != nil {
		return nil, err
	}
	return jv2.New(atlassianDoer(ua, prefix), baseURL)
}
