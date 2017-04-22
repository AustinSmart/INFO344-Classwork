package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/handlers"
)

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
)

func main() {
	//get the host and port from environment variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "443"
	}
	host := os.Getenv("HOST")
	addr := fmt.Sprintf("%s:%s", host, port)

	//get cloudflare origin cert and key
	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	http.HandleFunc(apiSummary, handlers.SummaryHandler)

	//start the server
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
