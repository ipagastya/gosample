package main

import (
	"flag"
	"log"
	"net/http"

  	"github.com/tokopedia/gosample/website"
	grace "gopkg.in/tokopedia/grace.v1"
	logging "gopkg.in/tokopedia/logging.v1"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

  	debug("app started") // message will not appear unless run with -debug switch

  	web := website.NewWebsiteModule()

	http.HandleFunc("/index", web.RenderWebpage)
	go logging.StatsLog()

	log.Fatal(grace.Serve(":9000", nil))
}
