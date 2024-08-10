// Package postgresstorageorm contains Postgres ORM models
package postgresstorageorm

type ToDoItemPG struct {
	ID         uint64 `gorm:"primaryKey;autoincrement;index:idx_todo_item"`
	OwnerID    uint64 `gorm:"index:idx_owner"`
	Owner      UserPG `gorm:"constraint:OnDelete:CASCADE"`
	Title      string `gorm:"size:255;not null"`
	IsComplete bool   `gorm:"default:false"`
}

func (ToDoItemPG) TableName() string {
	return "todoItems"
}
