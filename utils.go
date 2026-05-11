package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
)

func jsonRoundTrip(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func URLJoin(endpoint string, paths ...string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(fmt.Errorf("unable to parse endpoint: %s", endpoint))
	}
	paths = append([]string{u.Path}, paths...)
	u.Path = path.Join(paths...)
	return u.String()
}
