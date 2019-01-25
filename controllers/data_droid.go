package controllers

import (
	"ApiJServer/models"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"time"
)

//Note成功率分析系统
func NoteSuccessAnalysisSystem() {
	//遍历Note表
	o := orm.NewOrm()
	noteList := new([]models.JNote)
	_, err := o.QueryTable("j_note").RelatedSel().All(noteList)
	if err == nil {
		for _, item := range *noteList {
			rand.Seed(time.Now().UnixNano())
			if item.SuccessRate == 0 {
				item.SuccessRate = rand.Intn(100)
			} else {
				a := rand.Intn(10)
				event := rand.Intn(2)
				switch event {
				case 0: //+
					if item.SuccessRate+a < 100 {
						item.SuccessRate += a
					}
					break
				case 1: //-
					if item.SuccessRate-a > 0 {
						item.SuccessRate -= a
					}
					break
				}
			}
			//更新成功率参数
			o.Update(&item)
		}
	}
}
