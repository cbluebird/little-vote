package model

type User struct {
	UserId int `gorm:"primary_key;AUTO_INCREMENT"`
	Count  int
	Name   string
}
