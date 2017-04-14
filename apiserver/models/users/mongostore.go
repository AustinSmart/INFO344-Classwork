package users

import (
	"gopkg.in/mgo.v2"
)

// MongoStore represents a users.store backed by MongoDB
type MongoStore struct {
	Session mgo.Session
}

// GetAll returns all users
func GetAll() ([]*User, error) {
	return nil, nil
}

// GetByID returns the User with the given ID
func GetByID(id UserID) (*User, error) {
	return nil, nil

}

// GetByEmail returns the User with the given email
func GetByEmail(email string) (*User, error) {
	return nil, nil

}

// GetByUserName returns the User with the given user name
func GetByUserName(name string) (*User, error) {
	return nil, nil

}

// Insert inserts a new NewUser into the store
// and returns a User with a newly-assigned ID
func Insert(newUser *NewUser) (*User, error) {
	return nil, nil

}

// Update applies UserUpdates to the currentUser
func Update(updates *UserUpdates, currentuser *User) error {
	return nil
}
