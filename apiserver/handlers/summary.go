package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"golang.org/x/net/html"
)

//openGraphPrefix is the prefix used for Open Graph meta properties
const openGraphPrefix = "og:"

//openGraphProps represents a map of open graph property names and values
type openGraphProps map[string]string

func getPageSummary(url string) (openGraphProps, error) {
	//Get the URL
	//If there was an error, return it
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	//ensure that the response body stream is closed eventually
	//HINTS: https://gobyexample.com/defer
	//https://golang.org/pkg/net/http/#Response
	defer resp.Body.Close()

	//if the response StatusCode is >= 400
	//return an error, using the response's .Status
	//property as the error message
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	//if the response's Content-Type header does not
	//start with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return nil, errors.New("Recieved " + resp.Header.Get("Content-Type") + " expected text/html")
	}

	//create a new openGraphProps map instance to hold
	//the Open Graph properties you find
	//(see type definition above)
	openGraphMap := make(openGraphProps)

	//tokenize the response body's HTML and extract
	//any Open Graph properties you find into the map,
	//using the Open Graph property name as the key, and the
	//corresponding content as the value.
	//strip the openGraphPrefix from the property name before
	//you add it as a new key, so that the key is just `title`
	//and not `og:title` (for example).
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		// token type
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return nil, errors.New(tokenizer.Token().String())
		}
		token := tokenizer.Token()
		if token.Attr != nil {
			switch tokenType {
			case html.SelfClosingTagToken:
				if strings.HasPrefix(token.Attr[0].Val, openGraphPrefix) {
					attr := token.Attr[0].Val
					attrNoPrefix := strings.TrimPrefix(attr, openGraphPrefix)
					value := token.Attr[1].Val
					openGraphMap[attrNoPrefix] = value
				}
			}
		}
	}

	//HINTS: http://golang-examples.tumblr.com/post/47426518779/parse-html
	//https://godoc.org/golang.org/x/net/html

	return openGraphMap, nil
}

//SummaryHandler fetches the URL in the `url` query string parameter, extracts
//summary information about the returned page and sends those summary properties
//to the client as a JSON-encoded object.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	//Add the following header to the response
	//   Access-Control-Allow-Origin: *
	//this will allow JavaScript served from other origins
	//to call this API
	w.Header().Add("Access-Control-Allow-Origin", "*")

	//get the `url` query string parameter
	//if you use r.FormValue() it will also handle cases where
	//the client did POST with `url` as a form field
	//HINT: https://golang.org/pkg/net/http/#Request.FormValue
	url := r.FormValue("url")

	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	//HINT: https://golang.org/pkg/net/http/#Error
	if url == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	} else {

		//call getPageSummary() passing the requested URL
		//and holding on to the returned openGraphProps map
		//(see type definition above)
		openGraphMap, err := getPageSummary(url)

		//if you get back an error, respond to the client
		//with that error and an http.StatusBadRequest code
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		//otherwise, respond by writing the openGrahProps
		//map as a JSON-encoded object
		//add the following headers to the response before
		//you write the JSON-encoded object:
		//   Content-Type: application/json; charset=utf-8
		//this tells the client that you are sending it JSON
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		jsonString, err := json.Marshal(openGraphMap)
		w.Write(jsonString)
		w := httptest.NewRecorder()
		log.Printf("%d - %s", w.Code, w.Body.String())
	}
}
