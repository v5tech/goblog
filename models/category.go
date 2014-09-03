package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

// Category结构体
type Category struct {
	Id           int64     `orm:"auto" form:"id"`
	CategoryName string    `orm:"size(30);index" form:"categoryname"`
	Created      time.Time `orm:"index"`
	Updated      time.Time `orm:"index"`
}

// 查询所有的分类
func GetAllCategory() []*Category {
	o := orm.NewOrm()
	categorys := make([]*Category, 0, 10)
	_, err := o.QueryTable("category").All(&categorys)
	if err != nil {
		beego.Error("获取所有的分类失败:" + err.Error())
		return nil
	}
	return categorys
}

// 删除文章分类
func DeleteCategory(id int64) bool {
	o := orm.NewOrm()
	category := &Category{}
	err := o.QueryTable("category").Filter("id", id).One(category)
	if err != nil {
		beego.Error("根据文章分类id查询文章分类失败:" + err.Error())
		return false
	}
	_, err = o.Delete(category)
	if err != nil {
		beego.Error("根据文章分类id删除文章分类失败:" + err.Error())
		return false
	}
	return true
}

// 保存分类
func AddCategory(categoryName string) bool {
	o := orm.NewOrm()
	category := &Category{CategoryName: categoryName, Created: time.Now().Local(), Updated: time.Now().Local()}
	_, err := o.Insert(category)
	if err != nil {
		beego.Error("保存分类失败:" + err.Error())
		return false
	}
	return true
}

// 修改文章分类
func ModifyCategory(category *Category) bool {
	o := orm.NewOrm()
	cat := &Category{}
	err := o.QueryTable("category").Filter("id", category.Id).One(cat)
	if err != nil {
		beego.Error("获取文章分类失败:" + err.Error())
	}
	cat.CategoryName = category.CategoryName
	cat.Updated = time.Now().Local()
	_, err = o.Update(cat)
	if err != nil {
		beego.Error("修改文章分类失败:" + err.Error())
		return false
	}
	return true
}

// 根据分类id查询分类
func GetCategoryById(id int64) *Category {
	o := orm.NewOrm()
	category := &Category{}
	err := o.QueryTable("category").Filter("id", id).One(category)
	if err != nil {
		beego.Error("获取文章分类失败:" + err.Error())
		return nil
	}
	return category
}

// 根据分类名称查询分类
func GetCategoryByName(categoryName string) *Category {
	o := orm.NewOrm()
	category := &Category{}
	err := o.QueryTable("category").Filter("category_name", categoryName).One(category)
	if err != nil {
		beego.Error("获取文章分类失败:" + err.Error())
		return nil
	}
	return category
}
