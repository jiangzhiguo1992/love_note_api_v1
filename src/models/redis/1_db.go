package redis

import (
	"time"

	"libs/utils"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
)

// InitPool 开启连接
func InitPool() {
	// 获取参数
	address := utils.GetConfigStr("conf", "app.conf", "redis", "redis_address")
	//redisMaxIdle := utils.GetConfigInt("conf", "app.conf", "redis", "redis_max_idle")
	//redisIdleLifeMin := utils.GetConfigInt("conf", "app.conf", "redis", "redis_idle_life_min")
	// 开始连接
	pool = &redis.Pool{
		//MaxIdle:     redisMaxIdle,
		//IdleTimeout: time.Duration(redisIdleLifeMin) * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			utils.LogWarn("redis", err)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			utils.LogWarn("redis", err)
			return err
		},
	}
}

// ClosePool
func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}

// auth 密码验证
func auth(conn redis.Conn) {
	redisPwd := utils.GetConfigStr("conf", "app.conf", "redis", "redis_pwd")
	_, err := conn.Do("AUTH", redisPwd)
	utils.LogErr("redis", err)
}
