package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"iHome/src/utils"
	"log"
)

var MySQLPool *gorm.DB

var (
	UserNotExistErr = errors.New("user not exist")
	UserPasswordErr = errors.New("password not match")
)

// InitMySQL 初始化 MySQL 表，并返回 gorm 的 db 对象
func InitMySQL() (*gorm.DB, error) {
	//sql.Open()
	dsn := "root:lgh438@tcp(127.0.0.1:3306)/ihome?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 解决建立表后表名带复数的问题
		},
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(new(User), new(House), new(Facility), new(HouseImage), new(OrderHouse), new(Area))
	if err != nil {
		return nil, err
	}
	/*db.DB().SetMaxIdleConns(10)
	db.DB().SetConnMaxLifetime(100)*/
	MySQLPool = db
	return db, nil
}

// LoginJudgement 在 MySQL 中查找以判断用户是否存在
func LoginJudgement(phone, pwd string) (string, error) {
	var user User
	fmt.Println(phone)
	err := MySQLPool.Where("mobile = ?", phone).Take(&user).Error
	if err == nil {
		log.Println("存在该用户")
		cryptPhone, err := utils.AesCtrCrypt([]byte(pwd))
		if err != nil {
			log.Println(err)
			return "", err
		}

		if hex.EncodeToString(cryptPhone) != user.Password_hash {
			// 用户输入密码错误
			return "", UserPasswordErr
		}
		// 用户存在且密码正确返回 nil
		return user.Name, nil
	} else {
		log.Println("不存在该用户")
		return "", UserNotExistErr
	}
}

// GetUserInfo 在 MySQL 中获取用户信息
func GetUserInfo(userName string) (*User, error) {
	var user User
	err := MySQLPool.Where("name = ?", userName).Take(&user).Error
	return &user, err
}

// UpdateUserName 更新用户名
func UpdateUserName(newName, olderName string) error {
	return MySQLPool.Model(new(User)).Where("name = ?", olderName).Update("name", newName).Error
}
