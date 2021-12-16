package db

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type singleton struct {
	Db *gorm.DB
}

// 单例模式，获取gorm连接信息
// 失败返回nil
var instance *singleton 
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {

		if db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{}); err == nil {

			sqlDB, _ := db.DB()
			// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
			sqlDB.SetMaxIdleConns(11)

			// SetMaxOpenConns 设置打开数据库连接的最大数量。
			sqlDB.SetMaxOpenConns(51)

			// SetConnMaxLifetime 设置了连接可复用的最大时间。
			sqlDB.SetConnMaxLifetime(time.Hour)

			dbStatus, _ := db.DB()
			fmt.Println(dbStatus.Stats())

			instance = &singleton{
				Db: db,
			}

		} else {
			panic(err)
		}

	})
	return instance
}
