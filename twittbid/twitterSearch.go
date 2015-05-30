package twittbid

import (
	"fmt"
	"net/url"
)

//Search is for handle twitter search
func Search(hashtag string) (SearchResponse, error) {
	SetConsumerKey(key)
	SetConsumerSecret(secret)
	api := create(accessTokens, accessTokenSecret)
	fmt.Println(*api.Credentials)
	res, err := api.search(hashtag, nil)
	return res, err
}

//search performs a search based on queryString
func (t TwitterAPI) search(queryString string, v url.Values) (sr SearchResponse, err error) {
	v = CleanValues(v)
	v.Set("q", queryString)
	responseCh := make(chan Response)
	t.queryQueue <- Query{BaseURL + "/search/tweets.json", v, &sr, _GET, responseCh}

	resp := <-responseCh
	err = resp.err
	return sr, err
}

func (t TwitterAPI) performQuery(urlStr string, form url.Values, data interface{}, method int) error {
	switch method {
	case _GET:
		return t.Get(urlStr, form, data)
	case _Post:
		return t.Post(urlStr, form, data)
	default:
		return fmt.Errorf("HTTP method not yet supported")
	}
}
