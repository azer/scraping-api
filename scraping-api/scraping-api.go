package main

import (
	"flag"
	"github.com/azer/scraping-api"
	"fmt"
)

var port int

func main() {
	flag.IntVar(&port, "port", 8080, "Port to serve on")
	flag.Parse()
	scrapingAPI.Server.Start(fmt.Sprintf(":%d", port))
}
