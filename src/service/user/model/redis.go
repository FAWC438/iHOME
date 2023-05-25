package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"user/utils"
)

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

// CheckImgCode
//
//	@Description: 在 redis 中检查图片验证码输入是否正确
//	@param uuid
//	@param imgCode
//	@return bool
func CheckImgCode(uuid, imgCode string) bool {
	conn := RedisPool.Get()
	// 如果使用 go-redis 而非 redigo 则连接池的连接无需手动关闭
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	str, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return imgCode == str
}

// SaveSmsCode
//
//	@Description: 在 redis 中存储短信验证码，手机号码用 AES 加密算法加密
//	@param phoneNum
//	@param code
//	@return error
func SaveSmsCode(phoneNum, code string) error {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	// 用 AES 加密算法加密手机号码
	cryptPhoneNum, err := utils.AesCtrCrypt([]byte(phoneNum))
	if err != nil {
		return err
	}
	// 在 redis 中存储手机号码及其验证码
	_, err = conn.Do("setex", cryptPhoneNum, 3*60, code)
	if err != nil {
		return err
	}

	return nil
}

// CheckSmsCode
//
//	@Description: 通过获取 redis 信息校验短信验证码
//	@param phone
//	@param code
//	@return error
func CheckSmsCode(phone, code string) (error error, ok bool) {
	conn := RedisPool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("redis pool get error: ", err)
		}
	}(conn)

	cryptPhone, err := utils.AesCtrCrypt([]byte(phone))
	if err != nil {
		fmt.Println("phone crypt error: ", err)
		return err, false
	}
	trueSmsCode, err := redis.String(conn.Do("get", cryptPhone))
	if err != nil {
		fmt.Println("redis get error: ", err)
		return err, false
	}
	if trueSmsCode != code {
		return nil, false
	}

	return nil, true
}
