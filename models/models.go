package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func RegisterDB() {
	url := beego.AppConfig.String("url")
	orm.RegisterModel(new(Topic), new(Category), new(User))
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", url, 30) //注册数据库
}

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
}

type Topic struct {
	Id         int64     `orm:"auto" form:"id"`
	title      string    `valid:"Required" orm:"size(20);index" form:"title"`
	Content    string    `valid:"Required" orm:"size(5000);index" form:"content"`
	Category   string    `valid:"Required" orm:"size(100);index" form:"category"`
	Views      int64     //浏览次数
	ReplyCount int64     //回复次数
	Author     string    `orm:"size(32)"`
	Created    time.Time `orm:"index"`
	Updated    time.Time `orm:"index"`
}

type Category struct {
	Id           int64  `orm:"auto" form:"id"`
	CategoryName string `orm:"size(30);index" index`
}

/**
 * 发表博客
 */

func AddTopic(topic *Topic) error {
	o := orm.NewOrm()
	num, err := o.Insert(topic)
	if err != nil {
		beego.Error(err)
		return err
	}
	if num == 1 {
		return nil
	}
	return nil
}

/**
 * 用户注册
 */
func RegisterUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

/**
 * 检查用户名是否存在
 */
func CheckUser(username string) bool {
	o := orm.NewOrm()
	user := &User{}
	err := o.QueryTable(new(User)).Filter("username", username).One(user)
	if err != nil {
		beego.Error(err)
		return false
	}
	return user != nil
}

/**
 * 用户登录
 */
func Login(user *User) *User {
	u := new(User)
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("username", user.Username).Filter("password", user.Password).One(u)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return u
}

/**
 * 修改用户信息
 */
func UserModify(user *User) bool {
	o := orm.NewOrm()

	u := &User{}

	err := o.QueryTable("user").Filter("username", user.Username).Filter("password", user.Password).One(u) //使用用户名和密码查询用户信息

	if err != nil {
		return false
	}

	u.Lastlogin = user.Lastlogin //修改用户最后登录时间
	u.Loginip = user.Loginip     //修改用户最后登录ip

	num, err := o.Update(u)

	if err != nil {
		return false
	}

	return num == 1
}

/**
 * 根据用户名获取用户信息
 */
func GetUserInfo(username string) (*User, error) {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("username", username).One(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

/**
 * 删除用户
 */

func DeleteUser(id int64) bool {
	o := orm.NewOrm()
	user := User{Id: id}
	err := o.Read(&user)
	if err != nil {
		beego.Error(err)
		return false
	}
	num, err := o.Delete(&user)
	if err != nil {
		beego.Error(err)
		return false
	}
	beego.Info(num)
	return num == 1
}
