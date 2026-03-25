package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"go-demo-server/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() *redis.Client {
	RDSConfig := config.Conf.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:     RDSConfig.Addr,
		Password: RDSConfig.Password,
		DB:       RDSConfig.DB,
	})

	// 🔥 检查连接是否成功
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		panic("❌ Redis 连接失败: " + err.Error())
	}

	println("✅ Redis 连接成功")

	// demo()
	newRDB := RDB
	return newRDB
}

func demo() {
	ctx := context.Background()

	// 3. 设置缓存 (Key, Value, 过期时间)
	// 比如：缓存用户信息，5分钟后过期
	err := RDB.Set(ctx, "user:1001", "张三", 5*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	// 4. 获取缓存
	val, err := RDB.Get(ctx, "user:1001").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("键不存在")
		}
		panic(err)
	}
	fmt.Println("从 Redis 获取到的值:", val) // 输出: 张三
}

// func Set(key string, val interface{}, ttl time.Duration) error {
// 	return RDB.Set(context.Background(), key, val, ttl).Err()
// }

// func Get(key string) (string, error) {
// 	return RDB.Get(context.Background(), key).Result()
// }

func SetJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RDB.Set(context.Background(), key, data, ttl).Err()
}

func GetJSON(key string, dest interface{}) error {
	val, err := RDB.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}
