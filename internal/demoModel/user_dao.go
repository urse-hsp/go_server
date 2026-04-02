package demo

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/bootstrap"
	"go-server/internal/dao"
	"go-server/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	Create(ctx context.Context, username string, password string) (*model.User, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, info model.User, id uint) (*model.User, error)
	GetDetail(ctx context.Context, id uint) (*model.User, error)
	GetList(ctx context.Context) ([]model.User, error)
	GetLists(ctx context.Context, page, pageSize int) ([]model.User, int64, error)
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
	var user model.User

	err := r.DB(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= 根据关键字查询 =================

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.DB(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= 创建 =================

func (r *userRepository) Create(ctx context.Context, username string, password string) (*model.User, error) {
	// 判断是否已存在
	_, err := r.GetByUsername(ctx, username)

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

	if err := r.DB(ctx).Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ================= 更新 =================

func (r *userRepository) Update(ctx context.Context, user model.User, id uint) (*model.User, error) {
	result := r.DB(ctx).Model(&model.User{}).
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
	if err := r.DB(ctx).First(&updatedUser, id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// ================= 删除 =================

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB(ctx).Where("id = ?", id).Delete(&model.User{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// ================= 全部列表 =================

func (r *userRepository) GetList(ctx context.Context) ([]model.User, error) {
	var users []model.User

	if err := r.DB(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// ================= 分页列表 =================

func (r *userRepository) GetLists(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User

	db := r.DB(ctx).Model(&model.User{})

	// 分页
	total, err := dao.Paginate(db, &users, page, pageSize)

	return users, total, err
}
