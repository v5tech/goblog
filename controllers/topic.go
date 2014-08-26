package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
	"time"
)

type TopicController struct {
	beego.Controller
}

/**
 * 跳转到新增页面
 */
func (this *TopicController) Add() {
	this.Data["IsTopic"] = true
	this.Data["IsLogin"] = checkAccount(this.Ctx)
	if this.GetSession("user") != nil {
		user := this.GetSession("user").(*models.User) //从Session中获取用户信息
		this.Data["Nickname"] = user.Nickname
		this.Data["Username"] = user.Username
	}
	this.Data["Title"] = "添加文章"
	this.TplNames = "add_topic.html"
}

/**
 * 添加文章内容
 */
func (this *TopicController) AddTopic() {
	topic := &models.Topic{}
	err := this.ParseForm(topic)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	topic.Created = time.Now() //设置创建时间
	topic.Updated = time.Now() //设置更新时间

	if this.GetSession("user") != nil {
		topic.Author = this.GetSession("user").(*models.User).Nickname //设置作者
	}

	err = models.AddTopic(topic)

	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	this.Redirect("/", 302)
}

func (this *TopicController) Post() {

}
