package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AustinSmart/challenges-AustinSmart/apiserver/handlers"
)

const defaultPort = "80"

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
)

//main is the main entry point for this program
func main() {
	//read and use the following environment variables
	//when initializing and starting your web server
	// PORT - port number to listen on for HTTP requests (if not set, use defaultPort)
	// HOST - host address to respond to (if not set, leave empty, which means any host)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("HOST")

	//add your handlers.SummaryHandler function as a handler
	//for the apiSummary route
	//HINT: https://golang.org/pkg/net/http/#HandleFunc
	http.HandleFunc("/summary", handlers.SummaryHandler)

	//start your web server and use log.Fatal() to log
	//any errors that occur if the server can't start
	//HINT: https://golang.org/pkg/net/http/#ListenAndServe
	addr := host + ":" + port
	fmt.Printf("server is listening at %s:%s...\n", host, port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
