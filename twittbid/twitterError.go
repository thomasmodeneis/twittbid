package twittbid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//APIError struct
type APIError struct {
	StatusCode int
	Header     http.Header
	Body       string
	Decoded    TwitterErrorResponse
	URL        *url.URL
}

//Error the error interface
func (aerr APIError) Error() string {
	return fmt.Sprintf("Get %s returned status %d, %s", aerr.URL, aerr.StatusCode, aerr.Body)
}

//createError
func createError(resp *http.Response) *APIError {
	p, _ := ioutil.ReadAll(resp.Body)
	var twitterErrorResp TwitterErrorResponse
	_ = json.Unmarshal(p, &twitterErrorResp)
	return &APIError{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(p),
		Decoded:    twitterErrorResp,
		URL:        resp.Request.URL,
	}
}

//CheckRateLimit validate response for if status code 429 or 130
func (aerr *APIError) CheckRateLimit() (isRateLimitError bool, nextWindow time.Time) {
	if aerr.StatusCode == 429 {
		if reset := aerr.Header.Get("X-Rate-Limit-Reset"); reset != "" {
			if resetUnix, err := strconv.ParseInt(reset, 10, 64); err == nil {
				resetTime := time.Unix(resetUnix, 0)
				// Reject any time greater than an hour away
				if resetTime.Sub(time.Now()) > time.Hour {
					return true, time.Now().Add(15 * time.Minute)
				}

				return true, resetTime
			}
		}
	}

	return false, time.Time{}
}
