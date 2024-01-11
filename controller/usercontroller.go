package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jiansan_go_project/service"
	"jiansan_go_project/utils"
	"time"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

// SignIn 登录
func (u UserController) SignIn(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) <= 0 {
		utils.WriteResponseWithCode(c, "请填入用户名", nil, 0)
		return
	}

	if len(password) <= 0 {
		utils.WriteResponseWithCode(c, "请填入密码", nil, 0)
		return
	}

	user, err := service.GetUserFromDbByName(username)
	if err != nil {
		utils.WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	//if !has {
	//	utils.WriteResponseWithCode(c, "账号不存在", nil, 0)
	//	return
	//}

	if user.Passwd != utils.GenMD5WithSalt(password, user.Salt) {
		utils.WriteResponseWithCode(c, "密码不正确", nil, 0)
		return
	}

	//账密验证通过，生成session
	session := service.NewSession(user)
	err = session.Store() // 存储session到redis
	if err != nil {
		utils.WriteResponseWithCode(c, "登录失败，请重试", nil, 0)
		return
	}

	// 更新登录IP和登录时间
	var updates = map[string]interface{}{
		"last_login_unix": time.Now().Unix(),
		"last_login_ip":   c.ClientIP(),
	}
	err = service.UpdateUserToDB(session.Username, updates)
	if err != nil {
		fmt.Println("DBUpdateUser err: ", err)
	}

	//登录成功
	c.SetCookie("SESSION", session.ID, 0, "", "", false, true)
	c.SetCookie("USERNAME", session.Username, 0, "", "", false, true)
	c.SetCookie("UID", string(session.UID), 0, "", "", false, true)
	utils.WriteResponseWithCode(c, "", nil, 0)
	//c.Redirect(http.StatusFound, Conf.Common.HomePage) // 重定向跳回首页
}

func GetCtxUser(c *gin.Context) *service.Session {
	sess, exist := c.Get("USER")
	if !exist {
		return nil
	}
	return sess.(*service.Session)
}

// SignOut 登出
func (u UserController) SignOut(c *gin.Context) {
	sess := GetCtxUser(c)
	if sess == nil {
		utils.WriteResponseWithCode(c, "尚未登录", nil, 0)
		return
	}

	err := sess.Del()
	if err != nil {
		utils.WriteResponseWithCode(c, "注销失败，请重试", nil, 0)
		return
	}

	c.SetCookie("SESSION", "", 0, "", "", false, true)
	c.SetCookie("USERNAME", "", 0, "", "", false, true)
	c.SetCookie("UID", "", 0, "", "", false, true)

	utils.WriteResponseWithCode(c, "", nil, 0)
}

// SignUp 注册
func (u UserController) SignUp(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	retryPassword := c.PostForm("retry_password")
	email := c.PostForm("email")

	if retryPassword != password {
		utils.WriteResponseWithCode(c, "两次输入的密码不一致", nil, 0)
		return
	}

	err := utils.IsValidName(username)
	if err != nil {
		utils.WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	err = utils.IsValidEmail(email)
	if err != nil {
		utils.WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	err = utils.IsValidPasswd(password)
	if err != nil {
		utils.WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	if service.UserExists("Name", username) {
		utils.WriteResponseWithCode(c, "用户名已存在", nil, 0)
		return
	}
	if service.UserExists("Email", email) {
		utils.WriteResponseWithCode(c, "邮箱已被注册", nil, 0)
		return
	}

	// 将新用户插入数据库
	err = service.AddUserToDB(username, password, email)
	if err != nil {
		utils.WriteResponseWithCode(c, "注册失败，请稍后重新注册", nil, 0)
		return
	}

	utils.WriteResponseWithCode(c, "", nil, 0)
}

// UpdatePasswd 更新密码
func (u UserController) UpdatePasswd(c *gin.Context) {
	sess := GetCtxUser(c)
	if sess == nil {
		utils.WriteResponseWithCode(c, "尚未登录", nil, 0)
		return
	}

	retryPassword := c.PostForm("retry_password")
	password := c.PostForm("password")

	if retryPassword != password {
		utils.WriteResponseWithCode(c, "两次输入的密码不一致", nil, 0)
		return
	}

	err := utils.IsValidPasswd(password)
	if err != nil {
		utils.WriteResponseWithCode(c, err.Error(), nil, 0)
		return
	}

	user, err := service.GetUserFromDbByName(sess.Username)
	if err != nil {
		utils.WriteResponseWithCode(c, "用户不存在", nil, 0)
		return
	}

	var updates = map[string]interface{}{
		"passwd": utils.GenMD5WithSalt(password, user.Salt),
	}

	err = service.UpdateUserToDB(sess.Username, updates)
	if err != nil {
		utils.WriteResponseWithCode(c, "修改密码失败，请重试", nil, 0)
		return
	}

	//utils.WriteResponseWithCode(c, "修改密码成功", nil, 0)
	utils.WriteResponseWithCode(c, "", nil, 0)
}
