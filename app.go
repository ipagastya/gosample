package main

import (
	"flag"
	"fmt"
	"github.com/google/gops/agent"
	"log"
	"net/http"

	"github.com/tokopedia/gosample/hello"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)


func logEndpoint(h http.HandlerFunc) http.HandlerFunc {

	var count = 0

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var callx string
    	count++

    	if count == 11 || count == 12 || count == 13 {
    		callx = fmt.Sprintf("%dth", count)
    	} else if count % 10 == 1 {
    		callx = fmt.Sprintf("%dst", count)
    	} else if count % 10 == 2 {
    		callx = fmt.Sprintf("%dnd", count)
    	} else if count % 10 == 3 {
    		callx = fmt.Sprintf("%drd", count)
    	} else {
    		callx = fmt.Sprintf("%dth", count)
    	}
    	log.Printf("This is %s call of this endpoint", callx)

    	h.ServeHTTP(w, r)
  })
}

func haii(h http.HandlerFunc) http.HandlerFunc {
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Haiii~~")
		
		h.ServeHTTP(w, r)
	})
}

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

	http.HandleFunc("/hello", haii(logEndpoint(hwm.SayHelloWorld)))
	go logging.StatsLog()

	log.Fatal(grace.Serve(":9000", nil))
}
