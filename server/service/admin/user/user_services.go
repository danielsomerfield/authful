package user

type RegisterUserFn func(user User) error

type User struct {
	Username    string
	Password    string
	AuthMethods []string
}
