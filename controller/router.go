package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//Router 实例化router对象，可使用该对象点出首字母大写的方法(跨包调用)
var Router router

//定义router结构体
type router struct{}

//InitApiRouter 初始化路由，创建测试api接口
func (*router) InitApiRouter(r *gin.Engine) {
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "testApi success!",
			"data": nil,
		})
	}).
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog)
}
