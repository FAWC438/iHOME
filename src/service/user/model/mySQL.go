package model

import (
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"user/utils"
)

var MySQLPool *gorm.DB

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

// RegisterUserInMySQL
//
//	@Description: 在 MySQL 中存储用户注册信息
//	@param phone
//	@param pwd
//	@return exist
//	@return e
func RegisterUserInMySQL(phone, pwd string) (exist bool, e error) {
	var user User

	e = MySQLPool.Where("Mobile = ?", phone).Take(&user).Error
	if e == nil {
		exist = true
		return
	}

	exist = false
	user.Mobile = phone
	user.Name = phone // TODO: 自定义用户名
	cryptPwd, e := utils.AesCtrCrypt([]byte(pwd))
	if e != nil {
		return
	}
	user.Password_hash = hex.EncodeToString(cryptPwd)

	e = MySQLPool.Create(&user).Error
	return
}
