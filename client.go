package zeit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type rateLimit struct {
	limit     int
	remaining int
	reset     time.Time
	mutex     sync.Mutex
}

type Client struct {
	token      string
	rootUrl    string
	httpClient HttpClient
	rateLimit  *rateLimit
	team       string
}

var rateLimits = make(map[string]*rateLimit)

//go:generate mockgen -destination=mocks/mock_http_client.go -package=mocks github.com/kochie/zeit-api-go HttpClient

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewClient will create a new zeit client to apply api request to. Note that the team is defaulted to nothing, if you
// want to update the team then use Team.
func NewClient(token string) *Client {
	var rl *rateLimit
	if val, ok := rateLimits[token]; ok {
		rl = val
	} else {
		rl = &rateLimit{1, 1, time.Now(), sync.Mutex{}}
	}
	return &Client{
		token,
		"https://api.zeit.co",
		&http.Client{},
		rl,
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
	mutex := c.rateLimit.mutex

	// do a check to see if the rate limit has been hit, if so wait until a request can be sent again
doRequest:
	mutex.Lock()
	remaining := c.rateLimit.remaining
	reset := c.rateLimit.reset
	mutex.Unlock()
	now := time.Now()
	if remaining == 0 && now.Before(reset) {
		d := reset.Sub(now)
		log.Printf("Zeit rate limit hit, waiting for %s", d.String())
		time.Sleep(d)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		rateLimitError := RateLimitError{}
		err := json.NewDecoder(resp.Body).Decode(&struct {
			Error *RateLimitError `json:"error"`
		}{&rateLimitError})
		if err != nil {
			return nil, err
		}
		mutex.Lock()
		c.rateLimit.remaining = rateLimitError.Limit.Remaining
		c.rateLimit.reset = time.Unix(rateLimitError.Limit.Reset, 0)
		c.rateLimit.limit = rateLimitError.Limit.Total
		mutex.Unlock()

		goto doRequest
	}

	mutex.Lock()
	if RateLimitRemaining := resp.Header.Get("X-RateLimit-Remaining"); RateLimitRemaining != "" {
		if remaining, err := strconv.ParseInt(RateLimitRemaining, 10, 32); err != nil {
			c.rateLimit.remaining = int(remaining)
		}
	}
	if RateLimitLimit := resp.Header.Get("X-RateLimit-Limit"); RateLimitLimit != "" {
		if limit, err := strconv.ParseInt(RateLimitLimit, 10, 32); err != nil {
			c.rateLimit.limit = int(limit)
		}
	}
	if RateLimitReset := resp.Header.Get("X-RateLimit-Reset"); RateLimitReset != "" {
		if reset, err := strconv.ParseInt(RateLimitReset, 10, 64); err != nil {
			c.rateLimit.reset = time.Unix(reset, 0)
		}
	}
	mutex.Unlock()

	return resp, nil
}
