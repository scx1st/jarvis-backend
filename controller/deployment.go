package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"jarvis-backend/service"
	"net/http"
)

var Deployment deployment

type deployment struct{}

//GetDeployments 获取Deployment列表
func (d *deployment) GetDeployments(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	data, err := service.Deployment.GetDeployments(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Deployment列表成功",
		"data": data,
	})
}

//GetDeploymentDetail 获取Deployment详情
func (d *deployment) GetDeploymentDetail(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
		Cluster        string `form:"cluster"`
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	data, err := service.Deployment.GetDeploymentDetail(client, params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//测试更新Deployment使用
	//byte, _ := json.Marshal(data)
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"msg":  "获取Deployment详情成功",
	//	"data": string(byte),
	//})
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Deployment详情成功",
		"data": data,
	})
}

//DeleteDeployment 删除Deployment
func (d *deployment) DeleteDeployment(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	err = service.Deployment.DeleteDeployment(client, params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "删除Deployment成功",
		"data": nil,
	})
}

//UpdateDeployment 更新Deployment
func (d *deployment) UpdateDeployment(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		Content   string `json:"content"`
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	err = service.Deployment.UpdateDeployment(client, params.Namespace, params.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "更新Deployment成功",
		"data": nil,
	})
}

//ScaleDeployment 调整Deployment副本数
func (d *deployment) ScaleDeployment(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		ScaleNum       int    `json:"scale_num"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	data, err := service.Deployment.ScaleDeployment(client, params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "调整Deployment副本数成功",
		"data": data,
	})
}

//RestartDeployment 重启Deployment
func (d *deployment) RestartDeployment(c *gin.Context) {
	//1.接受参数，绑定参数
	//匿名结构体，get请求为form格式，其他请求为json格式
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		Cluster        string `json:"cluster"`
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	err = service.Deployment.RestartDeployment(client, params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "重启Deployment成功",
		"data": nil,
	})
}

//CreateDeployment 创建Deployment
func (d *deployment) CreateDeployment(c *gin.Context) {
	//1.接受参数，绑定参数
	var (
		deployCreate = new(service.DeployCreate)
		err          error
	)
	//form格式使用c.Bind方法， json格式使用c.ShouldBindJSON方法
	if err := c.ShouldBindJSON(deployCreate); err != nil {
		logger.Error(fmt.Sprintf("绑定参数失败,%v\n", err))
		c.JSON(http.StatusBadGateway, gin.H{
			"msg":  fmt.Sprintf("绑定参数失败,%v\n", err),
			"data": nil,
		})
	}
	//2.获取client
	client, err := service.K8s.GetClient(deployCreate.Cluster)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	//3.调用service方法，获取列表
	err = service.Deployment.CreateDeployment(client, deployCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "创建Deployment成功",
		"data": nil,
	})
}
