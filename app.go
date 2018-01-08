package main

import (
	"flag"
	"log"
	"net/http"

  "github.com/tokopedia/gosample/hello"
	grace "gopkg.in/tokopedia/grace.v1"
	logging "gopkg.in/tokopedia/logging.v1"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

  debug("app started") // message will not appear unless run with -debug switch

  hwm := hello.NewHelloWorldModule()

	http.HandleFunc("/hello", hwm.SayHelloWorld)
	http.HandleFunc("/set", hwm.SetRedis)
	http.HandleFunc("/get", hwm.GetRedis)
	go logging.StatsLog()

	log.Fatal(grace.Serve(":9000", nil))
}
