## scraping-api

Go [API server](http://github.com/azer/atlas) for scraping data with [CSS selectors](https://github.com/PuerkitoBio/goquery)

## Install

```bash
$ go get github.com/azer/scraping-api/scraping-api
```

## Usage

Start the server:

```bash
$ scraping-api -port 1234
```

And send JSON-Post requests to scrape data:

```bash
$ curl -X POST -d '{"url":"http://azer.io", "query": { "title": { "selector": "h1:first-child", "node":"text" } }}' http://localhost:8080
```

A request like above will output:

```
{
  "result": {
      "title": {
          "Key": "title",
          "Selector": "h1:first-child",
          "Value": "Azer Koçulu",
          "Node": "text"
      }
  },
  "ok": true
}%
```

You can get attribute values by choosing `attr:?` as node value in the query:

```
$ curl -X POST -d '{"url":"http://azer.io", "query": { "first-article": { "selector": ".articles li:first-child a", "node":"attr:href" } }}' http://localhost:8080
```

Optionally, results delivered to a callback URL by specifying the "callback" parameter:

```
$ curl -X POST -d '{"url":"http://azer.io", "callback":"http://localhost/save-results", "query": { "title": { "selector": "h1:first-child", "node":"text" } }}' http://localhost:8080
```

![](http://distilleryimage5.ak.instagram.com/51eb9256ba2611e3a63112f56a54141d_6.jpg)
