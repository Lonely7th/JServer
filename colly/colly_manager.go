package colly

import (
	"ApiJServer/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type JNoteFactory struct {
	NoteId    string `orm:"pk"` // 唯一标识
	CreatTime int64  // 创建时间
	ResPath   string `orm:"size(64)"`  // 图片地址
	GsResPath string `orm:"size(64)"`  // 高斯模糊图片地址
	Content   string `orm:"size(128)"` // 主要内容

	Sign     string `orm:"size(128)"` // Md5签名
	Released bool   // 是否已经发布
}

func init() {
	orm.RegisterModel(new(JNoteFactory))
}

//添加JNote到Factory
func AddJNote2Factory(resPath string, gsResPath string, content string) (bool, *JNoteFactory) {
	note := new(JNoteFactory)
	CurrentTime := int64(time.Now().UnixNano() / 1e6)
	note.CreatTime = CurrentTime

	NoteId := util.GetCurrentTime() + util.GetRandomString(8)
	note.NoteId = NoteId
	note.Content = content
	note.ResPath = resPath
	note.GsResPath = gsResPath

	note.Released = false

	o := orm.NewOrm()
	_, err := o.Insert(note)
	if err == nil {
		return true, note
	} else {
		fmt.Println(err)
		return false, nil
	}
}

//获取JNoteFactory列表
func GetJNoteByFactory() {
	o := orm.NewOrm()
	noteList := new([]JNoteFactory)

	_, err := o.QueryTable("j_note_factory").Filter("released", false).OrderBy("-creat_time").RelatedSel().All(noteList)
	if err == nil {
		for _, item := range *noteList {
			fmt.Println(item.NoteId)
			// 开始发布
		}
	} else {
		fmt.Println(err)
	}
}
