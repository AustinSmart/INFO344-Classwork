package messages

import "github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"

//MessageID defines the type for message ID's
type MessageID string

//Message represents a message in a channel
type Message struct {
	ID        MessageID    `json:"id" bson:"_id"`
	ChannelID ChannelID    `json:"channelID"`
	Body      string       `json:"body"`
	CreatedAt string       `json:"createdAt"`
	CreatorID users.UserID `json:"creatorID"`
	EditedAt  string       `json:"editedAt"`
}

//NewMessage represnts a new message to be added to a channel
type NewMessage struct {
	ChannelID ChannelID `json:"channelID"`
	Body      string    `json:"body"`
}

//MessageUpdates represents the field a message creator can modify
type MessageUpdates struct {
	Body string `json:"body"`
}
