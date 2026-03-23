package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"size:50;not null;unique"`
	Password  string `gorm:"size:255;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Avatar    string `gorm: "size:255"; not null; default: 'https://cdn-icons-png.flaticon.com/512/149/149071.png'`
	Age       int    `gorm:"not null; default: 0"`
}
