package scrapingAPI

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/azer/atlas"
	. "github.com/azer/debug"
	"github.com/franela/goreq"
	"math"
	"net/http"
	"net/url"
	"time"
)

type Query struct {
	Selector string
	Node     string
}

type Options struct {
	URL      string
	Callback string
	Query    map[string]Query
	StartTS  int64
}

type Result struct {
	Key      string
	Selector string
	Value    string
	Node     string
}

type Results map[string]Result

type Stats struct {
	Scraping             int
	Scraped              int
	Now                  int64
	ActiveRequest        int
	FailedRequest        int
	SuccessDelivery      int
	FailedDelivery       int
	AvgDeliverTime       int
	ShortestDeliveryTime int
	LongestDeliveryTime  int
	LastError            string
}

var (
	Scraping             = 0
	Scraped              = 0
	ActiveRequest        = 0
	FailedRequest        = 0
	LastError            = ""
	SuccessDelivery      = 0
	FailedDelivery       = 0
	TotalDeliverTime     = 0
	AvgDeliverTime       = 0
	ShortestDeliveryTime = 9999
	LongestDeliveryTime  = 0
)

var Server = atlas.New(atlas.Map{
	"/scrape": Scrape,
	"/stats":  GetStats,
})

func GetStats(request *atlas.Request) *atlas.Response {
	return atlas.Success(Stats{
		Scraping,
		Scraped,
		now(),
		ActiveRequest,
		FailedRequest,
		SuccessDelivery,
		FailedDelivery,
		AvgDeliverTime,
		ShortestDeliveryTime,
		LongestDeliveryTime,
		LastError,
	})
}

func Scrape(request *atlas.Request) *atlas.Response {
	opts := &Options{}
	err := request.JSONPost(&opts)

	opts.StartTS = now()

	Debug("Scraping %v", err)

	if err != nil {
		Debug("Failed to parse the JSON body: %v", err)
		return atlas.Error(500, err)
	}

	Scraping++

	if len(opts.Callback) > 0 {
		go Deliver(opts)
		return atlas.Success("Results will be posted to " + opts.Callback)
	}

	result, err := Select(opts)

	Scraping--
	Scraped++

	if err != nil {
		return atlas.Error(500, err)
	}

	return atlas.Success(result)
}

func Select(opts *Options) (result Results, err error) {
	var doc *goquery.Document
	result = make(Results)

	ActiveRequest++

	if doc, err = goquery.NewDocument(opts.URL); err != nil {
		ActiveRequest--
		FailedRequest++
		return nil, err
	}

	for key, query := range opts.Query {
		el := doc.Find(query.Selector)

		var value string

		if query.Node == "text" {
			value = el.Text()
		}

		if query.Node == "html" {
			value, _ = el.Html()
		}

		if len(query.Node) > 5 && query.Node[0:5] == "attr:" {
			value, _ = el.Attr(query.Node[5:])
		}

		result[key] = Result{
			Key:      key,
			Value:    value,
			Selector: query.Selector,
			Node:     query.Node,
		}
	}

	ActiveRequest--

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

	if err != nil {
		DeliverError(opts, "Failed to parse and extract the data.")
		return
	}

	Debug("Posting results to %s", opts.Callback)

	res, err := goreq.Request{
		Method:      "POST",
		Uri:         opts.Callback,
		Body:        result,
		Accept:      "application/json",
		ContentType: "application/json",
		Timeout:     10000 * time.Millisecond,
	}.Do()

	Scraping--
	Scraped++
	elapsed := int(now() - opts.StartTS)
	TotalDeliverTime = TotalDeliverTime + elapsed
	AvgDeliverTime = TotalDeliverTime / Scraped

	if elapsed < ShortestDeliveryTime {
		ShortestDeliveryTime = elapsed
	}

	if elapsed > LongestDeliveryTime {
		LongestDeliveryTime = elapsed
	}

	if err != nil {
		FailedDelivery++
		Debug("Unable to post to %s. Error: %v", opts.Callback, err)
		LastError = err.Error()
		return
	}

	SuccessDelivery++

	defer res.Body.Close()
}

func now() int64 {
	return int64(math.Floor(float64(time.Now().UnixNano()) / 1000000))
}
