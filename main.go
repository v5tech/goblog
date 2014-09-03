package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"goblog/models"
	_ "goblog/routers"
)

func init() {
	RegisterDB()
}

// 注册数据库信息
func RegisterDB() {
	url := beego.AppConfig.String("url")
	orm.RegisterModel(new(models.Topic), new(models.Category), new(models.User))
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", url, 30) //注册数据库
}

func main() {
	orm.Debug = true
	beego.SetLogFuncCall(true)            //开启日志输出
	beego.SessionOn = true                //开启Session
	orm.RunSyncdb("default", false, true) //设置创建表结构
	beego.Run()
}
