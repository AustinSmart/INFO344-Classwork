package messages

import (
	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStore represents a messages.store backed by MongoDB
type MongoStore struct {
	Session                *mgo.Session
	DatabaseName           string
	MessagesCollectionName string
	ChannelsCollectionName string
}

// NewMongoStore returns a new MongoStore
func NewMongoStore(session *mgo.Session, dbName string, MessagesCollectionName string, ChannelsCollectionName string) *MongoStore {
	return &MongoStore{
		Session:                session,
		DatabaseName:           dbName,
		MessagesCollectionName: MessagesCollectionName,
		ChannelsCollectionName: ChannelsCollectionName,
	}
}

//GetAllChannels returns all channels a user is allowed to see
func (ms *MongoStore) GetAllChannels() ([]*Channel, error) {
	channels := []*Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Find(nil).All(&channels)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

//GetMessages returns `n` number of messages from a channel
func (ms *MongoStore) GetMessages(n int, channel ChannelID) ([]*Message, error) {
	messages := []*Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessagesCollectionName).Find(bson.M{"channelID": channel}).Limit(n).All(&messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

//InsertChannel creates a new channel
func (ms *MongoStore) InsertChannel(user users.UserID, newChannel *NewChannel) (*Channel, error) {
	channel := newChannel.ToChannel()
	channel.ID = ChannelID(bson.NewObjectId().Hex())
	channel.CreatorID = user
	//TODO insert to db
	return channel, nil
}

//InsertMessage creates a new message
func (ms *MongoStore) InsertMessage(user users.UserID, newMessage *NewMessage) (*Message, error) {
	message := newMessage.ToMessage()
	message.ID = MessageID(bson.NewObjectId().Hex())
	message.CreatorID = user
	//TODO insert to db
	return message, nil
}

//InsertUser adds a user to a channels members list
func (ms *MongoStore) InsertUser(user *users.UserID, channel ChannelID) (*Channel, error) {
	return nil, nil
}

//UpdateChannel updates a channels name and description
func (ms *MongoStore) UpdateChannel(updates *ChannelUpdates, channel ChannelID) (*Channel, error) {
	return nil, nil
}

//UpdateMessage updates a messages body
func (ms *MongoStore) UpdateMessage(updates *MessageUpdates, message MessageID) (*Message, error) {
	return nil, nil
}

//DeleteChannel removes a channel and all messages within it
func (ms *MongoStore) DeleteChannel(channel ChannelID) error {
	return nil
}

//DeleteMessage removes a message
func (ms *MongoStore) DeleteMessage(message MessageID) error {
	return nil
}

//RemoveUser removes a user from a channels member list
func (ms *MongoStore) RemoveUser(user *users.UserID, channel ChannelID) error {
	return nil
}
