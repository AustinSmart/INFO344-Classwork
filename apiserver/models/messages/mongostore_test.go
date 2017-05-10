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

	c, err = s.GetChannel(c.ID)
	if err != nil {
		t.Errorf("error getting channel: %v\n", err)
	}
	if nil == c {
		t.Fatalf("nil returned from Store.GetChannel()")
	}

	allC, err := s.GetAllChannels(uID[0])
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
		t.Fatalf("getMessages failed. expected: %d. recieved: %d\n", 2, len(allM))
	}

	tm, err := s.GetMessage(allM[0].ID)
	if err != nil {
		t.Errorf("error getting message: %v\n", err)
	}
	if tm.ID != allM[0].ID {
		t.Fatalf("getMessage failed. expected: %s. recieved: %s\n", allM[0].ID, tm.ID)
	}

	uID = append(uID, "4444")
	err = s.AddUser(&uID[3], c.ID)
	if err != nil {
		t.Errorf("error adding user: %v\n", err)
	}
	allC, err = s.GetAllChannels(uID[0])
	if allC[0].Members[3] == "" || allC[0].Members[3] != uID[3] {
		t.Fatalf("AddUser failed. expected: %s. recieved: %s\n", uID[3], allC[0].Members[3])
	}

	err = s.RemoveUser(&uID[3], c.ID)
	if err != nil {
		t.Errorf("error removing user: %v\n", err)
	}
	allC, err = s.GetAllChannels(uID[0])
	if len(allC[0].Members) > 3 {
		if allC[0].Members[3] != "" || allC[0].Members[3] == uID[3] {
			t.Fatalf("AddUser failed. expected: %s. recieved: %s\n", "nothing", allC[0].Members[3])
		}
	}

	cup := ChannelUpdates{
		Name:        "Test Channel UPDATED",
		Description: "This is only a test UPDATED",
	}

	err = s.UpdateChannel(&cup, c.ID)
	if err != nil {
		t.Errorf("error udating channel: %v\n", err)
	}
	allC, err = s.GetAllChannels(uID[0])
	if allC[0].Name != cup.Name {
		t.Fatalf("UpdateChannel name failed. expected: %s. recieved: %s\n", cup.Name, allC[0].Name)
	}
	if allC[0].Description != cup.Description {
		t.Fatalf("UpdateChannel description failed. expected: %s. recieved: %s\n", cup.Description, allC[0].Description)
	}

	mup := MessageUpdates{
		Body: "This is a test message UPDATED",
	}

	err = s.UpdateMessage(&mup, allM[0].ID)
	if err != nil {
		t.Errorf("error udating message: %v\n", err)
	}
	allM, err = s.GetMessages(500, c.ID)
	if allM[0].Body != mup.Body {
		t.Fatalf("Update message failed. expected: %s. recieved: %s\n", mup.Body, allM[0].Body)
	}

	err = s.DeleteMessage(allM[1].ID)
	if err != nil {
		t.Errorf("error deleting message: %v\n", err)
	}
	allM, err = s.GetMessages(500, c.ID)
	if len(allM) > 1 {
		t.Fatalf("Update message failed. expected: %d. recieved: %d\n", 1, len(allM))
	}

	err = s.DeleteChannel(c.ID)
	if err != nil {
		t.Errorf("error deleting channel: %v\n", err)
	}
	allC, err = s.GetAllChannels(uID[0])
	if err != nil {
		t.Errorf("error getting all channels: %v\n", err)
	}
	if len(allC) != 0 {
		t.Fatalf("expected: %d channels, recieved: %d\n", 0, len(allC))
	}

	sess.DB(s.DatabaseName).C(s.ChannelsCollectionName).RemoveAll(nil)
	sess.DB(s.DatabaseName).C(s.MessagesCollectionName).RemoveAll(nil)
}
