package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

// People gorm 创建的字段一定要大写，模型定义见 https://gorm.io/zh_CN/docs/models.html
// 其中标签的更改只能在创建表时生效
type People struct {
	//Id   int // 主键
	gorm.Model
	Name string `gorm:"size:50"`
	Age  int    `gorm:"not null"`
}

var dbConn *gorm.DB

//// TableName 指定表名，同样能解决建立表后带复数的问题，见 https://www.zhyea.com/2022/05/24/gorm-plural-table-name.html
//func (People) TableName() string {
//	return "people"
//}

// InsertData 插入数据
func InsertData(data *People) error {
	return dbConn.Create(data).Error
}

// QueryData 查询数据，见 https://gorm.io/zh_CN/docs/query.html
func QueryData() error {
	return nil
}

func main() {
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// ? 后参数用于设置字符集与 MySQL 时区，否则默认 +0 时区
	dsn := "root:lgh438@tcp(127.0.0.1:3306)/world?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 解决建立表后带复数的问题
		},
	})
	if err != nil {
		fmt.Println("connect to mysql error: ", err)
		return
	}
	dbConn = db

	sqlDB, err := dbConn.DB()
	if err != nil {
		fmt.Println("DB object error: ", err)
		return
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		fmt.Println("DB use error: ", err)
		return
	}

	// 通过结构体在数据库中自动创建表（自动迁移功能）
	//err = dbConn.AutoMigrate(new(People))
	//if err != nil {
	//	fmt.Println("AutoMigrate error: ", err)
	//	return
	//}

	testData := People{
		Name: "kevin",
		Age:  22,
	}

	// 插入数据
	err = InsertData(&testData)
	if err != nil {
		fmt.Println("insert error", err)
		return
	}

	fmt.Println("运行结束...")
}
