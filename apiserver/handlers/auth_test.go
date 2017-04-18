package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

func testUsersCase(t *testing.T) {
	c := Context{
		SessionKey:   "thisisKEY",
		SessionStore: sessions.NewMemStore,
		UserStore:    users.NewMemStore,
	}
	req, err := http.NewRequest("POST", "/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(c.UsersHandler)
	handler.ServeHTTP(rr, req)
}
