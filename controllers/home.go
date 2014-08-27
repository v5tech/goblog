package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
)

type HomeController struct {
	beego.Controller
}

// 首页
func (this *HomeController) Get() {

	beego.ReadFromRequest(&this.Controller) //解析flash数据

	checkAccountSession(&this.Controller)

	this.Data["IsHome"] = true

	topics, err := models.GetAllTopics() //查询所有的文章

	categories := models.GetAllCategory() //查询所有的分类

	if err != nil {
		beego.Error("查询所有的文章出错:" + err.Error())
	}

	this.Data["Topics"] = topics
	this.Data["Categories"] = categories
	this.TplNames = "index.html"
}
