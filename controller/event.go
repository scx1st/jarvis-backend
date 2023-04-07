package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"jarvis-backend/service"
	"net/http"
)

var Event event

type event struct{}

//GetList 获取Event列表
func (*event) GetList(c *gin.Context) {
	params := new(struct {
		Name    string `form:"name"`
		Cluster string `form:"cluster"`
		Page    int    `form:"page"`
		Limit   int    `form:"limit"`
	})
	if err := c.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Event.GetList(params.Name, params.Cluster, params.Page, params.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取Event列表成功",
		"data": data,
	})
}
