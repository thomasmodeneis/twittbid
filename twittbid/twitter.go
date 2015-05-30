package twittbid

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/garyburd/go-oauth/oauth"
)

//SetConsumerKey this your API key
func SetConsumerKey(consumerKey string) {
	OauthClient.Credentials.Token = consumerKey
}

//SetConsumerSecret from your API secret
func SetConsumerSecret(consumerSecret string) {
	OauthClient.Credentials.Secret = consumerSecret
}

//create TwitterAPI  --> constructor
func create(accessToken string, accessTokenSecret string) *TwitterAPI {
	queue := make(chan Query)
	c := &TwitterAPI{
		Credentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: accessTokenSecret,
		},
		queryQueue:           queue,
		bucket:               nil,
		returnRateLimitError: false,
		HTTPClient:           http.DefaultClient,
	}
	go c.runQuery()
	return c
}

//runQuery will execute the query and stream result to q.responseCh
func (t *TwitterAPI) runQuery() {
	for q := range t.queryQueue {
		url := q.url
		form := q.form
		data := q.data
		method := q.method

		responseCh := q.responseCh

		if t.bucket != nil {
			<-t.bucket.SpendToken(1)
		}

		err := t.performQuery(url, form, data, method)

		// check for rate-limiting error
		if err != nil {
			if apiErr, ok := err.(*APIError); ok {
				if isRateLimitError, nextWindow := apiErr.CheckRateLimit(); isRateLimitError && !t.returnRateLimitError {
					log.Fatal(apiErr.Error())

					go func() {
						t.queryQueue <- q
					}()

					delay := nextWindow.Sub(time.Now())
					<-time.After(delay)

					// clean bucket
					if t.bucket != nil {
						t.bucket.Drain()
					}

					continue
				}
			}
		}

		responseCh <- Response{data, err}
	}
}

//CleanValues method
func CleanValues(v url.Values) url.Values {
	if v == nil {
		return url.Values{}
	}
	return nil
}

//Get http method
func (t TwitterAPI) Get(urlStr string, form url.Values, data interface{}) error {
	resp, err := OauthClient.Get(t.HTTPClient, t.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return DecodeResponse(resp, data)
}

//Post http method
func (t TwitterAPI) Post(urlStr string, form url.Values, data interface{}) error {
	resp, err := OauthClient.Post(t.HTTPClient, t.Credentials, urlStr, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return DecodeResponse(resp, data)
}

//DecodeResponse decode http response
func DecodeResponse(resp *http.Response, data interface{}) error {
	if resp.StatusCode != 200 {
		return createError(resp)
	}
	return json.NewDecoder(resp.Body).Decode(data)
}
