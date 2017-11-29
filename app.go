package main

import (
	"flag"
	"github.com/google/gops/agent"
	"log"
	"net/http"

	"github.com/tokopedia/gosample/hello"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true, // automatically closes on os.Interrupt
	}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule()

	http.HandleFunc("/hello", hwm.SayHelloWorld)
	
	//FOR TRAINING
	//http.HandleFunc("/private", hwm.thisIsPrivate)
	http.HandleFunc("/public", hwm.ThisIsPublic)

	//log.Println(hwm.private)
	log.Println(hwm.Public)

	go logging.StatsLog()

	log.Fatal(grace.Serve(":9000", nil))
}
