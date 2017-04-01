package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/handlers"
)

const defaultPort = "80"

const (
	apiRoot    = "/v1/"
	apiSummary = apiRoot + "summary"
)

//main is the main entry point for this program
func main() {
	//get the host and port from environment variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "443"
	}
	host := os.Getenv("HOST")
	addr := fmt.Sprintf("%s:%s", host, port)

	//get the TLS key and cert paths from environment variables
	//this allows us to use a self-signed cert/key during development
	//and the Let's Encrypt cert/key in production
	tlsKeyPath := "/etc/info344.austinsmart.com.key"
	tlsCertPath := "/etc/info344.austinsmart.com.pem"

	http.HandleFunc(apiSummary, handlers.SummaryHandler)

	//start the server
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
