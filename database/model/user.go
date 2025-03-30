package model

type User struct {
	Id       int    `gorm:"primaryKey"`
	Name     string `gorm:"column:name"`
	PassWord string `gorm:"column:password"`
}
