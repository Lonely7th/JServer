package models

import (
	"ApiJServer/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type JNote struct {
	NoteId      string `orm:"pk"`       // 唯一标识
	Content     string `orm:"size(64)"` // 内容简介
	Releaser    *User  `orm:"rel(fk)"`  // OneToOne relation 发布者
	CreatTime   int64  // 发布时间
	ResPath     string `orm:"size(64)"` // 图片地址
	DisplayNum  int    // 展示次数
	CompleteNum int    // 完成次数
	JType       int    // 限制类型(0.无限制 1.时间限制 2.次数限制)
	LimitNum    int    // 限制数值
	BestResults int    // 最优成绩
	//标签相关
	Label1      int
	LabelTitle1 string `orm:"size(16)"`
	Label2      int
	LabelTitle2 string `orm:"size(16)"`
	Label3      int
	LabelTitle3 string `orm:"size(16)"`

	HideUser   bool   // 匿名发布
	CropFormat string `orm:"size(16)"`
}

type JLable struct {
	Id    int
	Lid   string `orm:"size(64)"`
	Title string `orm:"size(64)"`
}

type RCateNote struct {
	Id  int
	Cid string `orm:"size(64)"`
	Nid string `orm:"size(64)"`
}

func init() {
	orm.RegisterModel(new(JNote), new(JLable), new(RCateNote))
}

//添加JNote
func AddJNote(content string, releaser string, resPath string, jtype int, limitNum int, hideUser bool, cropFormat string,
	label1 int, label2 int, label3 int, labelTitle1 string, labelTitle2 string, labelTitle3 string) (bool, *JNote) {
	note := new(JNote)
	note.BestResults = -1
	note.CompleteNum = 0
	note.Content = content
	note.CreatTime = int64(time.Now().UnixNano() / 1e6)
	note.CropFormat = cropFormat
	note.DisplayNum = 0
	note.HideUser = hideUser
	note.JType = jtype
	note.LimitNum = limitNum

	NoteId := releaser + util.GetCurrentTime()
	note.NoteId = NoteId
	note.Releaser = GetUserById(releaser)
	note.ResPath = resPath
	note.Label1 = label1
	note.Label2 = label2
	note.Label3 = label3
	note.LabelTitle1 = labelTitle1
	note.LabelTitle2 = labelTitle2
	note.LabelTitle3 = labelTitle3

	o := orm.NewOrm()
	_, err := o.Insert(note)
	if err == nil {
		//添加分类关联
		relation1 := new(RCateNote)
		relation1.Cid = "0" // 新的
		relation1.Nid = NoteId
		_, _ = o.Insert(relation1)

		relation2 := new(RCateNote)
		relation2.Nid = NoteId
		switch jtype {
		case 0:
			relation2.Cid = "4" // 无限制
			break
		case 1:
			relation2.Cid = "5" // 时间限制
			break
		case 2:
			relation2.Cid = "6" // 次数限制
			break
		}
		_, _ = o.Insert(relation2)
		return true, note
	} else {
		fmt.Println(err)
		return false, note
	}
}

//获取JNote列表
func GetJNoteList(categroy string) *[]JNote {
	o := orm.NewOrm()
	noteList := new([]JNote)
	if categroy == "0" {
		_, err := o.QueryTable("j_note").RelatedSel().All(noteList)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	return noteList
}

//获取JNote详情
func GetJNoteDetails(noteId string) *JNote {
	o := orm.NewOrm()
	note := new(JNote)
	err := o.QueryTable("j_note").RelatedSel().Filter("note_id", noteId).One(note)
	if err != nil {
		return nil
	}
	return note
}

//提交结果
func PostJNoteResult() bool {
	return true
}

//收藏JNote
func StarJNote(userNo string, NoteId string) bool {
	return true
}
