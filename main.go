package main

import (
	"encoding/json"
	"net/http"
	"twittbid/twittbid"

	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()

	m.Get("/", func() string {
		return "TwittBid REST API, usage --> search/golang"
	})

	m.Get("/search/:hashtag", listByHashtag)

	m.RunOnAddr(":3006")

}

func listByHashtag(res http.ResponseWriter, req *http.Request, params martini.Params) {

	//handle params
	hashtag := params["hashtag"]

	defer func() {
		if r := recover(); r != nil {
			http.Error(res, "", http.StatusNotFound)
			return
		}
	}()

	//execute search
	result, err := twittbid.Search(hashtag)
	if err != nil {
		http.Error(res, "", http.StatusNotFound)
		return
	}

	//process json
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	//render
	res.Header().Set("Content-Type", "application/json")
	res.Write(js)
}
