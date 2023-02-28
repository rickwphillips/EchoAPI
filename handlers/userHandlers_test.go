package handlers

import (
	"EchoAPI/user"
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestBodyToUse(t *testing.T) {

	valid := &user.User{
		ID:         bson.NewObjectId(),
		Name:       "jappleseed",
		First:      "Jonny",
		Last:       "Appleseed",
		Email:      "jappleseed@example.com",
		Department: "Marketing",
		Status:     "1",
	}

	valid2 := &user.User{
		ID:         valid.ID,
		Name:       "jappleseed",
		First:      "Jonny",
		Last:       "Appleseed",
		Email:      "jappleseed@example.com",
		Department: "Sales",
		Status:     "0",
	}

	js, err := json.Marshal(valid)
	if err != nil {
		t.Errorf("Error unmarshalling a valid user: %s", err)
		t.FailNow()
	}

	ts := []struct {
		txt string
		r   *http.Request
		u   *user.User
		err bool
		exp *user.User
	}{
		{
			txt: "nil request",
			err: true,
		},
		{
			txt: "nil request body",
			r:   &http.Request{},
			err: true,
		},
		{
			txt: "malformed data",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"id"=12}`)),
			},
			u:   &user.User{},
			err: true,
		},
		{
			txt: "valid request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBuffer(js)),
			},
			u:   &user.User{},
			exp: valid,
		},
		{
			txt: "valid partial request",
			r: &http.Request{
				Body: ioutil.NopCloser(
					bytes.NewBufferString(
						`{"user_status":"0","department":"Sales"}`)),
			},
			u:   valid,
			exp: valid2,
		},
	}

	// Test handlers
	for _, tc := range ts {
		t.Log(tc.txt)
		err := bodyToUser(tc.r, tc.u)
		if tc.err {
			if err == nil {
				t.Error("Expected error, got none")
			}
			continue
		}
		if err != nil {
			t.Errorf("Unexpected error: %s", tc.exp)
			continue
		}
		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different")
			t.Error(tc.u)
			t.Error(tc.exp)
		}
	}
}
