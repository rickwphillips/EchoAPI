package user

import (
	"errors"
	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID         bson.ObjectId `json:"user_id" storm:"id"`
	Name       string        `json:"user_name"`
	First      string        `json:"first_name"`
	Last       string        `json:"last_name"`
	Department string        `json:"department"`
	Email      string        `json:"email"`
	Status     string        `json:"user_status"`
}

const (
	dbPath = "users.db"
)

var (
	ErrRecordInvalid  = errors.New("record is invalid")
	ErrUniqueUsername = errors.New("username must be unique")
)

// All Get all users from the database
func All() ([]User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var users []User
	err = db.All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// One Returns a single user record or nil if not found
func One(id bson.ObjectId) (*User, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	u := new(User)
	err = db.One("ID", id, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Delete Remove a single user record from the database
func Delete(id bson.ObjectId) error {
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	u := new(User)
	err = db.One("ID", id, u)
	if err != nil {
		return err
	}

	return db.DeleteStruct(u)
}

// Save Create or update the provided user
func (u *User) Save() error {
	if err := u.validate(); err != nil {
		return err
	}

	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Save(u)
}

// validate User data validation for User
func (u *User) validate() error {
	if u.Name == "" || u.First == "" || u.Last == "" || u.Email == "" || u.Department == "" || u.Status == "" {
		return ErrRecordInvalid
	}
	return nil
}

// ValidateUsername Ensure user_name is unique
func (u *User) ValidateUsername() error {
	db, err := storm.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.One("Name", u.Name, u)

	if err == storm.ErrNotFound {
		return nil
	}

	return ErrUniqueUsername
}
