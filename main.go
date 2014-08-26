package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"goblog/models"
	_ "goblog/routers"
)

func init() {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	beego.SetLogFuncCall(true)            //开启日志输出
	beego.SessionOn = true                //开启Session
	orm.RunSyncdb("default", false, true) //设置创建表结构
	beego.Run()
}
