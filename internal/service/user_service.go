package service

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/bootstrap"
	"go-server/internal/client"
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

	HttpUserList(ctx context.Context) (any, error)
}

func NewUserService(
	service *Service,
	Repo dao.UserRepository,
	userClient client.UserClient,
) UserService {
	return &userService{
		Repo:       Repo,
		Service:    service,
		userClient: userClient,
	}
}

type userService struct {
	*Service
	Repo       dao.UserRepository
	userClient client.UserClient
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

// http数据
func (s *userService) HttpUserList(ctx context.Context) (any, error) {
	cacheKey := "user:1002"

	// 1. 先查缓存
	val, err := bootstrap.GetJSON[map[string]interface{}](s.Service.rdbCache, ctx, cacheKey)
	if err == nil {
		return val, nil
	}

	// 2. 缓存 miss -> 调 client
	result, err := s.userClient.GetRandomUser(ctx)

	if err != nil {
		return nil, err
	}

	// 3. 写缓存
	_ = s.rdbCache.SetJSON(ctx, cacheKey, result, 10*time.Second)

	return result, nil
}

// 没抽离的http
// func HttpUserList2(ctx context.Context) (any, error) {
// 	var u interface{}

// 	val, err := bootstrap.GetJSON[map[string]interface{}](s.Service.rdbCache, ctx, "user:1002")

// 	if err == nil {
// 		return val, nil
// 	}

// 	// 没缓存去获取数据保存
// 	if err != nil {
// 		// 1. 创建实例
// 		client := httpclient.New("https://randomuser.me/api", 30*time.Second)
// 		// 2. 尝试调用 (这里如果不报错，说明代码没问题，是你调用处的环境问题)
// 		// 注意：这里只是测试编译，不一定会运行成功
// 		var result interface{}

// 		fmt.Printf("httpclient 进行请求！\n")
// 		// 如果这里 IDE 报错说没有 Get 方法，说明上面的第 1、2 点有问题
// 		err := client.Get(context.Background(), "", &result)
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = s.rdbCache.SetJSON(ctx, "user:1002", result, 10*time.Second)

// 		if err != nil {
// 			fmt.Println("❌ Redis 写入失败:", err)
// 		}
// 		fmt.Printf("httpclient 请求成功！接收到的数据")

// 		// 方法存在，编译通过
// 		u = &result
// 	}
// 	return u, nil
// }
