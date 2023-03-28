package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"jarvis-backend/service"
	"net/http"
)

var Pod pod

type pod struct{}

//GetPods 获取pod列表
func (p *pod) GetPods(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
		Cluster    string `form:"cluster"`
	})
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	data, err := service.Pod.GetPods(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})
}

//GetPodDetail 获取pod详情
func (p *pod) GetPodDetail(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取详情列表
	data, err := service.Pod.GetPodDetail(client, params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//测试更新pod的使用
	//byte, _ := json.Marshal(data)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg":  "获取pod详情成功",
	//	"data": string(byte),
	//})
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})
}

//DeletePod 删除pod
func (p *pod) DeletePod(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
		Cluster   string `json:"cluster"`
	})
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，删除pod
	err = service.Pod.DeletePod(client, params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod成功",
		"data": nil,
	})
}

//UpdatePod 更新pod
func (p *pod) UpdatePod(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
	})
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，删除pod
	err = service.Pod.UpdatePod(client, params.Namespace, params.Content)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})
}

//GetPodContainer 获取pod容器
func (p *pod) GetPodContainer(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
		Cluster   string `form:"cluster"`
	})
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，删除pod
	data, err := service.Pod.GetPodContainer(client, params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器成功",
		"data": data,
	})
}

//GetPodLog 获取pod容器日志
func (p *pod) GetPodLog(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		ContainerName string `form:"container_name"`
		PodName       string `form:"pod_name"`
		Namespace     string `form:"namespace"`
		Cluster       string `form:"cluster"`
	})
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.Bind(params); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(params.Cluster)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，删除pod
	data, err := service.Pod.GetPodLog(client, params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod容器日志成功",
		"data": data,
	})
}
