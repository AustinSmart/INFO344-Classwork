package handlers

import (
	"net/http"
	"os/user"

	"encoding/json"

	"fmt"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

const (
	headerContentType   = "Content-Type"
	charsetUTF8         = "charset=utf-8"
	contentTypeJSON     = "application/json"
	contentTypeJSONUTF8 = contentTypeJSON + "; " + charsetUTF8
)

// UsersHandler adds and returns a new user if the request type is POST, returns all users if the request type is GET
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		nu := &users.NewUser{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(nu)
		if err != nil {
			http.Error(w, "Bad JSON:"+err.Error(), http.StatusBadRequest)
		}

		err = nu.Validate()
		if err != nil {
			http.Error(w, "Error validating form: "+err.Error(), http.StatusBadRequest)
			return
		}

		u, err := ctx.UserStore.GetByEmail(nu.Email)
		if err != nil && err != users.ErrUserNotFound {
			http.Error(w, "Error validaitng form: "+err.Error(), http.StatusBadRequest)
			return
		}
		if u != nil {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		u, err = ctx.UserStore.GetByUserName(nu.UserName)
		if err != nil && err != users.ErrUserNotFound {
			http.Error(w, "Error validaitng form: "+err.Error(), http.StatusBadRequest)
			return
		}
		if u != nil {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}

		u, err1 := ctx.UserStore.Insert(nu)
		if err1 != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}
		s := &SessionState{}
		_, err1 = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, s, w)
		if err1 != nil {
			http.Error(w, "Session error: "+err1.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, u)
	case "GET":
		u, err := ctx.UserStore.GetAll()
		if err != nil {
			http.Error(w, "Error retrieving users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		respond(w, u)
	}
}

// SessionsHandler allows existing users to sign in
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request type must be POST", http.StatusBadRequest)
		return
	}

	c := &users.Credentials{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(c)
	if err != nil {
		http.Error(w, "Bad JSON:"+err.Error(), http.StatusBadRequest)
	}

	u, err := ctx.UserStore.GetByEmail(c.Email)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = u.Authenticate(c.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	s := &SessionState{}
	sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, s, w)

	respond(w, u)
}

// SessionsMineHandler allows authenticated users to sign out
func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Request type must be DELETE", http.StatusBadRequest)
		return
	}

	sid, err := sessions.EndSession(r, ctx.SessionKey, ctx.SessionStore)
	if err == nil {
		w.Write([]byte(fmt.Sprintf("User %s has been signed out", sid)))
	} else {
		http.Error(w, "Error signing out", http.StatusInternalServerError)
		return
	}
}

// UsersMeHandler responds with the session state
func (ctx *Context) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	u := user.User{}
	s := SessionState{
		User: &u,
	}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, &s)
	if err != nil {
		http.Error(w, "Error retrieving state:"+err.Error(), http.StatusInternalServerError)
		return
	}
	respond(w, u)
}

func respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add(headerContentType, contentTypeJSONUTF8)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}
