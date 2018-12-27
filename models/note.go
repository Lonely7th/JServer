package models

import "github.com/astaxie/beego/orm"

type JNote struct {
	Id int
	NoteId string `orm:"size(64)"` // 唯一标识
	Content string `orm:"size(64)"` // 内容简介
	Releaser *User `orm:"rel(one)"` // OneToOne relation 发布者
	CreatTime int // 发布时间
	ResPath string `orm:"size(64)"` // 图片地址
	DisplayNum int // 展示次数
	CompleteNum int // 完成次数
	JType int // 限制类型(0.无限制 1.时间限制 2.次数限制)
	LimitNum int // 限制数值
	BestResults int // 最优成绩
	//标签相关
	Lable1 *JLable `orm:"rel(one)"`
	LableTitle1 string `orm:"size(16)"`
	Lable2 *JLable `orm:"rel(one)"`
	LableTitle2 string `orm:"size(16)"`
	Lable3 *JLable `orm:"rel(one)"`
	LableTitle3 string `orm:"size(16)"`
}

type JLable struct {
	Id int
	Lid string `orm:"size(64)"`
	Title string `orm:"size(64)"`
}

func init() {
	orm.RegisterModel(new(JNote),new(JLable))
}