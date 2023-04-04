package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"jarvis-backend/service"
	"net/http"
)

var Ingress ingress

type ingress struct{}

//CreateIngress 建deployment
func (i *ingress) CreateIngress(ctx *gin.Context) {
	var (
		ingressCreate = new(service.IngressCreate)
		err           error
	)
	//绑定参数
	//form格式使用ctx.Bind方法，json格式使用ctx.ShouldBindJSON方法
	if err := ctx.ShouldBindJSON(ingressCreate); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败, %v\n", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败, %v\n", err),
			"data": nil,
		})
		return
	}
	//获取client
	client, err := service.K8s.GetClient(ingressCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//调用service方法，获取列表
	err = service.Ingress.CreateIngress(client, ingressCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			//"code": 90500, //业务状态
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建Ingress成功",
		"data": nil,
	})
}
