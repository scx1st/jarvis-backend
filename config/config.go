package config

import "time"

const (
	WsAddr     = "0.0.0.0:8081"
	ListenAddr = "0.0.0.0:9090"
	//KubeConfigs 一个config，为了验证多集群用了两个一样的
	KubeConfigs = `{"Cluster-1":"D:\\client-config\\config","Cluster-2":"D:\\client-config\\config"}`
	//PodLogTailLine 查看容器日志时，显示的tail行数 tail -n 5000
	PodLogTailLine = 5000
	//DbType 数据库配置
	DbType = "mysql"
	DbHost = "localhost"
	DbPort = 3306
	DbName = "test_db"
	DbUser = "root"
	DbPwd  = "password"
	//LogMode 打印mysql debug sql日志
	LogMode = false
	//MaxIdleConns 连接池配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
	//helm上传路径
	UploadPath = "D:\\client-config\\config\\upload"
	//账号密码
	AdminUser = "admin"
	AdminPwd  = "123456"
)
