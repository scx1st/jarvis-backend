package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"jarvis-backend/config"
	"jarvis-backend/controller"
	"jarvis-backend/db"
	"jarvis-backend/middle"
	"jarvis-backend/service"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	//初始化gin对象
	r := gin.Default()
	//初始化数据库
	db.Init()
	//初始化 k8s client
	service.K8s.Init()
	//跨域中间件注册
	r.Use(middle.Cors())
	//r.Use(middle.JWTAuth())
	// 初始化路由规则
	controller.Router.InitApiRouter(r)
	//启动task
	go func() {
		service.Event.WatchEventTask("Cluster-1")
	}()
	//gin server启动
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
	//开启线程，否则ListenAndServe会无限循环，走不下去
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Info("listen: %s\n", err)
		}
	}()
	//等待中断信号，优雅关闭所有server
	//创建一个信号通道 quit，用于接收系统信号。
	quit := make(chan os.Signal)
	//将操作系统中的中断信号(os.Interrupt)转化为 quit 通道中的消息，以便在程序运行时能够捕获到该信号并做出相应的处理。
	signal.Notify(quit, os.Interrupt)
	<-quit
	//设置ctx超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//cancel 用于释放ctx
	defer cancel()
	//关闭gin
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Gin Server关闭异常：", err)
	}
	logger.Info("Gin Server退出成功")
	//关闭db
	if err := db.Close(); err != nil {
		logger.Fatal("DB关闭异常：", err)
	}
}
