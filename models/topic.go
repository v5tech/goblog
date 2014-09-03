package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

// Topic结构体
type Topic struct {
	Id         int64     `orm:"auto" form:"id"`
	Title      string    `valid:"Required" orm:"size(100);index" form:"title"`
	Content    string    `valid:"Required" orm:"size(5000);index" form:"content"`
	Category   string    `valid:"Required" orm:"size(100);index" form:"category"`
	Views      int64     //浏览次数
	ReplyCount int64     //回复次数
	Username   string    `orm:"size(20)" valid:"Required" form:"username"`
	Author     string    `orm:"size(32)"`
	Created    time.Time `orm:"index"`
	Updated    time.Time `orm:"index"`
}

// 发表文章
func AddTopic(topic *Topic) error {
	o := orm.NewOrm()
	_, err := o.Insert(topic)
	if err != nil {
		beego.Error("发表文章失败:" + err.Error())
		return err
	}
	return nil
}

// 修改文章
func ModifyTopic(topic *Topic) error {
	o := orm.NewOrm()

	t := ViewTopicById(topic.Id) //根据id从数据库查询文章

	t.Title = topic.Title
	t.Content = topic.Content
	t.Category = topic.Category
	t.Updated = time.Now().Local()

	_, err := o.Update(t)

	if err != nil {
		beego.Error("修改文章失败:" + err.Error())
		return err
	}
	return nil
}

// 查询所有的文章列表
func GetAllTopics() (int64, []*Topic) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0, 10)
	res, err := o.QueryTable("topic").OrderBy("-created").All(&topics) // 查询所有的文章 按降序排列
	if err != nil {
		beego.Error("获取文章失败")
		return 0, nil
	}
	return res, topics
}

// 根据文章id删除文章
func DeleteTopic(id int64) bool {
	o := orm.NewOrm()
	topic := new(Topic)
	topic.Id = id
	err := o.Read(topic)
	if err != nil {
		beego.Error("查询文章失败:" + err.Error())
		return false
	}
	_, err = o.Delete(topic)
	if err != nil {
		beego.Error("删除文章失败:" + err.Error())
		return false
	}
	return true
}

// 根据文章id查看文章
func ViewTopicById(id int64) *Topic {
	o := orm.NewOrm()
	topic := &Topic{Id: id}
	err := o.Read(topic)
	if err != nil {
		beego.Error("查看文章失败:" + err.Error())
		return nil
	}
	topic.Views++
	_, err = o.Update(topic)
	if err != nil {
		beego.Error("更新文章浏览次数失败:" + err.Error())
	}
	return topic
}

// 根据CategoryName查询文章列表
func QueryTopicByCategoryName(category string) []*Topic {
	o := orm.NewOrm()
	topics := make([]*Topic, 0, 10)
	_, err := o.QueryTable("topic").Filter("category", category).OrderBy("-created").All(&topics)
	if err != nil {
		beego.Error("根据分类名称查询文章失败" + err.Error())
		return nil
	}
	return topics
}
