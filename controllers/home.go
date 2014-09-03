package controllers

import (
	"github.com/astaxie/beego"
	"goblog/models"
	"goblog/utils"
)

type HomeController struct {
	beego.Controller
}

// 分页设置
func (this *HomeController) SetPaginator(per int, nums int64) *utils.Paginator {
	p := utils.NewPaginator(this.Ctx.Request, per, nums)
	this.Data["paginator"] = p
	return p
}

// 首页
func (this *HomeController) Get() {

	beego.ReadFromRequest(&this.Controller) //解析flash数据

	checkAccountSession(&this.Controller)

	this.Data["IsHome"] = true

	qs := models.QuerySeter("topic") //获取orm.QuerySeter

	// 可以在此处加上查询过滤条件

	// qs = qs.Filter("category", "git")

	count, _ := models.CountObjects(qs) //查询总记录数

	pager := this.SetPaginator(10, count) //设置分页信息页大小、总记录数

	qs = qs. /*Filter("category", "git").*/ OrderBy("-created").Limit(10, pager.Offset()).RelatedSel() //执行分页

	//var topics []models.Topic 声明一个slice 使用此方式声明的slice在使用的时候若需要取指针的地址需要加取地址&运算符

	topics := new([]models.Topic) //该方式创建的topics已经指向指针的地址了,使用该指针的地址时,不需要再加取地址运算符

	models.ListObjects(qs, topics) //查询结果带分页

	categories := models.GetAllCategory() //查询所有的分类

	this.Data["Topics"] = topics
	this.Data["Categories"] = categories

	this.TplNames = "index.html"
}
