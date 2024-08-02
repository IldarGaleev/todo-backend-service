package serviceDTO

//ToDoItem service DTO
type ToDoItem struct {
	Id         uint64
	Title      *string
	IsComplete *bool
	OwnerId    uint64
}
