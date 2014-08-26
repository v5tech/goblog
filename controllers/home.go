package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
)

type HomeController struct {
	beego.Controller
}

//首页
func (this *HomeController) Get() {
	this.Data["IsHome"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	if this.GetSession("user") != nil {
		user := this.GetSession("user").(*models.User) //从Session中获取用户信息
		this.Data["Nickname"] = user.Nickname
		this.Data["Username"] = user.Username
	}
	this.TplNames = "index.html"
}
