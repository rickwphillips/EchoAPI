package handlers

import (
	"EchoAPI/cache"
	"EchoAPI/user"
	"encoding/json"
	"errors"
	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

// usersGetAll Get list of all users
func usersGetAll(w http.ResponseWriter, r *http.Request) {
	if cache.Serve(w, r) {
		return
	}
	users, err := user.All()
	if err != nil {
		postError(w, http.StatusInternalServerError)
		return
	}
	cw := cache.NewWriter(w, r)
	postBodyResponse(cw, http.StatusOK, jsonResponse{"users": users})
}

func bodyToUser(r *http.Request, u *user.User) error {
	if r == nil {
		return errors.New("a request is required")
	}
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

// usersPostOne Create a single user
func usersPostOne(w http.ResponseWriter, r *http.Request) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	err = u.ValidateUsername()
	if err != nil {
		if err == user.ErrUniqueUsername {
			postBodyResponse(w, http.StatusBadRequest, jsonResponse{"error": "username must be unique"})
			return
		}
		postError(w, http.StatusInternalServerError)
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

	cache.Drop("/users")
	w.Header().Set("Location", "/users/"+u.ID.Hex())
	w.WriteHeader(http.StatusCreated)

}

// usersGetOne Get a single user
func usersGetOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	if cache.Serve(w, r) {
		return
	}
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cw := cache.NewWriter(w, r)
	postBodyResponse(cw, http.StatusOK, jsonResponse{"user": u})
}

// usersPutOne Replace a single user
func usersPutOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u := new(user.User)
	err := bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}

	err = u.ValidateUsername()
	if err != nil {
		if err == user.ErrUniqueUsername {
			postBodyResponse(w, http.StatusBadRequest, jsonResponse{"error": "username must be unique"})
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}

	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cw := cache.NewWriter(w, r)
	postBodyResponse(cw, http.StatusOK, jsonResponse{"user": u})

}

// usersPatchOne Update a single user
func usersPatchOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	err = bodyToUser(r, u)
	if err != nil {
		postError(w, http.StatusBadRequest)
		return
	}
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			postError(w, http.StatusBadRequest)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cw := cache.NewWriter(w, r)
	postBodyResponse(cw, http.StatusOK, jsonResponse{"user": u})

}

func usersDeleteOne(w http.ResponseWriter, r *http.Request, id bson.ObjectId) {
	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			postError(w, http.StatusNotFound)
			return
		}
		postError(w, http.StatusInternalServerError)
		return
	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(r))
	w.WriteHeader(http.StatusOK)
}
