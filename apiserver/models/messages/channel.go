package messages

import "github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"

//ChannelID defines the type for channel ID's
type ChannelID string

//Channel represents a channel in the database
type Channel struct {
	ID          ChannelID      `json:"id" bson:"_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   string         `json:"time"`
	CreatorID   users.UserID   `json:"creatorID"`
	Members     []users.UserID `json:"members"`
	Private     bool           `json:"private"`
}

//NewChannel represents a new channel to be created
type NewChannel struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Members     []users.UserID `json:"members"`
	Private     bool           `json:"private"`
}

// ChannelUpdates represents the fields a channel creator can modify
type ChannelUpdates struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
