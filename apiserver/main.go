package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
	redis "gopkg.in/redis.v5"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/handlers"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/middleware"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

const (
	apiRoot         = "/v1/"
	apiSummary      = apiRoot + "summary"
	apiUsers        = apiRoot + "users"
	apiSessions     = apiRoot + "sessions"
	apiSessionsMine = apiSessions + "/mine"
	apiUsersMe      = apiUsers + "/me"
	defaultPort     = "443"
)

func main() {
	//environment variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	addr := fmt.Sprintf("%s:%s", host, port)

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")

	sessionKey := os.Getenv("SESSIONKEY")
	redisAddr := os.Getenv("REDISADDR")
	dbAddr := os.Getenv("DBADDR")

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	mongoSession, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatalf("error dialing mongo: %v", err)
	}

	mongoStore := users.NewMongoStore(mongoSession, "users-db", "users")

	ctx := handlers.Context{
		SessionKey:   sessionKey,
		SessionStore: sessions.NewRedisStore(client, -1),
		UserStore:    mongoStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiUsers, ctx.UsersHandler)
	mux.HandleFunc(apiSessions, ctx.SessionsHandler)
	mux.HandleFunc(apiSessionsMine, ctx.SessionsMineHandler)
	mux.HandleFunc(apiUsersMe, ctx.UsersMeHandler)
	mux.HandleFunc(apiSummary, handlers.SummaryHandler)

	http.Handle(apiRoot, middleware.Adapt(mux, middleware.CORS("", "", "", "")))

	//start the server
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
