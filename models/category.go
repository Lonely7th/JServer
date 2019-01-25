package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type JCategroy struct {
	Id      int
	Cid     string `orm:"size(64)"`
	Title   string `orm:"size(64)"`
	Index   int
	Release bool // 是否展示
}

func init() {
	orm.RegisterModel(new(JCategroy))
}

func GetJCategroyList() *[]JCategroy {
	o := orm.NewOrm()
	categroy := new([]JCategroy)
	_, err := o.QueryTable("j_categroy").Filter("release", 1).OrderBy("index").All(categroy)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return categroy
}
