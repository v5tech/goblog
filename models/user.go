package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "github.com/astaxie/beego/validation"
	"time"
)

// User结构体
type User struct {
	Id        int64     `orm:"auto" form:"id"`
	Username  string    `orm:"size(20);MaxSize(20)" valid:"Required" form:"username"`
	Password  string    `orm:"size(32)" valid:"Required" form:"password"`
	Email     string    `orm:"size(32)" valid:"Required;Email;MaxSize(100)" form:"email"`
	Registed  time.Time `orm:"index"`    //注册时间
	Registeip string    `orm:"size(20)"` //注册ip
	Lastlogin time.Time `orm:"index"`    //最后登录时间
	Loginip   string    `orm:"size(20)"` //最后登录ip
	Nickname  string    `valid:"Required" orm:"size(20)" form:"nickname"`
	Uuid      string    `orm:"size(64)"` //用于验证找回密码,标识用户身份
	Exprise   string    //找回密码失效时间
}

// struct 实现接口 validation.ValidFormer
// 当 StructTag 中的测试都成功时，将会执行 Valid 函数进行自定义验证
// func (u *User) Valid(v *validation.Validation) {
// 	if strings.Index(u.Username, "admin") != -1 {
// 		// 通过 SetError 设置 username 的错误信息，HasErrors 将会返回 true
// 		v.SetError("username", "用户名不能为admin")
// 	}
// }

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
