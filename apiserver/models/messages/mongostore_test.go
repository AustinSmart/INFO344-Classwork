package messages

import (
	"testing"

	"github.com/info344-s17/challenges-AustinSmart/apiserver/models/users"

	mgo "gopkg.in/mgo.v2"
)

func TestCRUD(t *testing.T) {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		t.Fatalf("error dialing Mongo: %v", err)
	}
	defer sess.Close()

	s := NewMongoStore(sess, "devData", "devMessages", "devChannels")

	uID := []users.UserID{"11111", "22222", "33333"}

	nc := &NewChannel{
		Name:        "Test Channel",
		Description: "This is only a test",
		Members:     uID,
		Private:     false,
	}

	c, err := s.InsertChannel(uID[0], nc)
	if err != nil {
		t.Errorf("error inserting channel: %v\n", err)
	}
	if nil == c {
		t.Fatalf("nil returned from Store.InsertChannel()--you probably haven't implemented NewChannel.ToChannel() yet")
	}
	if len(string(c.ID)) == 0 {
		t.Errorf("new ID is zero-length\n")
	}
	if len(c.Members) != len(uID) {
		t.Errorf("inserted channels members is incorrect. expected: %d, recieved: %d\n", len(uID), len(c.Members))
	}

	allC, err := s.GetAllChannels()
	if err != nil {
		t.Errorf("error getting all channels: %v\n", err)
	}
	if len(allC) != 1 {
		t.Fatalf("expected: %d channels, recieved: %d\n", 1, len(allC))
	}

	nm := &NewMessage{
		ChannelID: c.ID,
		Body:      "This is a test message",
	}
	m, err := s.InsertMessage(uID[0], nm)
	m, err = s.InsertMessage(uID[0], nm)
	if err != nil {
		t.Errorf("error inserting message: %v\n", err)
	}
	if nil == m {
		t.Fatalf("nil returned from Store.InsertMessage()--you probably haven't implemented NewMessage.toMessage() yet")
	}

	limit := 1
	allM, err := s.GetMessages(limit, c.ID)
	if err != nil {
		t.Errorf("error getting messages: %v\n", err)
	}
	if len(allM) != limit {
		t.Fatalf("limit 1 on getMessage failed. recieved: %d\n", len(allM))
	}
	limit = 500
	allM, err = s.GetMessages(500, c.ID)
	if err != nil {
		t.Errorf("error getting messages: %v\n", err)
	}
	if len(allM) != 2 {
		t.Fatalf("getMessage failed. expected: %d. recieved: %d\n", 2, len(allM))
	}

	sess.DB(s.DatabaseName).C(s.ChannelsCollectionName).RemoveAll(nil)
	sess.DB(s.DatabaseName).C(s.MessagesCollectionName).RemoveAll(nil)
}
