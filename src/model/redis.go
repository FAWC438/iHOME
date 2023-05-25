package model

import "github.com/gomodule/redigo/redis"

// RedisPool 创建全局连接池
var RedisPool redis.Pool

// init
//
//	@Description: 利用包初始化函数，初始化 redis 连接池
func init() {
	RedisPool = redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
		MaxIdle:         20,
		MaxActive:       50,
		IdleTimeout:     60,
		Wait:            false, // 判断到达最大活跃连接数量后，如果有新连接请求，是否等待连接池有新的可用连接
		MaxConnLifetime: 60 * 5,
	}
}
