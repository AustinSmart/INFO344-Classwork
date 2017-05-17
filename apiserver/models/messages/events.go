package messages

import "github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"

type ChannelEvent struct {
	Type    string   `json:"type"`
	Channel *Channel `json:"channel"`
}

type MessageEvent struct {
	Type    string   `json:"type"`
	Message *Message `json:"message`
}

type UserEvent struct {
	Type string      `json:"type"`
	User *users.User `json:"user"`
}
