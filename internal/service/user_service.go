package service

import (
	"context"
	"fmt"
	"go-server/internal/dao"
	"go-server/internal/model"
	"go-server/pkg/bcrypt"
	"time"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*model.User, string, error)
	Create(ctx context.Context, username string, password string) (*model.User, error)
	GetUserDetail(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, info model.User, id uint) (*model.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetUserList(ctx context.Context) ([]model.User, error)
	GetUserLists(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
}

type userService struct {
	*Service
	userRepo dao.UserRepository
}

func NewUserService(
	service *Service,
	userRepo dao.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

// ================= 登录 =================

func (s *userService) Login(ctx context.Context, username string, password string) (*model.User, string, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if !bcrypt.CheckPassword(password, user.Password) {
		return nil, "", fmt.Errorf("密码错误")
	}

	// token, err := jwt.GenerateToken(user.ID, user.Username)
	// duration := time.Duration(s.) * time.Hour
	token, err := s.jwt.GenToken(user.ID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ================= 注册 =================

func (s *userService) Create(ctx context.Context, username string, password string) (*model.User, error) {
	hashedPwd, err := bcrypt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, username, hashedPwd)
}

// ================= 获取用户 =================

func (s *userService) GetUserDetail(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// ================= 更新用户 =================

func (s *userService) UpdateUser(ctx context.Context, info model.User, id uint) (*model.User, error) {
	return s.userRepo.UpdateUser(ctx, info, id)
}

// ================= 删除用户 =================

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	// 直接调用 DAO（DAO 已经判断是否存在）
	return s.userRepo.DeleteUser(ctx, id)
}

// ================= 用户列表 =================

func (s *userService) GetUserList(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetUserList(ctx)
}

// ================= 用户列表 分页 =================

func (s *userService) GetUserLists(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.GetUserLists(ctx, page, pageSize)
}
