package messages

import "errors"
import "github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"

//ErrMessageNotFound is returned when the requested message is not found in the store
var ErrMessageNotFound = errors.New("message not found")

//ErrChannelNotFound is returned when the requested channel is not found in the store
var ErrChannelNotFound = errors.New("channel not found")

//Store represents an abstract store for model.Channel and model.Message objects.
//This interface is used by the HTTP handlers to insert, get,
//and update new channels/messages. This interface can be implemented
//for any persistent database you want (e.g., MongoDB, PostgreSQL, etc.)
type Store interface {
	//GetAllChannels returns all channels a user is allowed to see
	GetAllChannels(user users.UserID) ([]*Channel, error)

	//GetMessages returns `n` number of messages from a channel
	GetMessages(n int, channel ChannelID) ([]*Message, error)

	//InsertChannel creates a new channel
	InsertChannel(user users.UserID, newChannel *NewChannel) (*Channel, error)

	//InsertMessage creates a new message
	InsertMessage(user users.UserID, newMessage *NewMessage) (*Message, error)

	//AddUser adds a user to a channels members list
	AddUser(user *users.UserID, channel ChannelID) error

	//RemoveUser removes a user from a channels member list
	RemoveUser(user *users.UserID, channel ChannelID) error

	//UpdateChannel updates a channels name and description
	UpdateChannel(updates *ChannelUpdates, channel ChannelID) error

	//UpdateMessage updates a messages body
	UpdateMessage(updates *MessageUpdates, message MessageID) error

	//DeleteChannel removes a channel and all messages within it
	DeleteChannel(channel ChannelID) error

	//DeleteMessage removes a message
	DeleteMessage(message MessageID) error
}
