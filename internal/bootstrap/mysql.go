package bootstrap

import (
	"context"
	"go-server/pkg/log"
	"go-server/pkg/zapgorm2"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	logger := zapgorm2.New(l.Logger)
	driver := conf.GetString("data.db.user.driver")
	dsn := conf.GetString("data.db.user.dsn")

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger,
		})
	// case "postgres":
	// 	db, err = gorm.Open(postgres.New(postgres.Config{
	// 		DSN:                  dsn,
	// 		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	// 	}), &gorm.Config{})
	// case "sqlite":
	// 	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		panic("unknown db driver")
	}
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	// Connection Pool config
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	println("✅ MySQL 连接成功")
	return db
}
