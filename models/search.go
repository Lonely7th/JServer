package models

import (
	"github.com/astaxie/beego/orm"
)

//查询JNote (基础搜索)
func SearchJNote2List(key string, page int) *[]JNote {
	noteList := new([]JNote)

	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond1 := cond.Or("label_title1__contains", key).Or("label_title2__contains", key).Or("label_title3__contains", key)

	qs := o.QueryTable("j_note")
	_, err := qs.SetCond(cond1).OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(noteList)
	if err != nil {
		return nil
	}
	return noteList
}

//查询JNote (高级搜索)
func SearchJNote2List2(key string) *[]JNote {
	return nil
}
