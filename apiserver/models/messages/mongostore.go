package messages

import (
	"time"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//MongoStore represents a messages.store backed by MongoDB
type MongoStore struct {
	Session                *mgo.Session
	DatabaseName           string
	MessagesCollectionName string
	ChannelsCollectionName string
}

//NewMongoStore returns a new MongoStore
func NewMongoStore(session *mgo.Session, dbName string, MessagesCollectionName string, ChannelsCollectionName string) *MongoStore {
	return &MongoStore{
		Session:                session,
		DatabaseName:           dbName,
		MessagesCollectionName: MessagesCollectionName,
		ChannelsCollectionName: ChannelsCollectionName,
	}
}

//GetAllChannels returns all channels a user is allowed to see
func (ms *MongoStore) GetAllChannels(user users.UserID) ([]*Channel, error) {
	channels := []*Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Find(bson.M{"$or": []bson.M{bson.M{"members": user}, bson.M{"private": false}}}).All(&channels)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

//GetChannel returns the requested channel
func (ms *MongoStore) GetChannel(channel ChannelID) (*Channel, error) {
	c := Channel{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Find(bson.M{"_id": channel}).One(&c)
	if err == mgo.ErrNotFound {
		return nil, ErrChannelNotFound
	}
	return &c, err
}

//GetMessages returns `n` number of messages from a channel
func (ms *MongoStore) GetMessages(n int, channel ChannelID) ([]*Message, error) {
	messages := []*Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessagesCollectionName).Find(bson.M{"channelid": channel}).Limit(n).All(&messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

//GetMessage returns the requested message
func (ms *MongoStore) GetMessage(message MessageID) (*Message, error) {
	m := Message{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessagesCollectionName).Find(bson.M{"_id": message}).One(&m)
	if err == mgo.ErrNotFound {
		return nil, ErrMessageNotFound
	}
	return &m, err
}

//InsertChannel creates a new channel
func (ms *MongoStore) InsertChannel(user users.UserID, newChannel *NewChannel) (*Channel, error) {
	channel := newChannel.ToChannel()
	channel.ID = ChannelID(bson.NewObjectId().Hex())
	channel.CreatorID = user
	channel.Members = append(channel.Members, user)
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Insert(channel)
	return channel, err
}

//InsertMessage creates a new message
func (ms *MongoStore) InsertMessage(user users.User, newMessage *NewMessage) (*Message, error) {
	message := newMessage.ToMessage()
	message.ID = MessageID(bson.NewObjectId().Hex())
	message.CreatorID = user.ID
	message.CreatorName = user.FirstName + " " + user.LastName
	message.CreatorPhotoURL = user.PhotoURL
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessagesCollectionName).Insert(message)
	return message, err
}

//AddUser adds a user to a channels members list
func (ms *MongoStore) AddUser(user *users.UserID, channel ChannelID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Update(bson.M{"_id": channel}, bson.M{"$push": bson.M{"members": user}})
	return err
}

//RemoveUser removes a user from a channels member list
func (ms *MongoStore) RemoveUser(user *users.UserID, channel ChannelID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Update(bson.M{"_id": channel}, bson.M{"$pull": bson.M{"members": user}})
	if err == mgo.ErrNotFound {
		return users.ErrUserNotFound
	}
	return err
}

//UpdateChannel updates a channels name and description
func (ms *MongoStore) UpdateChannel(updates *ChannelUpdates, channel ChannelID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Update(bson.M{"_id": channel}, bson.M{"$set": bson.M{"name": updates.Name, "description": updates.Description, "editedat": time.Now().String()}})
	if err == mgo.ErrNotFound {
		return ErrChannelNotFound
	}
	return err
}

//UpdateMessage updates a messages body
func (ms *MongoStore) UpdateMessage(updates *MessageUpdates, message MessageID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessagesCollectionName).Update(bson.M{"_id": message}, bson.M{"$set": bson.M{"body": updates.Body, "editedat": time.Now().String()}})
	if err == mgo.ErrNotFound {
		return ErrMessageNotFound
	}
	return err
}

//DeleteChannel removes a channel and all messages within it
func (ms *MongoStore) DeleteChannel(channel ChannelID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.ChannelsCollectionName).Remove(bson.M{"_id": channel})
	if err == mgo.ErrNotFound {
		return ErrChannelNotFound
	}
	return err
}

//DeleteMessage removes a message
func (ms *MongoStore) DeleteMessage(message MessageID) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.MessagesCollectionName).Remove(bson.M{"_id": message})
	if err == mgo.ErrNotFound {
		return ErrMessageNotFound
	}
	return err
}
