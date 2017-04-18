package handlers

import (
	"net/http"

	"encoding/json"

	"fmt"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

// UsersHandler adds and returns a new user if the request type is POST, returns all users if the request type is GET
func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		err := r.ParseForm
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		nu := users.NewUser{
			Email:        r.PostFormValue("email"),
			Password:     r.PostFormValue("password"),
			PasswordConf: r.PostFormValue("passwordconf"),
			UserName:     r.PostFormValue("username"),
			FirstName:    r.PostFormValue("firstname"),
			LastName:     r.PostFormValue("lastname"),
		}

		err = nu.Validate
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		_, err1 := ctx.UserStore.GetByEmail(nu.Email)
		if err1 != nil {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		_, err1 = ctx.UserStore.GetByUserName(nu.UserName)
		if err1 != nil {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}

		u, err1 := ctx.UserStore.Insert(&nu)
		if err1 != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}
		s := &SessionState{}
		_, err1 = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, s, w)
		if err1 != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}
		marshalResponse(w, u)
	case "GET":
		u, err := ctx.UserStore.GetAll()
		if err != nil {
			http.Error(w, "Error retrieving users", http.StatusInternalServerError)
			return
		}
		marshalResponse(w, u)
	}
}

// SessionsHandler allows existing users to sign in
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request type must be POST", http.StatusBadRequest)
		return
	}

	c := users.Credentials{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
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

	marshalResponse(w, u)
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
	}
}

// UsersMeHandler responds with the session state
func (ctx *Context) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	s := &SessionState{}
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, s)
	if err != nil {
		http.Error(w, "Error retrieving state", http.StatusInternalServerError)
		return
	}

	marshalResponse(w, s.User)
}

func marshalResponse(w http.ResponseWriter, data interface{}) {
	jdata, err1 := json.Marshal(data)
	if err1 != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
		return
	}

	_, err1 = w.Write(jdata)
	if err1 != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
