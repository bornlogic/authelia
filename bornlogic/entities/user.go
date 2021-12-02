package entities

type User struct {
	DB
	Username       string `gorm:"type:varchar(100);uniqueIndex"`
	Email          string `gorm:"type:varchar(100);uniqueIndex"`
	HashedPassword string `gorm:"size:100"`
	Name           string
}
