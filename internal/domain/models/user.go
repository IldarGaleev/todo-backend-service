package models

type User struct {
	userId       uint64
	email        string
	passwordHash []byte
}
