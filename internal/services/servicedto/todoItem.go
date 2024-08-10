// Package servicedto contains serviceDTO models
package servicedto

//ToDoItem service DTO
type ToDoItem struct {
	ID         uint64
	Title      *string
	IsComplete *bool
	OwnerID    uint64
}
