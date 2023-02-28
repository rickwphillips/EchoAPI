package user

import (
	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	err := os.Remove(dbPath)
	if err != nil {
		return
	}
}

func TestCRUD(t *testing.T) {
	t.Log("Create")
	u := &User{
		ID:         bson.NewObjectId(),
		Name:       "test",
		First:      "Test",
		Last:       "Example",
		Email:      "texample@example.com",
		Department: "Marketing",
		Status:     "1",
	}
	err := u.Save()
	if err != nil {
		t.Fatalf("Error saving user: %s", err)
	}
	t.Log("Read")
	u2, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retrieving record: %s", err)
	}
	if !reflect.DeepEqual(u2, u) {
		t.Error("Records do not match")
	}
	t.Log("Update")
	u.Department = "Sales"
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving record: %s", err)
	}
	u3, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retrieving record: %s", err)
	}
	if !reflect.DeepEqual(u3, u) {
		t.Error("Records do not match")
	}
	t.Log("Delete")
	err = Delete(u.ID)
	if err != nil {
		t.Fatalf("Error removing record: %s", err)
	}
	_, err = One(u.ID)
	if err == nil {
		t.Fatal("Record should not exit")
	}
	if err != storm.ErrNotFound {
		t.Fatalf("Error retrieving non-existant record: %s", err)
	}
	t.Log("Read All")
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving record: %s", err)
	}
	u2.ID = bson.NewObjectId()
	u2.Name = "test2"
	err = u2.Save()
	if err != nil {
		t.Fatalf("Error saving record: %s", err)
	}
	u3.ID = bson.NewObjectId()
	u3.Name = "test3"
	err = u3.Save()
	if err != nil {
		t.Fatalf("Error saving record: %s", err)
	}
	users, err := All()
	if err != nil {
		t.Fatalf("Error reading all records: %s", err)
	}
	if len(users) != 3 {
		t.Errorf("Different number of records retreived. Expected 3 Actual: %d", len(users))
	}
	t.Log("Check Duplicate Username")
	u4, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error retrieving record: %s", err)
	}
	u4.ID = bson.NewObjectId()
	err = u4.Save()
	if err == nil {
		t.Fatalf("Duplicate username should not save. Username: %s", u4.Name)
	}
}
