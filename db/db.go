package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/wonderivan/logger"
	"jarvis-backend/config"
	"jarvis-backend/model"
)

var (
	isInit bool //是否已经初始化
	GORM   *gorm.DB
	err    error
)

//Init db初始化函数
func Init() {
	//判断是否已初始化
	if isInit {
		return
	}
	//组装db连接的数据
	//parseTime是查询结果是否自动解析为时间
	//loc是mysql的时区设置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPwd,
		config.DbHost,
		config.DbPort,
		config.DbName)
	GORM, err = gorm.Open(config.DbType, dsn)
	//打印sql语句
	GORM.LogMode(config.LogMode)
	//开启连接池
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxIdleTime(config.MaxLifeTime)

	isInit = true
	GORM.AutoMigrate(model.Event{})
	logger.Info("连接数据库成功")
}

func Close() error {
	logger.Info("关闭数据库成功")
	return GORM.Close()
}
