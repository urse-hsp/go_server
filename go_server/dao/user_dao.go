package dao

import (
	"errors"
	"fmt"
	"go-demo-server/bootstrap"
	"go-demo-server/model"
	"time"

	"gorm.io/gorm"
)

func GetUserByID(id uint) (*model.User, error) {
	var user model.User

	result := bootstrap.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetSeatchInfo(username string) (*model.User, error) {
	var user model.User

	result := bootstrap.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s不存在", username) // 或者返回特定的错误
		}
		// 如果是其他数据库错误（如连接断开），也返回
		return nil, result.Error
	}

	return &user, nil
}

func CreateUser(username string, password string) (*model.User, error) {
	_, err := GetSeatchInfo(username)
	if err == nil {
		return nil, fmt.Errorf("用户名已经存在")
	}

	newUser := &model.User{
		Username:  username, // 注意：这里字段名首字母必须大写，且要和数据库对应
		Password:  password, // 实际开发中这里通常需要加密
		CreatedAt: time.Now(),
	}

	result := bootstrap.DB.Create(&newUser)
	if result.Error != nil {
		// 这里通常不需要判断 ErrRecordNotFound，因为 Create 失败通常是因为主键冲突或字段长度不够
		// 直接返回数据库错误即可
		return nil, result.Error
	}

	return newUser, nil
}

func UpdateUser(user model.User, id uint) (*model.User, error) {
	result := bootstrap.DB.Where("id = ?", id).Updates(&user) // ✅ 更清晰
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func DeleteUser(id uint) error {
	result := bootstrap.DB.Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
