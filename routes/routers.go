package routes

import (
	"github.com/gin-gonic/gin"
	"jiansan_go_project/controller"
	"jiansan_go_project/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

type Result struct {
	Ctx *gin.Context
}

type ResultCont struct {
	Code int         `json:"code"` // 自增
	Msg  string      `json:"msg"`  //
	Data interface{} `json:"data"`
}

var files []string

func SetRouter() *gin.Engine {
	r := gin.Default()
	r.NoRoute(HandleNotFound)
	r.NoMethod(HandleNotFound)
	r.Use(Recover)
	workDir, _ := os.Getwd()
	filepath.Walk(workDir+"/assets/template", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	r.LoadHTMLFiles(files...)
	//r.LoadHTMLGlob("/assets/template/*.html")
	r.Static("/assets/bootstrap", workDir+"/assets/bootstrap")

	//no login
	r.GET("/admin_login.html", GetTemplate)
	r.GET("/admin_regist.html", GetTemplate)
	r.GET("/login_regist.html", GetTemplate)

	needLimitGroup := r.Group("/sign")
	//needLimit.Use(CommonBlacklist()) // 黑名单
	//needLimit.Use(CommonRateLimit()) // 频率控制
	{
		u := controller.NewUserController()
		needLimitGroup.POST("/sign_in", u.SignIn) // 登录
		needLimitGroup.POST("/sign_up", u.SignUp) // 注册
	}

	//物品item路由组
	itemGroup := r.Group("/item")
	{
		//查询物品信息
		i := controller.NewItemController()
		itemGroup.GET("/items", i.GetItemInfoByName)
	}
	//账号account路由组
	accountGroup := r.Group("/account")
	{
		//查询物品信息
		g := controller.NewAccountController()
		accountGroup.GET("/accounts", g.GetAccountInfoByID)
	}

	// needLogin 以下接口需要登录态才可访问
	needLoginGroup := r.Group("/user")
	needLoginGroup.Use(middleware.NeedLogin())
	//needLoginGroup.Use(UserBlacklist())
	//needLoginGroup.Use(UserRateLimit())
	{
		u := controller.NewUserController()
		needLoginGroup.GET("/update_password_page.html", GetTemplate) // 更新密码页
		needLoginGroup.GET("/home.html", GetTemplate)                 // 用户首页
		needLoginGroup.POST("/update_passwd", u.UpdatePasswd)         // 更新密码
		needLoginGroup.POST("/sign_out", u.SignOut)                   // 登出
	}
	return r
}

func HandleNotFound(c *gin.Context) {
	NewResult(c).Error(404, "资源未找到")
	return
}

func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			log.Printf("panic: %v\n", r)
			debug.PrintStack()
			NewResult(c).Error(500, "服务器内部错误")
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

func NewResult(ctx *gin.Context) *Result {
	return &Result{Ctx: ctx}
}

// 返回失败
func (r *Result) Error(code int, msg string) {
	res := ResultCont{}
	res.Code = code
	res.Msg = msg
	res.Data = gin.H{}
	r.Ctx.JSON(http.StatusOK, res)
	r.Ctx.Abort()
}

func GetTemplate(c *gin.Context) {
	path := c.Request.URL.Path
	arr := strings.Split(path, "/")
	html := arr[len(arr)-1]
	c.HTML(http.StatusOK, html, nil)
}

//router.LoadHTMLFiles(files...)
