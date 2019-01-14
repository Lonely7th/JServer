package colly

import (
	"ApiJServer/models"
	"ApiJServer/util"
	"bytes"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/esimov/stackblur-go"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type JNoteFactory struct {
	Id        int
	NoteId    string `orm:"size(64)"` // 唯一标识
	CreatTime int64  // 创建时间
	ResPath   string `orm:"size(256)"` // 图片地址
	Content   string `orm:"size(128)"` // 主要内容

	Sign     string `orm:"size(128)"` // Md5签名
	Released bool   // 是否已经发布
}

func init() {
	orm.RegisterModel(new(JNoteFactory))
}

//添加JNote到Factory
func AddJNote2Factory(resPath string, sign string, content string) bool {
	note := new(JNoteFactory)

	currentTime := int64(time.Now().UnixNano() / 1e6)
	note.CreatTime = currentTime

	noteId := util.GetRandomString(8) + util.GetCurrentTime()
	note.NoteId = noteId

	note.Content = content
	note.ResPath = resPath
	note.Sign = sign

	note.Released = false

	//note := JNoteFactory{NoteId: noteId,CreatTime: currentTime,ResPath: resPath,Content: content,Sign: sign,Released : false}

	o := orm.NewOrm()
	if created, id, err := o.ReadOrCreate(note, "res_path", "sign"); err == nil {
		if created {
			fmt.Println("New Factory:", id)
		} else {
			fmt.Println("Get Factory:", id)
		}
		return true
	}
	return false
}

//获取JNoteFactory列表
func ReleaseJNoteByFactory(num int) {
	o := orm.NewOrm()
	noteList := new([]JNoteFactory)

	rand.Seed(time.Now().UnixNano())
	start := rand.Intn(10000)
	_, err := o.QueryTable("j_note_factory").Filter("released", false).Limit(num, start).RelatedSel().All(noteList)
	if err == nil {
		for _, item := range *noteList {
			fmt.Println(item)
			//随机发布人信息
			releaser := models.GetRandReleaser()
			if releaser != nil {
				//加载图片
				picRes, gsRes, err := LoadNetPic(releaser.UserNo, item.ResPath)
				if err == nil {
					//随机NoteType
					ntype := GetRandNoteType()
					//随机限制长度
					limit := GetRandLimitNum()
					//随机裁剪格式
					format := GetRandFormat()
					//随机标签
					label1, label2, label3, labelTitle1, labelTitle2, labelTitle3 := GetRandLabels(item.Content)
					models.AddJNote("我发布了一条新动态，快来点击看看吧~", releaser.UserNo, picRes, gsRes, ntype, limit,
						false, format, label1, label2, label3, labelTitle1, labelTitle2, labelTitle3)
					//将工厂设置成已发布
					//item.Released = true
					//o.Update(item)
				} else {
					fmt.Println(err)
				}
			}
		}
	} else {
		fmt.Println(err)
	}
}

//获取随机标签
func GetRandLabels(content string) (int, int, int, string, string, string) {
	var label1, label2, label3 int
	var labelTitle1, labelTitle2, labelTitle3 string
	o := orm.NewOrm()

	rs := []rune(content)
	rl := len(rs)
	for i := 0; i < rl; i++ {
		for j := rl; j > i; j-- {
			if j-i > 1 {
				key := string(rs[i:j])

				label := new(models.JLabel)
				err := o.QueryTable("j_label").Filter("title", key).One(label)
				if err == nil && label != nil {
					fmt.Println("匹配到", label)
					if labelTitle1 == "" {
						label1 = label.Id
						labelTitle1 = label.Title
					} else if labelTitle2 == "" {
						label2 = label.Id
						labelTitle2 = label.Title
					} else if labelTitle3 == "" {
						label3 = label.Id
						labelTitle3 = label.Title
					}
					break
				}
			}
		}
	}
	if labelTitle1 != "" {
		return label1, label2, label3, labelTitle1, labelTitle2, labelTitle3
	} else {
		return 451, 0, 0, "其它", "", ""
	}
}

//保存图片到本地
func LoadNetPic(releaser string, imagPath string) (r string, g string, e error) {
	filePath := releaser + util.GetCurrentTime() + ".jpg"        // 原图路径
	gaussianPath := releaser + util.GetCurrentTime() + "-gs.jpg" // 模糊图路径

	resp, _ := http.Get(imagPath)
	body, _ := ioutil.ReadAll(resp.Body)

	//保存当前图片
	outRes, _ := os.Create(util.PicDir + filePath)
	_, err := io.Copy(outRes, bytes.NewReader(body))
	if err == nil {
		//保存高斯模糊后的图片
		src, _ := util.LoadImage(util.PicDir + filePath)

		var done = make(chan struct{}, 25)
		err1 := util.SaveImage(util.PicDir+gaussianPath, stackblur.Process(src, 20, done))

		if err1 != nil {
			fmt.Println(err1)
		}
	}

	return filePath, gaussianPath, err
}

//获取随机裁剪格式
func GetRandFormat() string {
	a := rand.Intn(10)
	var format string
	if a < 1 { //10%
		format = "4-3"
	} else if a < 3 { //20%
		format = "4-4"
	} else if a < 8 { //50%
		format = "6-4"
	} else { //20%
		format = "6-6"
	}
	return format
}

//获取随机NoteType
func GetRandNoteType() int {
	a := rand.Intn(2)
	var ntype int
	switch a {
	case 0:
		ntype = 1
		break
	case 1:
		ntype = 2
		break
	}
	return ntype
}

//获取随机限制长度
func GetRandLimitNum() int {
	a := rand.Intn(6)
	var num int
	switch a {
	case 0:
		num = 90
		break
	case 1:
		num = 120
		break
	case 2:
		num = 180
		break
	case 3:
		num = 240
		break
	case 4:
		num = 360
		break
	case 5:
		num = 500
		break
	}
	return num
}
