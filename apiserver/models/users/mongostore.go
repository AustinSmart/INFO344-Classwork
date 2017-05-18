package users

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoStore represents a users.store backed by MongoDB
type MongoStore struct {
	Session        *mgo.Session
	DatabaseName   string
	CollectionName string
}

// NewMongoStore returns a new MongoStore
func NewMongoStore(session *mgo.Session, dbName string, collectionName string) *MongoStore {
	return &MongoStore{
		Session:        session,
		DatabaseName:   dbName,
		CollectionName: collectionName,
	}
}

// GetAll returns all users
func (ms *MongoStore) GetAll() ([]*User, error) {
	users := []*User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(nil).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID returns the User with the given ID
func (ms *MongoStore) GetByID(id UserID) (*User, error) {
	u := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).FindId(id).One(u)
	if err == mgo.ErrNotFound {
		return nil, ErrUserNotFound
	}
	return u, err
}

// GetByEmail returns the User with the given email
func (ms *MongoStore) GetByEmail(email string) (*User, error) {
	u := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"email": email}).One(u)
	if err == mgo.ErrNotFound {
		return nil, ErrUserNotFound
	}
	return u, err
}

// GetByUserName returns the User with the given user name
func (ms *MongoStore) GetByUserName(name string) (*User, error) {
	u := &User{}
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Find(bson.M{"username": name}).One(u)
	if err == mgo.ErrNotFound {
		return nil, ErrUserNotFound
	}
	return u, err
}

// Insert inserts a new NewUser into the store
// and returns a User with a newly-assigned ID
func (ms *MongoStore) Insert(newUser *NewUser) (*User, error) {
	u, err := newUser.ToUser()
	u.ID = UserID(bson.NewObjectId().Hex())
	err = ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Insert(u)
	return u, err
}

// Update applies UserUpdates to the currentUser
func (ms *MongoStore) Update(updates *UserUpdates, currentuser *User) error {
	err := ms.Session.DB(ms.DatabaseName).C(ms.CollectionName).Update(bson.M{"_id": currentuser.ID}, bson.M{"$set": bson.M{"firstname": updates.FirstName, "lastname": updates.LastName}})
	return err
}
