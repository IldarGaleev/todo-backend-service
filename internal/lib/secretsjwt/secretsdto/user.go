// Package secretsdto contains DTO for secrets
package secretsdto

type User struct {
	UserID   *uint64
	Username *string
	Payload  interface{}
}