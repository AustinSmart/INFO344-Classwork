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
	if r.Method == "POST" {
		err := r.ParseForm
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
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
		}

		_, err1 := ctx.UserStore.GetByEmail(nu.Email)
		if err1 != nil {
			http.Error(w, "Email already exists in store", http.StatusBadRequest)
		}

		_, err1 = ctx.UserStore.GetByUserName(nu.UserName)
		if err1 != nil {
			http.Error(w, "Username already exists in store", http.StatusBadRequest)
		}

		u, err1 := ctx.UserStore.Insert(&nu)
		if err1 != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
		}
		_, err1 = sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, nil, w)
		if err1 != nil {
			http.Error(w, "Error starting session", http.StatusInternalServerError)
		}

		ju, err1 := json.Marshal(u)
		if err1 != nil {
			http.Error(w, "JSON error", http.StatusInternalServerError)
		}

		_, err1 = w.Write(ju)
		if err1 != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
		}
	}
	if r.Method == "GET" {
		u, err := ctx.UserStore.GetAll()

		ju, err := json.Marshal(u)
		if err != nil {
			http.Error(w, "JSON error", http.StatusInternalServerError)
		}

		_, err = w.Write(ju)
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
		}
	}
}

// SessionsHandler allows existing users to sign in
func (ctx *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Request type must be POST", http.StatusBadRequest)
	}

	c := users.Credentials{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	u, err := ctx.UserStore.GetByEmail(c.Email)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	err = u.Authenticate(c.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	sessions.BeginSession(ctx.SessionKey, ctx.SessionStore, nil, w)

	ju, err1 := json.Marshal(u)
	if err1 != nil {
		http.Error(w, "JSON error", http.StatusInternalServerError)
	}

	_, err1 = w.Write(ju)
	if err1 != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

// SessionsMineHandler allows authenticated users to sign out
func (ctx *Context) SessionsMineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Request type must be DELETE", http.StatusBadRequest)
	}

	sid, err := sessions.EndSession(r, ctx.SessionKey, ctx.SessionStore)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("User %s has been signed out", sid)))
	}
}

// UsersMeHandler responds with the session state
func (ctx *Context) UsersMeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := sessions.GetState(r, ctx.SessionKey, ctx.SessionStore, nil)
	if err != nil {
		http.Error(w, "Error retrieving state", http.StatusInternalServerError)
	}

	//TODO: Respond to the client with the session state's User field, encoded as a JSON object
}

//TODO: Change state accesses to something other than nil. most likely after MongoDB store is implemented.
