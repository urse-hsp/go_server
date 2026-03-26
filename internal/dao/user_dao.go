package dao

import (
	"errors"
	"fmt"
	"go-demo-server/internal/bootstrap"
	"go-demo-server/internal/model"

	"gorm.io/gorm"
)

// ================= 根据ID查询 =================

func GetUserByID(id uint) (*model.User, error) {
	var user model.User

	err := bootstrap.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= 根据用户名查询 =================

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	err := bootstrap.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= 创建用户 =================

func CreateUser(username string, password string) (*model.User, error) {
	// 判断是否已存在
	_, err := GetUserByUsername(username)

	if err == nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 如果不是“未找到”，说明是数据库错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: password,
		// CreatedAt: time.Now(),
	}

	if err := bootstrap.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ================= 更新用户 =================

func UpdateUser(user model.User, id uint) (*model.User, error) {
	result := bootstrap.DB.Model(&model.User{}).
		Where("id = ?", id).
		Updates(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	// 重新查询最新数据（关键）
	var updatedUser model.User
	if err := bootstrap.DB.First(&updatedUser, id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// ================= 删除用户 =================

func DeleteUser(id uint) error {
	result := bootstrap.DB.Where("id = ?", id).Delete(&model.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// ================= 用户列表 =================

func GetUserList() ([]model.User, error) {
	var users []model.User

	if err := bootstrap.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// ================= 用户列表 分页 =================

func GetUserLists(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User

	db := bootstrap.DB.Model(&model.User{})

	// 分页
	total, err := Paginate(db, &users, page, pageSize)

	return users, total, err
}
