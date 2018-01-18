package utils

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	//定义常量

	//RedisClient redis连接池
	RedisClient *redis.Pool
)

func init() {
	initConfig() //初始化配置信息
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     Conf.RedisMaxidle,
		MaxActive:   Conf.RedisMaxActive,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Conf.RedisHost)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", Conf.RedisDb)
			return c, nil
		},
	}
}

//RedisLoginKey 登入信息key格式
func RedisLoginKey(key string) string {
	return "api:login_member:" + key
}
