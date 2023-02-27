package handlers

import (
	"EchoAPI/user"
	"encoding/json"
	"errors"
	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

// usersGetAll Get list of all users
func usersGetAll(w http.ResponseWriter, _ *http.Request) {
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	postBodyResponse(w, http.StatusNotFound, jsonResponse{"users": users})
}

func bodyToUser(r *http.Request, u *user.User) error {
	if r.Body == nil {
		return errors.New("request body in empty")
	}
	if u == nil {
		return errors.New("user is required")
	}
	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bd, u)
}

// usersPostOne Create or update a single user
func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)

}

// usersGetOne Get a single user
func usersGetOne(w http.ResponseWriter, _ *http.Request, id bson.ObjectId) {
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	postBodyResponse(w, http.StatusOK, jsonResponse{"user": u})
}
