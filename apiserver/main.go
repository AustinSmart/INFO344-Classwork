package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	redis "gopkg.in/redis.v5"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/handlers"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/middleware"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/messages"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

const (
	apiRoot = "/v1/"

	apiSummary = apiRoot + "summary"

	apiUsers        = apiRoot + "users"
	apiSessions     = apiRoot + "sessions"
	apiSessionsMine = apiSessions + "/mine"
	apiUsersMe      = apiUsers + "/me"

	apiChannels        = apiRoot + "channels"
	apiSpecificChannel = apiChannels + "/"
	apiMessages        = apiRoot + "messages"
	apiSpecificMessage = apiMessages + "/"

	defaultPort = "443"
)

func main() {
	//Environment variables
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
	defer mongoSession.Close()

	mongoUsersStore := users.NewMongoStore(mongoSession, "users-db", "users")
	mongoMessagesStore := messages.NewMongoStore(mongoSession, "messages-db", "messages", "channels")

	ctx := handlers.Context{
		SessionKey:    sessionKey,
		SessionStore:  sessions.NewRedisStore(client, -1),
		UserStore:     mongoUsersStore,
		MessagesStore: mongoMessagesStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiSummary, handlers.SummaryHandler)
	mux.HandleFunc(apiUsers, ctx.UsersHandler)
	mux.HandleFunc(apiSessions, ctx.SessionsHandler)
	mux.HandleFunc(apiSessionsMine, ctx.SessionsMineHandler)
	mux.HandleFunc(apiUsersMe, ctx.UsersMeHandler)
	mux.HandleFunc(apiChannels, ctx.ChannelsHandler)
	mux.HandleFunc(apiSpecificChannel, ctx.SpecificChannelHandler)
	mux.HandleFunc(apiMessages, ctx.MessagesHandler)
	mux.HandleFunc(apiSpecificMessage, ctx.SpecificMessageHandler)

	http.Handle(apiRoot, middleware.Adapt(mux, middleware.CORS("", "", "", "")))

	//General channel for all users
	mongoSession.DB(mongoMessagesStore.DatabaseName).C(mongoMessagesStore.ChannelsCollectionName).RemoveAll(bson.M{"name": "General"})
	nc := &messages.NewChannel{
		Name:        "General",
		Description: "A public channel for all users of Taut",
		Private:     false,
	}
	ctx.MessagesStore.InsertChannel("admin", nc)

	//Start the server
	fmt.Printf("listening on %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, nil))
}
