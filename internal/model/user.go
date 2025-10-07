package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Email     string `gorm:"unique"`
	Roles     []Role `gorm:"many2many:user_roles"`
	Tasks     []Task
	CreatedAt time.Time
}

type Role struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"unique;not null"`
	Users []User `gorm:"many2many:user_roles"`
}
