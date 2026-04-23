package dao

import (
	"context"
	"fmt"
	demodto "go-server/internal/dto/demo"
	"go-server/internal/model"

	"gorm.io/gorm"
)

type DemoRepository interface {
	Create(ctx context.Context, data *model.Demo) (*model.Demo, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, data *model.Demo, id uint) (*model.Demo, error)
	GetDetail(ctx context.Context, id uint) (*model.Demo, error)
	GetList(ctx context.Context, q demodto.RequestQuery) ([]model.Demo, error)
	GetPageList(ctx context.Context, q demodto.RequestPageQuery) ([]model.Demo, int64, error)

	buildQuery(ctx context.Context, q demodto.RequestQuery) *gorm.DB
	GetByKeyWhere(ctx context.Context, name string) (*model.Demo, error)
}

func NewDemoRepository(
	r *Repository,
) DemoRepository {
	return &demoRepository{
		Repository: r,
	}
}

type demoRepository struct {
	*Repository
}

// ================= 根据ID查询 =================

func (r *demoRepository) GetDetail(ctx context.Context, id uint) (*model.Demo, error) {
	var data model.Demo

	err := r.DB(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ================= 根据关键字查询 =================

func (r *demoRepository) GetByKeyWhere(ctx context.Context, name string) (*model.Demo, error) {
	var data model.Demo

	err := r.DB(ctx).Where("username = ?", name).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// ================= 创建 =================

func (r *demoRepository) Create(ctx context.Context, data *model.Demo) (*model.Demo, error) {

	if err := r.DB(ctx).Create(data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// ================= 更新 =================

func (r *demoRepository) Update(ctx context.Context, data *model.Demo, id uint) (*model.Demo, error) {
	if err := r.DB(ctx).
		Model(&model.Demo{}).
		Where("id = ?", id).
		Updates(data).Error; err != nil {
		return nil, err
	}

	return r.GetDetail(ctx, id)
}

// ================= 删除 =================
// model有 DeletedAt 的情况下，会自动变成软删除，实际数据库记录并未删除，而是更新了 DeletedAt 字段的值
// 没有 DeletedAt 字段的情况下，会真正删除记录

func (r *demoRepository) Delete(ctx context.Context, id uint) error {
	result := r.DB(ctx).Where("id = ?", id).Delete(&model.User{})
	// r.DB(ctx).Unscoped().Where("id = ?", id).Delete(&model.User{}) // 物理删除，直接删除记录，不更新 DeletedAt 字段

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("id不存在")
	}

	return nil
}

// ================= 全部列表 =================

func (r *demoRepository) GetList(ctx context.Context, q demodto.RequestQuery) ([]model.Demo, error) {
	var users []model.Demo

	db := r.buildQuery(ctx, q)

	// r.DB(ctx).Unscoped().Find(&users) // 物理查询，包含已软删除的记录

	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// ================= 分页列表 =================

func (r *demoRepository) GetPageList(ctx context.Context, q demodto.RequestPageQuery) ([]model.Demo, int64, error) {
	var users []model.Demo

	db := r.buildQuery(ctx, q.RequestQuery)

	total, err := Paginate(db, &users, q.Page, q.PageSize)

	return users, total, err
}

// ================= 公共查询 =================
func (r *demoRepository) buildQuery(ctx context.Context, q demodto.RequestQuery) *gorm.DB {
	db := r.DB(ctx).Model(&model.Demo{})

	if q.Query != "" {
		db = db.Where("username LIKE ?", "%"+q.Query+"%")
	}

	if q.Status != nil {
		db = db.Where("status = ?", *q.Status)
	}

	if q.Type != nil {
		db = db.Where("type = ?", *q.Type)
	}

	return db
}
