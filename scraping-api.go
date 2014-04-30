package scrapingAPI

import (
	"github.com/azer/atlas"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"encoding/json"
	"strings"
	"fmt"
	. "github.com/azer/debug"
)

type Query struct {
  Selector string
  Node string
}

type Options struct {
  URL string
	Callback string
  Query map[string]Query
}

type Result struct {
  Key string
  Selector string
  Value string
	Node string
}

type Results map[string]Result

var Server = atlas.New(atlas.Map{
	"/": Scrape,
})

func Scrape(request *atlas.Request) *atlas.Response {
	opts := &Options{}
	err := request.JSONPost(&opts)

	if err != nil {
		return atlas.Error(500, err)
	}

	if len(opts.Callback) > 0 {
		go Deliver(opts)
		return atlas.Success("Results will be posted to " + opts.Callback)
	}

	result, err := Select(opts)

	if err != nil {
		return atlas.Error(500, err)
	}

	return atlas.Success(result)
}

func Select(opts *Options) (result Results, err error) {
	var doc *goquery.Document
	result = make(Results)

	if doc, err = goquery.NewDocument(opts.URL); err != nil {
		return nil, err
	}

	for key, query := range opts.Query {
		el := doc.Find(query.Selector)

		var value string

		if query.Node == "text" {
			value = el.Text()
		}

		if query.Node == "html" {
			value, _ := el.Html()
			value = strings.Replace(value, "\u003c", "<", -1)
			value = strings.Replace(value, "\u003e", ">", -1)
			value = strings.Replace(value, "<br>", "\n", -1)
			value = strings.Replace(value, "<br/>", "\n", -1)
			value = strings.Replace(value, "<br />", "\n", -1)
		}

		if len(query.Node) > 5 && query.Node[0:5] == "attr:" {
			value, _ = el.Attr(query.Node[5:])
		}

		result[key] = Result{
			Key: key,
			Value: value,
   		Selector: query.Selector,
   		Node: query.Node,
		}
	}

	return result, nil
}

func DeliverError(opts *Options, err string) {
	Debug("Request to %s failed. Error: %s", opts.Callback, err)

	data := url.Values{}
	data.Set(fmt.Sprintf("{ \"error\": \"%s\" }", err), "")
	_, derr := http.PostForm(opts.Callback, data)

	if derr != nil {
		Debug("Failed to post the error message '%s' to %s", err, opts.URL)
	}
}

func Deliver(opts *Options) {
	Debug("Results of %s will be delivered to %s", opts.URL, opts.Callback)

	result, err := Select(opts)

	data := url.Values{}

	if err != nil {
		DeliverError(opts, "Failed to parse and extract the data.")
		return
	}

	strResult, err := json.Marshal(result)

	if err != nil {
		DeliverError(opts, "Apologizes, scraping-api screwed up.")
		return
	}

	data.Set(string(strResult), "")

	Debug("Posting results to %s", opts.Callback)
	_, err = http.PostForm(opts.Callback, data)

	if err != nil {
		Debug("Unable to post to %s", opts.Callback)
	}
}
