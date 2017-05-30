package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

//ChatBotHandler lets users ask NLP based questions
func (ctx *Context) ChatBotHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// if r.Method != "POST" {
	// 	http.Error(w, "Request type must be POST", http.StatusBadRequest)
	// 	return
	// }

	u := users.User{}
	s := SessionState{
		User: &u,
	}
	sid, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &s)
	if err != nil || sid.String() == "" {
		http.Error(w, "Error retrieving state:"+err.Error(), http.StatusInternalServerError)
		return
	}

	chatBotAddr := os.Getenv("CHATBOTADDR")
	if len(chatBotAddr) == 0 {
		log.Fatal("You must supply a value for CHATBOTADDR")
	}

	service := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = chatBotAddr
			j, _ := json.Marshal(u)
			r.Header.Add("User", string(j))
		},
	}

	service.ServeHTTP(w, r)
}
