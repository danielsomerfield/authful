package user

type UserRecord struct {
	Username string
	HashedPassword []byte
	AuthMethods []string
}

type SaveUser func (userRecord UserRecord) error