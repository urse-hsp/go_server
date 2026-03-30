package service

import (
	"fmt"
	"go-server/internal/dao"
	"go-server/internal/model"
	"go-server/pkg/bcrypt"
	"go-server/pkg/jwt"
)

type UserService interface {
	Login(username string, password string) (*model.User, string, error)
	Create(username string, password string) (*model.User, error)
	GetUserDetail(id uint) (*model.User, error)
	UpdateUser(info model.User, id uint) (*model.User, error)
	DeleteUser(id uint) error
	GetUserList() ([]model.User, error)
	GetUserLists(page, pageSize int) ([]model.User, int64, error)
}

type userService struct {
	userRepo dao.UserRepository
}

func NewUserService(repo dao.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}

// ================= 登录 =================

func (s *userService) Login(username string, password string) (*model.User, string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, "", err
	}

	if !bcrypt.CheckPassword(password, user.Password) {
		return nil, "", fmt.Errorf("密码错误")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ================= 注册 =================

func (s *userService) Create(username string, password string) (*model.User, error) {
	hashedPwd, err := bcrypt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(username, hashedPwd)
}

// ================= 获取用户 =================

func (s *userService) GetUserDetail(id uint) (*model.User, error) {
	return s.userRepo.GetUserByID(id)
}

// ================= 更新用户 =================

func (s *userService) UpdateUser(info model.User, id uint) (*model.User, error) {
	return s.userRepo.UpdateUser(info, id)
}

// ================= 删除用户 =================

func (s *userService) DeleteUser(id uint) error {
	// 直接调用 DAO（DAO 已经判断是否存在）
	return s.userRepo.DeleteUser(id)
}

// ================= 用户列表 =================

func (s *userService) GetUserList() ([]model.User, error) {
	return s.userRepo.GetUserList()
}

// ================= 用户列表 分页 =================

func (s *userService) GetUserLists(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.GetUserLists(page, pageSize)
}
