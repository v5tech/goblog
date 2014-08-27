package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
	"strconv"
)

type CategoryController struct {
	beego.Controller
}

// 查询所有的分类信息
func (this *CategoryController) GetAllCategory() {

	if checkAccountSession(&this.Controller) {
		beego.ReadFromRequest(&this.Controller) //解析flash数据
		this.Data["Title"] = "分类列表"
		this.Data["IsCategory"] = true
		this.Data["Categories"] = models.GetAllCategory()
		this.TplNames = "categories.html"
	} else {
		flash := beego.NewFlash()
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302)
		return
	}

}

// 添加文章分类
func (this *CategoryController) AddCategory() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) {
		category := &models.Category{}
		err := this.ParseForm(category)
		if err != nil {
			beego.Error("收集表单失败:" + err.Error())
			flash.Error("收集表单失败!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //重定向到分类列表页面
			return
		}

		if !models.GetCategoryByName(category.CategoryName) { //查询该分类是否存在
			//不存在,则保存该分类
			if !models.AddCategory(category.CategoryName) {
				beego.Error("添加文章分类失败:" + err.Error())
				flash.Error("添加文章分类失败!")
				flash.Store(&this.Controller)
			} else {
				flash.Notice("添加成功!")
				flash.Store(&this.Controller)
				this.Redirect("/category", 302) //添加成功,重定向到文章分类列表页面
			}
		} else {
			flash.Error("该分类已存在!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //添加成功,重定向到文章分类列表页面
		}
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //尚未登录,重定向到登录页面
		return
	}
}

// 修改文章分类
func (this *CategoryController) ModifyCategory() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) {
		// id, err := strconv.ParseInt(this.Input().Get("id"), 10, 64) //获取表单中的id
		// if err != nil {
		// 	beego.Error("获取文章分类Id失败:" + err.Error())
		// 	return
		// }
		category := &models.Category{}
		err := this.ParseForm(category)
		if err != nil {
			beego.Error("获取表单信息失败:" + err.Error())
			flash.Error("获取表单信息失败!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //重定向到分类列表页面
			return
		}
		if models.ModifyCategory(category) {
			flash.Notice("修改成功!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //修改成功,重定向到文章分类列表页面
			return
		} else {
			flash.Error("修改失败!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //重定向到分类列表页面
			return
		}
	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //尚未登录,重定向到登录页面
		return
	}

}

// 删除文章分类
func (this *CategoryController) DeleteCategory() {
	flash := beego.NewFlash()
	if checkAccountSession(&this.Controller) {
		id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64) //获取文章分类Id
		if err != nil {
			beego.Error("获取文章分类Id失败:" + err.Error())
			flash.Error("获取文章分类Id失败!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //重定向到分类列表页面
			return
		}

		if models.DeleteCategory(id) {
			flash.Notice("删除文章分类成功!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //重定向到分类列表页面
			return
		} else {
			flash.Error("删除文章分类失败!")
			flash.Store(&this.Controller)
			this.Redirect("/category", 302) //重定向到分类列表页面
			return
		}

	} else {
		flash.Error("您尚未登录,请登录!")
		flash.Store(&this.Controller)
		this.Redirect("/login", 302) //尚未登录,重定向到登录页面
		return
	}
}
