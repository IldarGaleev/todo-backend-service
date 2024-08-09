package storageDTO

type User struct {
	Id           uint64
	Username     string
	PasswordHash []byte
}
