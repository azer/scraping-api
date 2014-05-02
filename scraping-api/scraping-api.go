package main

import (
	"flag"
	"fmt"
	"github.com/azer/scraping-api"
)

var port int

func main() {
	flag.IntVar(&port, "port", 8080, "Port to serve on")
	flag.Parse()
	scrapingAPI.Server.Start(fmt.Sprintf(":%d", port))
}
