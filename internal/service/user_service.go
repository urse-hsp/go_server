package service

import (
	"fmt"
	"go-demo-server/internal/dao"
	"go-demo-server/internal/model"
	"go-demo-server/pkg/bcrypt"
	"go-demo-server/pkg/jwt"
)

// ================= 登录 =================

func Login(username string, password string) (*model.User, string, error) {
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if !bcrypt.CheckPassword(password, user.Password) {
		return nil, "", fmt.Errorf("密码错误")
	}

	token := jwt.GenerateToken(user.ID, user.Username)

	return user, token, nil
}

// ================= 注册 =================

func Create(username string, password string) (*model.User, error) {
	hashedPwd, err := bcrypt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return dao.CreateUser(username, hashedPwd)
}

// ================= 获取用户 =================

func GetUserDetail(id uint) (*model.User, error) {
	return dao.GetUserByID(id)
}

// ================= 更新用户 =================

func UpdateUser(info model.User, id uint) (*model.User, error) {
	return dao.UpdateUser(info, id)
}

// ================= 删除用户 =================

func DeleteUser(id uint) error {
	// 直接调用 DAO（DAO 已经判断是否存在）
	return dao.DeleteUser(id)
}

// ================= 用户列表 =================

func GetUserList() ([]model.User, error) {
	return dao.GetUserList()
}

// ================= 用户列表 分页 =================

func GetUserLists(page, pageSize int) ([]model.User, int64, error) {
	return dao.GetUserLists(page, pageSize)
}
