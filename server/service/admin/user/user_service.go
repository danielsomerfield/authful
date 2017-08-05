package user

import userRepo "github.com/danielsomerfield/authful/server/repository/user"

type RegisterUserFn func(user User) error

type User struct {
	Username    string
	Password    string
	AuthMethods []string
}

type BuildHash func(string) []byte

func NewRegisterUserFn(saveUserFn userRepo.SaveUser, hashFn BuildHash) RegisterUserFn {
	return func(user User) error {
		return saveUserFn(userRepo.UserRecord{
			Username: user.Username,
			HashedPassword: hashFn(user.Password),
			AuthMethods: user.AuthMethods,
		})
	}
}
