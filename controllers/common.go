package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"goblog/models"
)

// 判断用户是否已登录 从cookie中验证
func checkAccountCookie(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("username")
	if err != nil {
		return false
	}

	username := ck.Value

	ck, err = ctx.Request.Cookie("password")

	if err != nil {
		return false
	}

	password := ck.Value

	user := &models.User{Username: username, Password: password}

	return models.Login(user) != nil
}

// 判断用户是否已登录 从Session中验证
func checkAccountSession(ctl *beego.Controller) bool {
	if ctl.GetSession("user") != nil {
		user := ctl.GetSession("user").(*models.User) //从Session中获取用户信息
		ctl.Data["Nickname"] = user.Nickname
		ctl.Data["Username"] = user.Username
		ctl.Data["IsLogin"] = true
		return true
	} else {
		return false
	}
}

// 生成该字符串的MD5串
func MD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
