package dao

import (
	"context"
	"go-server/pkg/log"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	// rdbCache *bootstrap.RDBCache
	//mongo  *mongo.Client
	logger *log.Logger
}

// Repository 负责数据库连接和事务管理
func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
	// rdbCache *bootstrap.RDBCache,
	// mongo *mongo.Client,
) *Repository {
	return &Repository{
		db: db,
		// rdbCache: rdbCache,
		//mongo:  mongo,
		logger: logger,
	}
}

// 事物接口。将数据库事务的执行逻辑抽象化
type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

const ctxTxKey = "TxKey"

// DB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.DB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}
