package bootstrap

import (
	"go-server/pkg/log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	rdb *redis.Client
	//mongo  *mongo.Client
	logger *log.Logger
}

// Repository 负责数据库连接和事务管理
func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
	rdb *redis.Client,
	// mongo *mongo.Client,
) *Repository {
	return &Repository{
		db:  db,
		rdb: rdb,
		//mongo:  mongo,
		logger: logger,
	}
}
