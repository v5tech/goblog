package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dchest/captcha"
	"goblog/models"
	"strconv"
	"time"
)

type UserController struct {
	beego.Controller
}

// 生成验证码
func (this *UserController) NewCaptcha() {

	captchaStr := captcha.New()                                                                  //生成唯一的id
	this.SetSession("captchaStr", captchaStr)                                                    //将该字符串放入session中
	captcha.WriteImage(this.Ctx.ResponseWriter, captchaStr, captcha.StdWidth, captcha.StdHeight) //将验证码输出到浏览器
	return

}

// 登出
func (this *UserController) Logout() {

	if checkAccountSession(&this.Controller) { //判断用户是否登录
		this.Ctx.SetCookie("username", "", -1, "/") //设置cookie失效
		this.Ctx.SetCookie("password", "", -1, "/") //设置cookie失效
		flash := beego.NewFlash()
		flash.Notice("用户" + this.GetSession("user").(*models.User).Nickname + "已成功退出!")
		flash.Store(&this.Controller)
		this.DelSession("user") //从Session中移除当前登录的用户信息
		this.Redirect("/", 302) //重定向到主页
		return
	}
	this.Redirect("/", 302) //重定向到主页
	return

}

// 登录页
func (this *UserController) Login() {

	beego.ReadFromRequest(&this.Controller)
	this.Data["Title"] = "用户登录"
	this.TplNames = "login.html" //登录页

}

// 用户注册页
func (this *UserController) Register() {

	beego.ReadFromRequest(&this.Controller)
	this.Data["Title"] = "用户注册"
	this.TplNames = "register.html"

}

// 用户注册表单提交
func (this *UserController) RegisterAction() {
	flash := beego.NewFlash()
	user := &models.User{}
	err := this.ParseForm(user)
	if err != nil {
		beego.Error("用户注册失败:" + err.Error())
		flash.Error("注册用户失败!")
		flash.Store(&this.Controller)
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
		flash.Error("验证码不正确!")
		flash.Store(&this.Controller)
		this.DelSession("captchaStr")   //从session中清空
		this.Redirect("/register", 302) //验证码不正确,重定向到登录页
		return
	} else {
		isExists := models.CheckUser(user.Username) //判断该用户是否已经存在
		if isExists {
			flash.Error("该用户已存在!")
			flash.Store(&this.Controller)
			this.Redirect("/register", 302) //该用户已存在,重定向到注册页
			return
		} else {
			err = models.RegisterUser(user) //用户注册
			if err != nil {
				flash.Error("注册用户失败!")
				flash.Store(&this.Controller)
				this.Redirect("/register", 302) //验证码不正确,重定向到注册页
				return
			}
		}

	}
	flash.Notice("注册成功!")
	flash.Store(&this.Controller)
	this.Redirect("/login", 302) //注册成功,重定向到登录页
	return
}

// 用户登录表单提交
func (this *UserController) LoginAction() {
	flash := beego.NewFlash()
	user := &models.User{}
	err := this.ParseForm(user)

	if err != nil {
		beego.Error("用户登录失败:" + err.Error())
		flash.Error("用户登录失败!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //登录失败,重定向到登录页
		return
	}

	user.Password = MD5(user.Password) //将密码以MD5加密存储
	captchaCode := this.Input().Get("captcha")

	//判断验证码是否正确
	if !captcha.VerifyString(this.GetSession("captchaStr").(string), captchaCode) {
		flash.Error("验证码不正确!")
		flash.Store(&this.Controller)
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
			this.Ctx.SetCookie("username", user.Username, maxAge, "/") //设置cookie
			this.Ctx.SetCookie("password", user.Password, maxAge, "/") //设置cookie

			u.Lastlogin = time.Now()        //设置最后登录时间
			u.Loginip = this.Ctx.Input.IP() //获取客户端IP

			if !models.UserModify(u) { //用户登录成功后更新最后登录时间
				beego.Error("更新用户最后登录时间失败")
				flash.Error("更新用户最后登录时间失败!")
				flash.Store(&this.Controller)
			}

			this.SetSession("user", u) //将用户信息存放到Session中
			flash.Notice("登录成功!")
			flash.Store(&this.Controller)
			this.Redirect("/", 302) //登录成功
			return
		} else {
			flash.Error("用户名或密码不正确!")
			flash.Store(&this.Controller)
			this.Redirect("/login", 302) //登录失败,重定向到登录页
			return
		}
	}

}

// 删除用户
func (this *UserController) DeleteUser() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) {
		id := this.Ctx.Input.Param(":id") // /user/:id 删除用户的路径
		idNum, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			beego.Error("删除用户失败:" + err.Error())
			flash.Error("删除用户失败!")
			flash.Store(&this.Controller)
			return
		}
		models.DeleteUser(idNum) //删除用户
		this.DelSession("user")  //清空session
		flash.Notice("用户删除成功!")
		flash.Store(&this.Controller)
		this.Redirect("/", 302) //重定向到主页
		return
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}

}

// 根据用户名查看用户详细信息
func (this *UserController) GetUserInfo() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) {
		username := this.Ctx.Input.Param(":username") //user/:username
		user, err := models.GetUserInfo(username)
		if err != nil {
			beego.Error("获取用户信息失败:" + err.Error())
			flash.Error("获取用户信息失败!")
			flash.Store(&this.Controller)
			this.Redirect("/", 302)
			return
		}
		if this.GetSession("user") != nil {
			user := this.GetSession("user").(*models.User) //从Session中获取用户信息
			this.Data["Nickname"] = user.Nickname
			this.Data["Username"] = user.Username
			this.Data["IsLogin"] = true
		}
		this.Data["User"] = user
		this.TplNames = "user.html"
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}

}
