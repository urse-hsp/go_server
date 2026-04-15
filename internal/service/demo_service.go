package service

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/dao"
	demodto "go-server/internal/dto/demo"
	"go-server/internal/model"

	"gorm.io/gorm"
)

type DemoService interface {
	Create(ctx context.Context, req demodto.CreateRequest) (*model.Demo, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, id uint, req demodto.UpdateRequest) (*model.Demo, error)
	GetDetail(ctx context.Context, id uint) (*model.Demo, error)
	GetList(ctx context.Context, q demodto.RequestQuery) ([]model.Demo, error)
	GetPageList(ctx context.Context, q demodto.RequestPageQuery) ([]model.Demo, int64, error)
}

func NewDemoService(
	service *Service,
	Repo dao.DemoRepository,
) DemoService {
	return &demoService{
		Repo:    Repo,
		Service: service,
	}
}

type demoService struct {
	*Service
	Repo dao.DemoRepository
}

// ================= 创建 =================

func (s *demoService) Create(ctx context.Context, req demodto.CreateRequest) (*model.Demo, error) {
	// 判断是否已存在
	data, err := s.Repo.GetByKeyWhere(ctx, req.Username)
	if err == nil && data != nil {
		return nil, fmt.Errorf("数据已存在")
	}

	data_ := &model.Demo{}

	return s.Repo.Create(ctx, data_)
}

// ================= 删除 =================

func (s *demoService) Delete(ctx context.Context, id uint) error {
	// 直接调用 DAO（DAO 已经判断是否存在）
	return s.Repo.Delete(ctx, id)
}

// ================= 更新 =================

func (s *demoService) Update(ctx context.Context, id uint, req demodto.UpdateRequest) (*model.Demo, error) {
	data, err := s.Repo.GetDetail(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("数据不存在")
		}
		return nil, err
	}

	// AssignIfNotNil(&data.Username, req.Username)

	return s.Repo.Update(ctx, data, id)
}

// ================= 获取 =================

func (s *demoService) GetDetail(ctx context.Context, id uint) (*model.Demo, error) {
	return s.Repo.GetDetail(ctx, id)
}

// ================= 全部列表 =================

func (s *demoService) GetList(ctx context.Context, q demodto.RequestQuery) ([]model.Demo, error) {
	return s.Repo.GetList(ctx, q)
}

// ================= 分页列表 =================

func (s *demoService) GetPageList(ctx context.Context, q demodto.RequestPageQuery) ([]model.Demo, int64, error) {
	return s.Repo.GetPageList(ctx, q)
}
