package colly

import (
	"ApiJServer/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/gocolly/colly"
	"log"
	"time"
)

type JNoteFactory struct {
	NoteId    string `orm:"pk"` // 唯一标识
	CreatTime int64  // 创建时间
	ResPath   string `orm:"size(64)"` // 图片地址
	GsResPath string `orm:"size(64)"` // 高斯模糊图片地址
	//标签相关
	Label1      int
	LabelTitle1 string `orm:"size(16)"`
	Label2      int
	LabelTitle2 string `orm:"size(16)"`
	Label3      int
	LabelTitle3 string `orm:"size(16)"`

	Released bool // 是否已经发布
}

func init() {
	orm.RegisterModel(new(JNoteFactory))
}

func GetJNoteByZol() {
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("%+v\r\n%+v\r\n", *r, *(r.Headers))
	})

	c.OnScraped(func(_ *colly.Response) {
		bData, _ := json.MarshalIndent("", "", "\t")
		fmt.Println(string(bData))
	})

	var err = c.Visit("http://sj.zol.com.cn/bizhi/p2/")
	if err != nil {
		log.Fatal(err)
	}
}

//添加JNote到Factory
func AddJNote2Factory(releaser string, resPath string, gsResPath string, label1 int, label2 int, label3 int,
	labelTitle1 string, labelTitle2 string, labelTitle3 string) (bool, *JNoteFactory) {
	note := new(JNoteFactory)
	CurrentTime := int64(time.Now().UnixNano() / 1e6)
	note.CreatTime = CurrentTime

	NoteId := releaser + util.GetCurrentTime()
	note.NoteId = NoteId
	note.ResPath = resPath
	note.GsResPath = gsResPath
	note.Label1 = label1
	note.Label2 = label2
	note.Label3 = label3
	note.LabelTitle1 = labelTitle1
	note.LabelTitle2 = labelTitle2
	note.LabelTitle3 = labelTitle3

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
