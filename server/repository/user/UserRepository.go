package user

type UserRecord struct {
	Username       string
	HashedPassword string
	AuthMethods    []string
}

type SaveUser func(userRecord UserRecord) error

func NewInMemoryUserRepository() inMemoryUserRepository {
	return inMemoryUserRepository{}
}

type inMemoryUserRepository struct {
	userRecords map[string]UserRecord
}

func (repo *inMemoryUserRepository) SaveUser(userRecord UserRecord) error {
	repo.userRecords[userRecord.Username] = userRecord
	return nil
}
