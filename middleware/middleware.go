package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jiansan_go_project/service"
	"jiansan_go_project/utils"
	"net/http"
)

// NeedLogin 必须登录的请求，从session读user写入context
func NeedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		err, sess := service.GetSession(c)
		if err == nil && sess != nil {
			fmt.Println("NeedLogin in login, ", sess)
			c.Set("USER", sess)
			sess.Store() // 已登录，每次请求都会续期
			c.Next()
			return
		} else {
			// 未登录
			utils.WriteResponseWithCode(c, "未登录", nil, http.StatusUnauthorized)
			//c.Redirect(http.StatusFound, Conf.Common.EnterPage)
			c.Abort()
			return
		}
	}
}
