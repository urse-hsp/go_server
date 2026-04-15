package model

import "gorm.io/gorm"

type Demo struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"` // 模型里有 gorm.DeletedAt，Delete 自动就是软删除
}

func (m *Demo) TableName() string {
	return "demo"
}
