package postgresStorageORM

type ToDoItemPG struct {
	Id         uint64 `gorm:"primaryKey;autoincrement;index:idx_user"`
	OwnerId    uint64 `gorm:"index:idx_owner"`
	Owner      UserPG `gorm:"constraint:OnDelete:CASCADE"`
	Title      string `gorm:"size:255;not null"`
	IsComplete bool   `gorm:"default:false"`
}

func (ToDoItemPG) TableName() string {
	return "todoItems"
}
