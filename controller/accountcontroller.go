package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"jiansan_go_project/service"
	"net/http"
)

type AccountController struct{}

func NewAccountController() AccountController {
	return AccountController{}
}

func (a AccountController) GetAccountInfoByID(c *gin.Context) {
	//下面这个方法可以获得参数中的id字段
	id := c.DefaultQuery("id", "")
	//将id传给service层的GetItemInfoFromDbByName法，获取物品信息
	if res, err := service.GetAccountInfoFromDbByID(id); err == nil {
		if data, marshalerr := json.Marshal(res); marshalerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": marshalerr.Error()})
		} else {
			c.String(http.StatusOK, fmt.Sprintf("%s", string(data)))
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
