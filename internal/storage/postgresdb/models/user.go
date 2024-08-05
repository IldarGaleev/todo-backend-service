package postgresStorageORM

type UserPG struct {
	Id           uint64 `gorm:"primaryKey;autoincrement"`
	Username     string `gorm:"size:40;not null"`
	PasswordHash []byte `gorm:""`
}

func (UserPG) TableName() string {
	return "users"
}
