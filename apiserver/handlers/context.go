package handlers

import (
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/messages"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	"github.com/info344-s17/challenges-AustinSmart/apiserver/sessions"
)

// Context provides context for http handlers
type Context struct {
	SessionKey    string
	SessionStore  sessions.Store
	UserStore     users.Store
	MessagesStore messages.Store
	Notifier      Notifier
}
