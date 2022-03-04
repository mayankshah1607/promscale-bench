package prom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PromQLResponse struct {
	Response http.Response
	ExecTime time.Duration
	Err      error
}

type PromQLRequest struct {
	Query string
	Start string
	End   string
	Step  string
}

// CreateHTTPRequst builds a PromQL HTTP request for range queries.
// This request can be initiated by calling ExecRequestWithClient.
func CreateHTTPRequest(promURL string, q PromQLRequest) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/query_range", promURL))
	if err != nil {
		return &http.Request{},
			fmt.Errorf("error parsing URL: %v", err)
	}

	v := url.Values{}
	v.Add("query", q.Query)
	v.Add("start", q.Start)
	v.Add("end", q.End)
	v.Add("step", q.Step)

	u.RawQuery = v.Encode()

	return http.NewRequest(
		"GET", u.String(), nil,
	)
}

// ParseHTTPResponse is a helper function for parsing HTTP response.
func ParseHTTPResponse(body io.ReadCloser) (map[string]interface{}, error) {
	buffer, err := ioutil.ReadAll(body)
	if err != nil {
		return make(map[string]interface{}), err
	}
	defer body.Close()

	var payload interface{}
	json.Unmarshal(buffer, &payload)
	return payload.(map[string]interface{}), nil
}

// measureRequestDuration is a middleware for measuring the execution time
// of a given function `f`.
func measureRequestDuration(f func() (*http.Response, error)) (*http.Response,
	time.Duration, error) {
	start := time.Now()
	res, err := f()
	elapsed := time.Since(start)

	if err != nil {
		return &http.Response{}, elapsed, err
	}
	return res, elapsed, err
}

// ExecRequestWithClient initiates a PromQL HTTP request
func ExecRequestWithClient(req *http.Request,
	c *http.Client) PromQLResponse {

	res, execTime, err := measureRequestDuration(func() (*http.Response, error) {
		return c.Do(req)
	})
	if err != nil {
		return PromQLResponse{
			Err: fmt.Errorf("error making HTTP request: %v", err),
		}
	}

	return PromQLResponse{
		Response: *res,
		ExecTime: execTime,
	}
}
