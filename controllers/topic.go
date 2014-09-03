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

// 根据分类名称查询文章
func (this *TopicController) ViewTopicByCategoryName() {

	category := this.Ctx.Input.Param(":category")

	topics := models.QueryTopicByCategoryName(category)

	categories := models.GetAllCategory() //查询所有的分类

	this.Data["Topics"] = topics
	this.Data["Categories"] = categories
	this.Data["Category"] = category

	this.TplNames = "view_topic_cat.html"
}

// 根据文章id删除文章
func (this *TopicController) DeleteTopic() {

	flash := beego.NewFlash()

	if checkAccountSession(&this.Controller) { //验证用户是否已登录
		id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
		if err != nil {
			beego.Error("转换文章id失败")
			flash.Error("删除文章失败!")
			flash.Store(&this.Controller)
			return
		}
		if !models.DeleteTopic(id) {
			beego.Error("删除文章失败")
			flash.Error("删除文章失败!")
			flash.Store(&this.Controller)
			return
		}
		this.Redirect("/", 302) //删除成功回首页
		return
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}

}

// 修改文章页面
func (this *TopicController) ModifyTopic() {

	beego.ReadFromRequest(&this.Controller)
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) { //验证用户是否已登录
		id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
		if err != nil {
			beego.Error("获取文章id失败")
			flash.Error("获取文章id失败!")
			flash.Store(&this.Controller)
			return
		}
		topic := models.ViewTopicById(id)
		this.Data["Topic"] = topic
		this.TplNames = "modify_topic.html"
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}
}

// 修改文章Action
func (this *TopicController) ModifyTopicAction() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) { //验证用户是否已登录

		topic := &models.Topic{}
		err := this.ParseForm(topic)

		if err != nil {
			beego.Error("收集表单数据失败!" + err.Error())
		} else {

			category := models.GetCategoryByName(topic.Category) //查询该分类名称查询该分类是否存在

			if category == nil {
				//不存在,则添加分类
				if !models.AddCategory(topic.Category) {
					beego.Error("添加文章分类失败:" + err.Error())
					flash.Error("添加文章分类失败!")
					flash.Store(&this.Controller)
					this.Redirect("/topic/modify/"+strconv.FormatInt(topic.Id, 10), 302) //重新定向到修改页面
					return
				}
			} else { //存在则修改分类
				category.CategoryName = topic.Category
				if !models.ModifyCategory(category) {
					beego.Error("修改文章分类失败:" + err.Error())
					flash.Error("修改文章分类失败!")
					flash.Store(&this.Controller)
					this.Redirect("/topic/modify/"+strconv.FormatInt(topic.Id, 10), 302) //重新定向到修改页面
					return
				}
			}

			//修改文章
			err = models.ModifyTopic(topic)
			if err != nil {
				beego.Error("修改文章失败!" + err.Error())
				flash.Error("修改文章失败!")
				flash.Store(&this.Controller)
				this.Redirect("/topic/modify/"+strconv.FormatInt(topic.Id, 10), 302) //重新定向到修改页面
				return
			}
			flash.Notice("修改文章成功!")
			flash.Store(&this.Controller)
			this.Redirect("/topic/view/"+strconv.FormatInt(topic.Id, 10), 302) //修改成功重定向到查看页面
		}
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}
}

// 根据文章id查看文章
func (this *TopicController) ViewTopic() {

	beego.ReadFromRequest(&this.Controller)
	flash := beego.NewFlash()

	if this.GetSession("user") != nil {
		user := this.GetSession("user").(*models.User) //从Session中获取用户信息
		this.Data["Nickname"] = user.Nickname
		this.Data["Username"] = user.Username
		this.Data["IsLogin"] = true
		this.Data["IsTopic"] = true
	}
	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		beego.Error("获取文章id失败" + err.Error())
		flash.Error("获取文章id失败!")
		flash.Store(&this.Controller)
		return
	}
	topic := models.ViewTopicById(id)
	this.Data["Topic"] = topic
	this.TplNames = "view_topic.html"
}

// 跳转到新增页面
func (this *TopicController) Add() {

	beego.ReadFromRequest(&this.Controller)

	flash := beego.NewFlash()

	if checkAccountSession(&this.Controller) { //验证用户是否已登录
		this.Data["IsTopic"] = true
		this.Data["Title"] = "添加文章"
		this.TplNames = "add_topic.html"
		return
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}
}

// 添加文章内容
func (this *TopicController) AddTopic() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) { //验证用户是否已登录
		topic := &models.Topic{}
		err := this.ParseForm(topic)
		if err != nil {
			beego.Error("添加文章内容失败:" + err.Error())
			flash.Error("添加文章失败!")
			flash.Store(&this.Controller)
			this.Redirect("/", 302)
			return
		}

		topic.Created = time.Now().Local() //设置创建时间
		topic.Updated = time.Now().Local() //设置更新时间

		if this.GetSession("user") != nil {
			topic.Author = this.GetSession("user").(*models.User).Nickname   //设置作者
			topic.Username = this.GetSession("user").(*models.User).Username //设置用户名
		}

		if nil == models.GetCategoryByName(topic.Category) { //查询该分类是否存在
			//不存在,则保存该分类
			if !models.AddCategory(topic.Category) {
				beego.Error("添加文章分类失败:" + err.Error())
				flash.Error("添加文章分类失败!")
				flash.Store(&this.Controller)
			}
		}

		err = models.AddTopic(topic)

		if err != nil {
			beego.Error("添加文章失败:" + err.Error())
			flash.Error("添加文章失败!")
			flash.Store(&this.Controller)
			this.Redirect("/topic/add", 302)
			return
		}

		this.Redirect("/", 302) //添加成功到文章列表页
		return
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //跳转到登录页
		return
	}

}
