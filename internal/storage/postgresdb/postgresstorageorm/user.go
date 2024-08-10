package postgresstorageorm

type UserPG struct {
	ID           uint64 `gorm:"primaryKey;autoincrement;index:idx_user"`
	Username     string `gorm:"size:40;not null;unique"`
	PasswordHash []byte `gorm:"not null"`
}

func (UserPG) TableName() string {
	return "users"
}
