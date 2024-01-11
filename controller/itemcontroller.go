package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"jiansan_go_project/service"
	"net/http"
)

type ItemController struct{}

func NewItemController() ItemController {
	return ItemController{}
}

func (a ItemController) GetItemInfoByName(c *gin.Context) {
	//定义一个User变量
	//item := &model.ITEM{}
	//将调用后端的request请求中的body数据根据json格式解析到User结构变量中
	//err := c.BindJSON(item)
	//下面这个方法可以获得参数中的Name字段
	name := c.DefaultQuery("Name", "")
	//将name传给service层的GetInfoFromDbByName法，获取物品信息
	if res, err := service.GetItemInfoFromDbByName(name); err == nil {
		if data, marshalerr := json.Marshal(res); marshalerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": marshalerr.Error()})
		} else {
			c.String(http.StatusOK, fmt.Sprintf("%s", string(data)))
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
