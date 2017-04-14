package handlers

import (
	"os/user"
	"time"
)

// SessionState is used to track session state
type SessionState struct {
	BeginAt    time.Time
	ClientAddr string
	User       *user.User
}
