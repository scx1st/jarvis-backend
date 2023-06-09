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
		//登录验证
		POST("/api/login", Login.Auth).
		//集群
		GET("/api/k8s/clusters", Cluster.GetClusters).
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		//deployment操作
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		//service操作
		POST("/api/k8s/service/create", Servicev1.CreateService).
		//ingress操作
		POST("/api/k8s/ingress/create", Ingress.CreateIngress).
		//event操作
		GET("/api/k8s/events", Event.GetList).
		//allres
		GET("/api/k8s/allres", AllRes.GetAllNum).
		//helm应用商店
		GET("/api/helmstore/releases", HelmStore.ListReleases).
		GET("/api/helmstore/release/detail", HelmStore.DetailRelease).
		POST("/api/helmstore/release/install", HelmStore.InstallRelease).
		DELETE("/api/helmstore/release/uninstall", HelmStore.UninstallRelease).
		GET("/api/helmstore/charts", HelmStore.ListCharts).
		POST("/api/helmstore/chart/add", HelmStore.AddChart).
		PUT("/api/helmstore/chart/update", HelmStore.UpdateChart).
		DELETE("/api/helmstore/chart/del", HelmStore.DeleteChart).
		POST("/api/helmstore/chartfile/upload", HelmStore.UploadChartFile).
		DELETE("/api/helmstore/chartfile/del", HelmStore.DeleteChartFile)
}
