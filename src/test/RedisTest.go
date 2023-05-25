package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 测试 redis 客户端
func main() {
	redisConn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(redisConn)
	str, err := redis.String(redisConn.Do("set", "name", "kevin"))
	if err != nil {
		return
	}
	fmt.Println(str)

}
