## scraping-api

[API server](http://github.com/azer/atlas) for [goquery](https://github.com/PuerkitoBio/goquery)

## Install

```bash
$ go get github.com/azer/scraping-api
```

## Usage

Start the server:

```bash
$ scraping-api
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
