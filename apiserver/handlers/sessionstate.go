package handlers

import (
	"time"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
)

// SessionState is used to track session state
type SessionState struct {
	BeginAt    time.Time
	ClientAddr string
	User       *users.User
}
