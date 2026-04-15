package dao

import (
	"context"
	"fmt"
	"go-server/internal/bootstrap"
	userdto "go-server/internal/dto/user"
	"go-server/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, data *model.User) (*model.User, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, data *model.User, id uint) (*model.User, error)
	GetDetail(ctx context.Context, id uint) (*model.User, error)
	GetList(ctx context.Context, q userdto.RequestQuery) ([]model.User, error)
	GetPageList(ctx context.Context, q userdto.RequestPageQuery) ([]model.User, int64, error)

	buildQuery(ctx context.Context, q userdto.RequestQuery) *gorm.DB
	GetByKeyWhere(ctx context.Context, name string) (*model.User, error)
}

func NewUserRepository(
	r *bootstrap.Repository,
) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*bootstrap.Repository
}

// ================= 根据ID查询 =================

func (r *userRepository) GetDetail(ctx context.Context, id uint) (*model.User, error) {
	var data model.User

	err := r.DB(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ================= 根据关键字查询 =================

func (r *userRepository) GetByKeyWhere(ctx context.Context, name string) (*model.User, error) {
	var data model.User

	err := r.DB(ctx).Where("username = ?", name).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ================= 创建 =================

func (r *userRepository) Create(ctx context.Context, data *model.User) (*model.User, error) {

	if err := r.DB(ctx).Create(data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// ================= 更新 =================

func (r *userRepository) Update(ctx context.Context, data *model.User, id uint) (*model.User, error) {
	if err := r.DB(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(data).Error; err != nil {
		return nil, err
	}

	return r.GetDetail(ctx, id)
}

// ================= 删除 =================

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB(ctx).Where("id = ?", id).Delete(&model.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("id不存在")
	}

	return nil
}

// ================= 全部列表 =================

func (r *userRepository) GetList(ctx context.Context, q userdto.RequestQuery) ([]model.User, error) {
	var users []model.User

	db := r.buildQuery(ctx, q)

	err := db.Find(&users).Error
	return users, err
}

// ================= 分页列表 =================

func (r *userRepository) GetPageList(ctx context.Context, q userdto.RequestPageQuery) ([]model.User, int64, error) {
	var users []model.User

	db := r.buildQuery(ctx, q.RequestQuery)

	total, err := Paginate(db, &users, q.Page, q.PageSize)

	return users, total, err
}

// ================= 公共查询 =================
func (r *userRepository) buildQuery(ctx context.Context, q userdto.RequestQuery) *gorm.DB {
	db := r.DB(ctx).Model(&model.User{})

	if q.Query != nil && *q.Query != "" {
		db = db.Where("username LIKE ?", "%"+*q.Query+"%")
	}

	return db
}
