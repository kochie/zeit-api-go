package zeit

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type rateLimit struct {
	limit     int
	remaining int
	reset     time.Time
}

type Client struct {
	token      string
	rootUrl    string
	httpClient HttpClient
	rateLimit  *rateLimit
	team       string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient will create a new zeit client to apply api request to. Note that the team is defaulted to nothing, if you
// want to update the team then use Team.
func NewClient(token string) *Client {
	rl := rateLimit{1, 1, time.Now()}
	return &Client{
		token,
		"https://api.zeit.co",
		&http.Client{},
		&rl,
		"",
	}
}

// Team will set the team associated with the api client, to not use a team set with empty string.
func (c Client) Team(team string) {
	c.team = team
}

// closeResponseBody is a helper function to close the body of a http response and panic if there is an error closing
// the io writer.
func closeResponseBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		panic("http response couldn't be closed")
	}
}

// makeAndDoRequest will create the appropriate request and then send it to the endpoint specified. It will handle
// authentication, headers, and rate limiting.
func (c Client) makeAndDoRequest(httpMethod, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.rootUrl, endpoint)
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// if team is defined, add it to the url query
	if c.team != "" {
		q := req.URL.Query()
		q.Add("teamId", c.team)
		req.URL.RawQuery = q.Encode()
	}

	// do a check to see if the rate limit has been hit, if so wait until a request can be sent again
	if c.rateLimit.remaining == 0 && time.Now().Before(c.rateLimit.reset) {
		d := time.Now().Sub(c.rateLimit.reset)
		fmt.Println(fmt.Sprintf("Zeit rate limit hit, waiting for %f seconds", d.Seconds()))
		time.Sleep(d)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	c.rateLimit.remaining-- // in-case there are multiple threads

	if remaining, err := strconv.Atoi(resp.Header.Get("X-RateLimit-remaining")); err != nil {
		c.rateLimit.remaining = remaining
	}
	if limit, err := strconv.Atoi(resp.Header.Get("X-RateLimit-limit")); err != nil {
		c.rateLimit.limit = limit
	}
	if reset, err := strconv.Atoi(resp.Header.Get("X-RateLimit-reset")); err != nil {
		c.rateLimit.reset = time.Unix(int64(reset), 0)
	}

	return resp, nil
}
