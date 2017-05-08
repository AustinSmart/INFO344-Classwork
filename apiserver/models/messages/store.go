package messages

import "errors"

//ErrMessageNotFound is returned when the requested user is not found in the store
var ErrMessageNotFound = errors.New("message not found")

//ErrChannelNotFound is returned when the requested user is not found in the store
var ErrChannelNotFound = errors.New("channel not found")

//Store represents an abstract store for model.Channel and model.Message objects.
//This interface is used by the HTTP handlers to insert, get,
//and update new channels/messages. This interface can be implemented
//for any persistent database you want (e.g., MongoDB, PostgreSQL, etc.)
type Store interface {
	//GetAllChannels returns all channels a user is allowed to see
	GetAllChannels() ([]*Channel, error)

	//GetMessages returns `n` number of messages from a channel
	GetMessages(n int, channel ChannelID) (*[]Message, error)

	InsertChannel(newChannel *NewChannel) (*Channel, error)

	InsertMessage(newMessage *NewMessage) (*Message, error)

	AddUser(user *User, channel *Channel) (*Channel, err)

	//Insert inserts a new NewUser into the store
	//and returns a User with a newly-assigned ID
	Insert(newUser *NewUser) (*User, error)

	//Update applies UserUpdates to the currentUser
	Update(updates *UserUpdates, currentuser *User) error
}
