package main

import (
	"fmt"
	"testing"
	"twittbid/twittbid"
)

//TestSearch test that twitter search method is returning valid data
func TestSearch(t *testing.T) {
	res, err := twittbid.Search("golang")
	if err != nil {
		t.Errorf("result returned error: %v", res)
		return
	} else if len(res.Statuses) == 0 {
		t.Errorf("result is 0 --> %v", res)
	}
	for _, tweet := range res.Statuses {
		fmt.Print(tweet.Text)
	}
}
