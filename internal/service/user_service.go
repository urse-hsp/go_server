package service

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/dao"
	userdto "go-server/internal/dto/user"
	"go-server/internal/model"
	"go-server/pkg/bcrypt"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	Create(ctx context.Context, req userdto.CreateRequest) (*model.User, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, id uint, req userdto.UpdateRequest) (*model.User, error)
	GetDetail(ctx context.Context, id uint) (*model.User, error)
	GetList(ctx context.Context, q userdto.RequestQuery) ([]model.User, error)
	GetPageList(ctx context.Context, q userdto.RequestPageQuery) ([]model.User, int64, error)

	Login(ctx context.Context, username string, password string) (*model.User, string, error)
}

func NewUserService(
	service *Service,
	Repo dao.UserRepository,
) UserService {
	return &userService{
		Repo:    Repo,
		Service: service,
	}
}

type userService struct {
	*Service
	Repo dao.UserRepository
}

// ================= 登录 =================

func (s *userService) Login(ctx context.Context, username string, password string) (*model.User, string, error) {
	data, err := s.Repo.GetByKeyWhere(ctx, username)
	if err != nil {
		return nil, "", err
	}

	if !bcrypt.CheckPassword(password, data.Password) {
		return nil, "", fmt.Errorf("密码错误")
	}

	// duration := time.Duration(s.) * time.Hour
	token, err := s.jwt.GenToken(data.ID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return nil, "", err
	}

	return data, token, nil
}

// ================= 创建 =================

func (s *userService) Create(ctx context.Context, req userdto.CreateRequest) (*model.User, error) {
	// 判断是否已存在
	data, err := s.Repo.GetByKeyWhere(ctx, req.Username)
	if err == nil && data != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	hashedPwd, err := bcrypt.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	data_ := &model.User{
		Username:  req.Username,
		Password:  hashedPwd,
		CreatedAt: time.Now(),
	}

	return s.Repo.Create(ctx, data_)
}

// ================= 删除 =================

func (s *userService) Delete(ctx context.Context, id uint) error {
	// 直接调用 DAO（DAO 已经判断是否存在）
	return s.Repo.Delete(ctx, id)
}

// ================= 更新 =================

func (s *userService) Update(ctx context.Context, id uint, req userdto.UpdateRequest) (*model.User, error) {
	data, err := s.Repo.GetDetail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	AssignIfNotNil(&data.Username, req.Username)
	AssignIfNotNil(&data.Avatar, req.Avatar)
	// if req.Username != nil {
	// 	data.Username = *req.Username
	// }

	return s.Repo.Update(ctx, data, id)
}

// ================= 获取 =================

func (s *userService) GetDetail(ctx context.Context, id uint) (*model.User, error) {
	return s.Repo.GetDetail(ctx, id)
}

// ================= 全部列表 =================

func (s *userService) GetList(ctx context.Context, q userdto.RequestQuery) ([]model.User, error) {
	return s.Repo.GetList(ctx, q)
}

// ================= 分页列表 =================

func (s *userService) GetPageList(ctx context.Context, q userdto.RequestPageQuery) ([]model.User, int64, error) {
	return s.Repo.GetPageList(ctx, q)
}
