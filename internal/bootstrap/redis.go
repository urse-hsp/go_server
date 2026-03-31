package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// var RDB *redis.Client

func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	println("✅ Redis 连接成功")

	return rdb
}

// func demo() {
// 	ctx := context.Background()

// 	// 3. 设置缓存 (Key, Value, 过期时间)
// 	// 比如：缓存用户信息，5分钟后过期
// 	err := RDB.Set(ctx, "user:1001", "张三", 5*time.Minute).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 4. 获取缓存
// 	val, err := RDB.Get(ctx, "user:1001").Result()
// 	if err != nil {
// 		if err == redis.Nil {
// 			fmt.Println("键不存在")
// 		}
// 		panic(err)
// 	}
// 	fmt.Println("从 Redis 获取到的值:", val) // 输出: 张三
// }

// func Set(key string, val interface{}, ttl time.Duration) error {
// 	return RDB.Set(context.Background(), key, val, ttl).Err()
// }

// func Get(key string) (string, error) {
// 	return RDB.Get(context.Background(), key).Result()
// }

// func SetJSON(key string, value interface{}, ttl time.Duration) error {
// 	data, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}
// 	return RDB.Set(context.Background(), key, data, ttl).Err()
// }

// func GetJSON(key string, dest interface{}) error {
// 	val, err := RDB.Get(context.Background(), key).Result()
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal([]byte(val), dest)
// }
