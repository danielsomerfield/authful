package user

import (
	"testing"
	userrepo "github.com/danielsomerfield/authful/server/repository/user"
	"github.com/danielsomerfield/authful/util"
)

//Duplicate username fails
//Invalid auth method fails

var userRecords = []userrepo.UserRecord{}

func mockSaveFn(userRecord userrepo.UserRecord) error {
	userRecords = append(userRecords, userRecord)
	return nil
}

func mockHashFn(pwd string)[]byte {
	return []byte(pwd)
}

func setup() {
	userRecords = []userrepo.UserRecord{}
}

func TestRegisterUserFunction_RegistersValidUser(t *testing.T) {

	setup()
	user := User{
		Username: "username1",
		Password:    "password1",
		AuthMethods: []string{"username-password"},
	}

	registerFn := NewRegisterUserFn(mockSaveFn, mockHashFn)
	registerFn(user)

	expected := userrepo.UserRecord{
		Username: "username1",
		HashedPassword: []byte("password1"),
		AuthMethods: []string{"username-password"},
	}

	util.AssertTrue(len(userRecords) == 1, "A use record was entered", t)
	util.AssertEquals(expected, userRecords[0], t)

}
