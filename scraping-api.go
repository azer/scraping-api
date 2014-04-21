package main

import (
	"github.com/azer/atlas"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type Query struct {
  Selector string
  Node string
}

type Options struct {
  URL string
  Query map[string]Query
}

type Result struct {
  Key string
  Selector string
  Value string
	Node string
}

type Results map[string]Result

var api = atlas.New(atlas.Map{
	"/": Scrape,
})

func main() {
	api.Start(":8080")
}

func Scrape(request *atlas.Request) *atlas.Response {
	opts := &Options{}
	err := request.JSONPost(&opts)

	if err != nil {
		return atlas.Error(500, err)
	}

	var doc *goquery.Document
	var result = make(Results)

	if doc, err = goquery.NewDocument(opts.URL); err != nil {
		return atlas.Error(500, err)
	}

	for key, query := range opts.Query {
		el := doc.Find(query.Selector)

		var value string

		if query.Node == "text" {
			value = el.Text()
		}

		if query.Node == "html" {
			value, err = el.Html()
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

	return atlas.Success(result)
}
