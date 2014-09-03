package models

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 生成该字符串的MD5串
func MD5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// 根据表名称 返回与该表对应的QuerySeter
func QuerySeter(tablename string) orm.QuerySeter {
	o := orm.NewOrm()
	return o.QueryTable(tablename)
}

// 检查结果集是否存在
func CheckIsExist(qs orm.QuerySeter, field string, value interface{}, skipId int) bool {
	qs = qs.Filter(field, value)
	if skipId > 0 {
		qs = qs.Exclude("Id", skipId)
	}
	return qs.Exist()
}

// 获取总记录数
func CountObjects(qs orm.QuerySeter) (int64, error) {
	cnt, err := qs.Count()
	if err != nil {
		beego.Error("models.CountObjects ", err)
		return 0, err
	}
	return cnt, err
}

// 获取查询列表
func ListObjects(qs orm.QuerySeter, objs interface{}) (int64, error) {
	nums, err := qs.All(objs)
	if err != nil {
		beego.Error("models.ListObjects ", err)
		return 0, err
	}
	return nums, err
}
