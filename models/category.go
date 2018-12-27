package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type JCategroy struct {
	Id int
	Cid string `orm:"size(64)"`
	Title string `orm:"size(64)"`
	Index string `orm:"size(64)"`
}

type RelationCaNote struct {
	Id int
	Cid *JCategroy `orm:"rel(one)"`
	Nid *JNote `orm:"rel(one)"`
}

func init() {
	orm.RegisterModel(new(JCategroy),new(RelationCaNote))
}

func GetJCategroyList() *[]JCategroy {
	o := orm.NewOrm()
	categroy := new([]JCategroy)
	_, err := o.QueryTable("j_categroy").All(categroy)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return categroy
}