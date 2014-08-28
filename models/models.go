package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// 注册数据库信息
func RegisterDB() {
	url := beego.AppConfig.String("url")
	orm.RegisterModel(new(Topic), new(Category), new(User))
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", url, 30) //注册数据库
}

// User结构体
type User struct {
	Id        int64     `orm:"auto" form:"id"`
	Username  string    `orm:"size(20)" valid:"Required" form:"username"`
	Password  string    `orm:"size(32)" valid:"Required" form:"password"`
	Email     string    `orm:"size(32)" valid:"Required" form:"email"`
	Registed  time.Time `orm:"index"`    //注册时间
	Registeip string    `orm:"size(20)"` //注册ip
	Lastlogin time.Time `orm:"index"`    //最后登录时间
	Loginip   string    `orm:"size(20)"` //最后登录ip
	Nickname  string    `valid:"Required" orm:"size(20)" form:"nickname"`
	Uuid      string    `orm:"size(64)"` //用于验证找回密码,标识用户身份
	Exprise   string    //找回密码失效时间
}

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

// Category结构体
type Category struct {
	Id           int64     `orm:"auto" form:"id"`
	CategoryName string    `orm:"size(30);index" form:"categoryname"`
	Created      time.Time `orm:"index"`
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
	category := &Category{CategoryName: categoryName, Created: time.Now().Local()}
	//category.CategoryName = categoryName
	//category.Created = time.Now().Local()
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
func GetCategoryByName(categoryName string) bool {
	o := orm.NewOrm()
	category := &Category{}
	err := o.QueryTable("category").Filter("category_name", categoryName).One(category)
	if err != nil {
		beego.Error("获取文章分类失败:" + err.Error())
		return false
	}
	if category != nil {
		return true //存在
	}
	return false
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

// 获取所有的文章
func GetAllTopics() ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0, 10)
	_, err := o.QueryTable("topic").OrderBy("-created").All(&topics) // 查询所有的文章 按降序排列
	return topics, err
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

// 用户注册
func RegisterUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

// 检查用户名是否存在
func CheckUser(username string) bool {
	o := orm.NewOrm()
	return o.QueryTable("user").Filter("username", username).Exist()
}

// 用户登录
func Login(user *User) *User {
	u := &User{}
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("username", user.Username).Filter("password", user.Password).One(u)
	if err != nil {
		beego.Error("用户登录验证失败:" + err.Error())
		return nil
	}
	return u
}

// 根据用户名和uuid更新密码

func UpdatePassWord(u *User) bool {

	o := orm.NewOrm()
	p, err := o.Raw("UPDATE user SET password = ? where username = ? and uuid = ? ").Prepare()
	_, err = p.Exec(u.Password, u.Username, u.Uuid)
	if err != nil {
		beego.Error(err.Error())
		return false
	}
	defer p.Close()
	return true //修改成功

}

// 密码找回时更新用户信息
func UpdateUser(u *User) bool {
	//用于找回密码时更新uuid和找回密码失效时间
	o := orm.NewOrm()
	p, err := o.Raw(" UPDATE user SET uuid = ?,exprise = ? where username = ? ").Prepare()
	_, err = p.Exec(u.Uuid, u.Exprise, u.Username)
	if err != nil {
		beego.Error(err.Error())
		return false
	}
	defer p.Close()
	return true //修改成功
}

// 修改用户信息
func UserModify(user *User) bool {
	o := orm.NewOrm()
	_, err := o.Update(user)
	if err != nil {
		beego.Error("修改用户信息失败!" + err.Error())
		return false
	}
	return true
}

//根据用户名和uuid查询用户信息
func QueryUserByUsernameAndUUID(username, uuid string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("username", username).Filter("uuid", uuid).One(user)
	if err != nil {
		beego.Error("根据用户名和uuid查询用户失败！" + err.Error())
		return nil
	}
	return user
}

// 根据用户名和电子邮件查询用户是否存在
func CheckUserExists(username, email string) bool {
	o := orm.NewOrm()
	return o.QueryTable("user").Filter("username", username).Filter("email", email).Exist()
}

// 根据用户名获取用户信息
func GetUserInfo(username string) (*User, error) {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("username", username).One(user)
	if err != nil {
		beego.Error("根据用户名获取用户信息失败:" + err.Error())
		return nil, err
	}
	return user, nil
}

// 删除用户
func DeleteUser(id int64) bool {
	o := orm.NewOrm()
	user := User{Id: id}
	err := o.Read(&user)
	if err != nil {
		beego.Error("读取用户信息失败:" + err.Error())
		return false
	}
	_, err = o.Delete(&user)
	if err != nil {
		beego.Error("删除用户失败:" + err.Error())
		return false
	}
	return true
}
