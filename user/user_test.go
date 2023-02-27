package user_test

import (
	"EchoAPI/user"
	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	It("can get a list of current users", func() {
		u := &user.User{
			bson.NewObjectId(), "admin", "Admin", "Example", "IT", "admin@example.com", "1",
		}

		err := u.ValidateUsername()
		Expect(err).To(HaveOccurred())
	})
})
