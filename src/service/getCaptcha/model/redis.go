package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func SaveImgCode(code, uuid string) error {
	redisConn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(redisConn)

	// 设置图片验证码在 redis 中5分钟过期
	_, err = redis.String(redisConn.Do("setex", uuid, 60*3, code))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
