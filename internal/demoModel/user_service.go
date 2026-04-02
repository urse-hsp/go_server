package demo

import (
	"context"
	"fmt"
	"go-server/internal/dao"
	"go-server/internal/model"
	"go-server/internal/service"
	"go-server/pkg/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*model.User, string, error)

	Create(ctx context.Context, username string, password string) (*model.User, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, info model.User, id uint) (*model.User, error)
	GetDetail(ctx context.Context, id uint) (*model.User, error)
	GetList(ctx context.Context) ([]model.User, error)
	GetLists(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
}

func NewUserService(
	service *service.Service,
	userRepo dao.UserRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

type userService struct {
	*service.Service
	userRepo dao.UserRepository
}

// ================= 登录 =================

func (s *userService) Login(ctx context.Context, username string, password string) (*model.User, string, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if !bcrypt.CheckPassword(password, user.Password) {
		return nil, "", fmt.Errorf("密码错误")
	}

	// token, err := s.jwt.GenToken(user.ID, time.Now().Add(time.Hour*24*90))
	// if err != nil {
	// 	return nil, "", err
	// }

	// return user, token, nil

	return user, "", nil
}

// ================= 注册 =================

func (s *userService) Create(ctx context.Context, username string, password string) (*model.User, error) {
	hashedPwd, err := bcrypt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Create(ctx, username, hashedPwd)
}

// ================= 删除 =================

func (s *userService) Delete(ctx context.Context, id uint) error {
	// 直接调用 DAO（DAO 已经判断是否存在）
	return s.userRepo.Delete(ctx, id)
}

// ================= 更新 =================

func (s *userService) Update(ctx context.Context, info model.User, id uint) (*model.User, error) {
	return s.userRepo.Update(ctx, info, id)
}

// ================= 获取 =================

func (s *userService) GetDetail(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetDetail(ctx, id)
}

// ================= 全部列表 =================

func (s *userService) GetList(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetList(ctx)
}

// ================= 分页列表 =================

func (s *userService) GetLists(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.GetLists(ctx, page, pageSize)
}
