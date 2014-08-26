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

	checkAccountSession(&this.Controller)

	topics, err := models.GetAllTopics()  //查询所有的文章
	categories := models.GetAllCategory() //查询所有的分类

	if err != nil {
		beego.Error(err)
	}

	this.Data["Topics"] = topics
	this.Data["Categories"] = categories

	this.TplNames = "index.html"
}
