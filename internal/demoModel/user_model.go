package demo

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;`
	Username  string    `gorm:"size:50;not null;unique;comment:用户名1"`
	Password  string    `gorm:"size:255;not null;comment:密码"`
	CreatedAt time.Time `gorm:"comment:注册时间"`
	UpdatedAt time.Time `gorm:"comment:更新时间"`
	Avatar    string    `gorm:"size:255;not null;default:'https://cdn-icons-png.flaticon.com/512/149/149071.png';comment:头像"`
	Age       int       `gorm:"not null; default: 0;comment:年龄"`
	Email     string    `gorm:"type:varchar(100);not null;comment:电子邮件"`
	Phone     string    `gorm:"type:varchar(20);not null;comment:手机号"`
	Gender    string    `gorm:"type:varchar(10);not null;comment:性别"`
}

// gorm 默认表名是结构体名的复数形式，重写 TableName 方法可以自定义表名
func (User) TableName() string {
	return "users"
}
