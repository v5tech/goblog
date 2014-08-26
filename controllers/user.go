package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dchest/captcha"
	"goblog/models"
	"strconv"
	"time"
)

type UserController struct {
	beego.Controller
}

//生成验证码
func (this *UserController) NewCaptcha() {
	captchaStr := captcha.New()
	this.SetSession("captchaStr", captchaStr)
	captcha.WriteImage(this.Ctx.ResponseWriter, captchaStr, captcha.StdWidth, captcha.StdHeight)
	return
}

//登出
func (this *UserController) Logout() {
	this.Ctx.SetCookie("username", "", -1, "/")
	this.Ctx.SetCookie("password", "", -1, "/")
	this.DelSession("user") //从Session中移除当前登录的用户信息
	this.Redirect("/", 302) //重定向到主页
	return
}

//登录页
func (this *UserController) Login() {
	this.Data["Title"] = "用户登录"
	IsError := (this.GetSession("Error") != nil) //从Session中获取错误消息
	if IsError {
		this.Data["IsError"] = IsError
		this.Data["Error"] = this.GetSession("Error").(string)
	}
	this.DelSession("Error")     //从session中移除错误信息
	this.TplNames = "login.html" //登录页
}

//用户注册页
func (this *UserController) Register() {
	this.Data["Title"] = "用户注册"
	IsError := (this.GetSession("Error") != nil) //从Session中获取错误消息
	if IsError {
		this.Data["IsError"] = IsError
		this.Data["Error"] = this.GetSession("Error").(string)
	}
	this.DelSession("Error") //从session中移除错误信息
	this.TplNames = "register.html"
}

//用户注册表单提交
func (this *UserController) RegisterAction() {
	user := &models.User{}
	err := this.ParseForm(user)
	if err != nil {
		beego.Error(err)
		this.Redirect("/register", 302) //注册失败,重定向到注册页
		return
	}

	user.Password = MD5(user.Password)   //将密码以MD5加密存储
	user.Registed = time.Now()           //用户注册时间
	user.Lastlogin = time.Now()          //用户最后登录时间
	user.Registeip = this.Ctx.Input.IP() //用户注册的ip

	captchaCode := this.Input().Get("captcha")

	//判断验证码是否正确
	if !captcha.VerifyString(this.GetSession("captchaStr").(string), captchaCode) {
		this.SetSession("Error", "验证码不正确!")
		this.DelSession("captchaStr")   //从session中清空
		this.Redirect("/register", 302) //验证码不正确,重定向到登录页
		return
	} else {
		isExists := models.CheckUser(user.Username) //判断该用户是否已经存在
		if isExists {
			this.SetSession("Error", "该用户已存在!")
			this.Redirect("/register", 302) //该用户已存在,重定向到注册页
			return
		} else {
			err = models.RegisterUser(user) //用户注册
			if err != nil {
				this.SetSession("Error", err)
				this.Redirect("/register", 302) //验证码不正确,重定向到注册页
				return
			}
		}

	}
	this.Redirect("/", 302) //注册成功,重定向到首页
	return
}

//用户登录表单提交
func (this *UserController) LoginAction() {

	user := &models.User{}
	err := this.ParseForm(user)

	if err != nil {
		beego.Error(err)
		this.Redirect("/login", 302) //登录失败,重定向到登录页
		return
	}

	user.Password = MD5(user.Password) //将密码以MD5加密存储
	captchaCode := this.Input().Get("captcha")

	//判断验证码是否正确
	if !captcha.VerifyString(this.GetSession("captchaStr").(string), captchaCode) {
		this.SetSession("Error", "验证码不正确!")
		this.DelSession("captchaStr") //从session中清空
		this.Redirect("/login", 302)  //验证码不正确,重定向到登录页
		return
	} else {
		isAutoLogin := this.Input().Get("isAutoLogin") == "on" //是否自动登录

		u := models.Login(user) //成功返回user,失败返回nil

		if u != nil {
			maxAge := 0
			if isAutoLogin {
				maxAge = 72 * 24 * 60
			}
			this.Ctx.SetCookie("username", user.Username, maxAge, "/")
			this.Ctx.SetCookie("password", user.Password, maxAge, "/")

			u.Lastlogin = time.Now()        //设置最后登录时间
			u.Loginip = this.Ctx.Input.IP() //获取客户端IP

			if !models.UserModify(u) { //用户登录成功后更新最后登录时间
				beego.Error("更新用户最后登录时间出错")
			}

			this.SetSession("user", u) //将用户信息存放到Session中
			this.Redirect("/", 302)    //登录成功
			return
		} else {
			this.SetSession("Error", "用户名或密码不正确!")
			this.Redirect("/login", 302) //登录失败,重定向到登录页
			return
		}
	}

}

/**
 * 删除用户
 */
func (this *UserController) DeleteUser() {
	id := this.Ctx.Input.Param(":id") // /user/:id 删除用户的路径
	idNum, err := strconv.ParseInt(id, 10, 64)
	beego.Info(idNum)
	if err != nil {
		beego.Error(err)
		return
	}

	models.DeleteUser(idNum)
	this.DelSession("user")
	this.Redirect("/", 302)
	return
}

/**
 * 根据用户名查看用户详细信息
 */
func (this *UserController) GetUserInfo() {
	username := this.Ctx.Input.Param(":username") //user/:username
	user, err := models.GetUserInfo(username)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	if this.GetSession("user") != nil {
		user := this.GetSession("user").(*models.User) //从Session中获取用户信息
		this.Data["Nickname"] = user.Nickname
		this.Data["Username"] = user.Username
	}
	this.Data["User"] = user
	this.TplNames = "user.html"
}

//判断用户是否已登录 从cookie中验证
func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("username")
	if err != nil {
		return false
	}

	username := ck.Value

	ck, err = ctx.Request.Cookie("password")

	if err != nil {
		return false
	}

	password := ck.Value

	user := &models.User{Username: username, Password: password}

	return models.Login(user) != nil
}

/**
 * 生成该字符串的MD5串
 */
func MD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
