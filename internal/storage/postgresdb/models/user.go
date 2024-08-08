package postgresStorageORM

type UserPG struct {
	Id           uint64 `gorm:"primaryKey;autoincrement;index:idx_user"`
	Username     string `gorm:"size:40;not null;unique"`
	PasswordHash []byte `gorm:"not null"`
}

func (UserPG) TableName() string {
	return "users"
}
