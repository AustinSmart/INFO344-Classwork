package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"encoding/json"

	"fmt"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

const (
	headerAuthorization = "Authorization"
	schemeBearer        = "Bearer "
)

var c = Context{
	SessionKey:   "thisisKEY",
	SessionStore: sessions.NewMemStore(-1),
	UserStore:    users.NewMemStore(),
}
var nu = users.NewUser{
	Email:        "gh@hopper.com",
	Password:     "monkey123",
	PasswordConf: "monkey123",
	UserName:     "GH",
	FirstName:    "Grace",
	LastName:     "Hopper",
}

var session string

func TestUsersCase(t *testing.T) {
	// Test POST
	handler := http.HandlerFunc(c.UsersHandler)
	rr := httptest.NewRecorder()

	form := url.Values{}
	form.Add("email", nu.Email)
	form.Add("password", nu.Password)
	form.Add("passwordconf", nu.PasswordConf)
	form.Add("username", nu.UserName)
	form.Add("firstname", nu.FirstName)
	form.Add("lastname", nu.LastName)

	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBufferString(form.Encode()))
	req.Header.Add(headerContentType, "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, rr.Code)
	}

	ct := rr.Header().Get(headerContentType)
	if ct != contentTypeJSONUTF8 {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", contentTypeJSONUTF8, ct)
	}

	if nil == rr.Body || 0 == rr.Body.Len() {
		t.Errorf("handler returned empty response body")
	}

	u := users.User{}
	json.NewDecoder(rr.Body).Decode(&u)
	if len(u.ID) == 0 {
		t.Errorf("returned user does not have an id")
	}
	if u.Email != nu.Email {
		t.Errorf("returned email does not match: expected %s; got %s", nu.Email, u.Email)
	}
	if u.UserName != nu.UserName {
		t.Errorf("returned username does not match: expected %s; got %s", nu.UserName, u.UserName)
	}
	if u.FirstName != nu.FirstName {
		t.Errorf("returned first name does not match: expected %s; got %s", nu.FirstName, u.FirstName)
	}
	if u.LastName != nu.LastName {
		t.Errorf("returned last name does not match: expected %s; got %s", nu.LastName, u.LastName)
	}

	// Test GET
	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, rr.Code)
	}

	ct = rr.Header().Get(headerContentType)
	if ct != contentTypeJSONUTF8 {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", contentTypeJSONUTF8, ct)
	}

	if nil == rr.Body || 0 == rr.Body.Len() {
		t.Errorf("handler returned empty response body")
	}

	users := []users.User{}
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 1 {
		t.Errorf("incorrect number of users returned by GetAll: expected `%d`; got `%d`", 1, len(users))
	}
}

func TestSessionsCase(t *testing.T) {
	handler := http.HandlerFunc(c.SessionsHandler)
	rr := httptest.NewRecorder()
	// Test incorrect request type
	req, err := http.NewRequest("GET", "/v1/sessions", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusBadRequest, rr.Code)
	}

	// Test correct request type
	rr = httptest.NewRecorder()
	form := url.Values{}
	form.Add("email", nu.Email)
	form.Add("password", nu.Password)

	req, err = http.NewRequest("POST", "/v1/sessions", bytes.NewBufferString(form.Encode()))
	req.Header.Add(headerContentType, "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, rr.Code)
	}

	ct := rr.Header().Get(headerContentType)
	if ct != contentTypeJSONUTF8 {
		t.Errorf("incorrect Content-Type response header: expected %s; got %s", contentTypeJSONUTF8, ct)
	}

	if nil == rr.Body || 0 == rr.Body.Len() {
		t.Errorf("handler returned empty response body")
	}

	u := users.User{}
	json.NewDecoder(rr.Body).Decode(&u)
	if u.Email != nu.Email {
		t.Errorf("returned user does not match: expected %s; got %s", nu.UserName, u.UserName)
	}
	if len(rr.Header().Get(headerAuthorization)) == 0 {
		t.Error("response does not have session header")
	}
	session = rr.Header().Get(headerAuthorization)
}

func TestSessionsMineCase(t *testing.T) {
	handler := http.HandlerFunc(c.SessionsMineHandler)
	rr := httptest.NewRecorder()
	// Test incorrect request type
	req, err := http.NewRequest("GET", "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusBadRequest, rr.Code)
	}

	// Test correct request type
	rr = httptest.NewRecorder()

	req, err = http.NewRequest("DELETE", "/v1/sessions/mine", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(headerAuthorization, session)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, rr.Code)
	}

	if nil == rr.Body || 0 == rr.Body.Len() {
		t.Errorf("handler returned empty response body")
	}
}

func TestUsersMeCase(t *testing.T) {
	s := SessionState{}
	rr := httptest.NewRecorder()
	sid, err := sessions.BeginSession(c.SessionKey, c.SessionStore, &s, rr)
	handler := http.HandlerFunc(c.UsersMeHandler)
	auth := rr.Header().Get(headerAuthorization)
	if !strings.HasPrefix(auth, schemeBearer) {
		t.Errorf("Authorization header value does not start with %s: got %s\n", schemeBearer, auth)
	}

	rr = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v1/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(headerAuthorization, schemeBearer+sid.String())

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected `%d` but got `%d`\n", http.StatusOK, rr.Code)
	}

	u := users.User{}
	json.NewDecoder(rr.Body).Decode(&u)
	fmt.Println(rr.Body)
	if u.Email != nu.Email {
		t.Errorf("returned user does not match: expected %s; got %s", nu.UserName, u.UserName)
	}
}
