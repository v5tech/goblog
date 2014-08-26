package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
	"strconv"
	"time"
)

type TopicController struct {
	beego.Controller
}

/**
 * 根据文章id删除文章
 */
func (this *TopicController) DeleteTopic() {
	if checkAccountSession(&this.Controller) { //验证用户是否已登录
		id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
		if err != nil {
			beego.Error("转换文章id出错")
			return
		}
		if !models.DeleteTopic(id) {
			beego.Error("删除文章出错")
			return
		}
		this.Redirect("/", 302)
		return
	} else {
		this.SetSession("Error", "您尚未登录,请登录!")
		this.Redirect("/login", 302) //跳转到登录页
		return
	}

}

/**
 * 根据文章id查看文章
 */
func (this *TopicController) ViewTopic() {
	if this.GetSession("user") != nil {
		user := this.GetSession("user").(*models.User) //从Session中获取用户信息
		this.Data["Nickname"] = user.Nickname
		this.Data["Username"] = user.Username
		this.Data["IsLogin"] = true
		this.Data["IsTopic"] = true
	}
	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		beego.Error("转换文章id出错")
		return
	}
	topic := models.ViewTopicById(id)
	this.Data["Topic"] = topic
	this.TplNames = "view_topic.html"
}

/**
 * 跳转到新增页面
 */
func (this *TopicController) Add() {
	if checkAccountSession(&this.Controller) { //验证用户是否已登录
		this.Data["IsTopic"] = true
		this.Data["Title"] = "添加文章"
		this.TplNames = "add_topic.html"
		return
	} else {
		this.SetSession("Error", "您尚未登录,请登录!")
		this.Redirect("/login", 302) //跳转到登录页
		return
	}
}

/**
 * 添加文章内容
 */
func (this *TopicController) AddTopic() {
	if checkAccountSession(&this.Controller) { //验证用户是否已登录
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
			topic.Author = this.GetSession("user").(*models.User).Nickname   //设置作者
			topic.Username = this.GetSession("user").(*models.User).Username //设置用户名
		}

		if !models.GetCategoryByName(topic.Category) { //查询该分类是否存在
			//不存在,则保存该分类
			if !models.AddCategory(topic.Category) {
				beego.Error("保存分类出错:" + err.Error())
			}
		}

		err = models.AddTopic(topic)

		if err != nil {
			beego.Error(err)
			this.Redirect("/", 302)
			return
		}

		this.Redirect("/", 302)
		return
	} else {
		this.SetSession("Error", "您尚未登录,请登录!")
		this.Redirect("/login", 302) //跳转到登录页
		return
	}

}

func (this *TopicController) Post() {

}
