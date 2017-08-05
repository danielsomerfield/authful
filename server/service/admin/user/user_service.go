package user

import (
	userRepo "github.com/danielsomerfield/authful/server/repository/user"
	"github.com/danielsomerfield/authful/server/service/crypto"
)

type RegisterUserFn func(user User) error

type User struct {
	Username    string
	Password    string
	AuthMethods []string
}

func NewRegisterUserFn(saveUserFn userRepo.SaveUser, hashFn crypto.BuildPwdHash) RegisterUserFn {
	return func(user User) error {
		return saveUserFn(userRepo.UserRecord{
			Username: user.Username,
			HashedPassword: hashFn(user.Password),
			AuthMethods: user.AuthMethods,
		})
	}
}
