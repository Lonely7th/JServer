package models

import (
	"ApiJServer/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type JNote struct {
	NoteId      string `orm:"pk"`       // 唯一标识
	Content     string `orm:"size(64)"` // 内容简介
	Releaser    *User  `orm:"rel(fk)"`  // relation 发布者
	CreatTime   int64  // 发布时间
	ResPath     string `orm:"size(64)"` // 图片地址
	GsResPath   string `orm:"size(64)"` // 高斯模糊图片地址
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

	FavoriteId  int
	SuccessRate int // 成功率
}

type RCategoryNote struct {
	Id        int
	Cid       string `orm:"size(64)"`
	Note      *JNote `orm:"rel(fk)"` //关联的JNote
	CreatTime int64  // 发布时间
}

type JNoteStar struct {
	Id        int
	UserNo    string `orm:"size(64)"` // 用户Id
	Note      *JNote `orm:"rel(fk)"`  // NoteId
	CreatTime int64  // 收藏时间
}

type JNoteScore struct {
	Id     int
	UserNo string `orm:"size(64)"` // 用户Id
	NoteId string `orm:"size(64)"` // NoteId
	Status int    // 完成状态 0.失败 1.成功 2.刷新最好成绩
	Score  int    // 成绩
}

func init() {
	orm.RegisterModel(new(JNote), new(RCategoryNote), new(JNoteStar), new(JNoteScore))
}

//添加JNote
func AddJNote(content string, releaser string, resPath string, gsResPath string, jtype int, limitNum int, hideUser bool, cropFormat string,
	label1 int, label2 int, label3 int, labelTitle1 string, labelTitle2 string, labelTitle3 string) (bool, *JNote) {
	note := new(JNote)
	note.BestResults = limitNum
	note.CompleteNum = 0
	note.Content = content
	CurrentTime := int64(time.Now().UnixNano() / 1e6)
	note.CreatTime = CurrentTime
	note.CropFormat = cropFormat
	note.DisplayNum = 0
	note.HideUser = hideUser
	note.JType = jtype
	note.LimitNum = limitNum

	NoteId := releaser + util.GetCurrentTime()
	note.NoteId = NoteId
	note.Releaser = getUserById(releaser)
	note.ResPath = resPath
	note.GsResPath = gsResPath
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
		if jtype != 0 && limitNum < 120 { // 低于120次的默认为困难
			relation1 := new(RCategoryNote)
			relation1.Cid = "2" // 困难
			relation1.CreatTime = CurrentTime
			relation1.Note = note
			_, _ = o.Insert(relation1)
		}

		relation2 := new(RCategoryNote)
		relation2.Note = note
		relation2.CreatTime = CurrentTime
		switch jtype {
		case 0:
			relation2.Cid = "3" // 无限制
			break
		case 1:
			relation2.Cid = "4" // 时间限制
			break
		case 2:
			relation2.Cid = "5" // 次数限制
			break
		}
		_, _ = o.Insert(relation2)
		return true, note
	} else {
		fmt.Println(err)
		return false, note
	}
}

const PageSize int = 10

//获取JNote列表
func GetJNoteList(categroy string, phoneSign string, page int) *[]JNote {
	o := orm.NewOrm()
	noteList := new([]JNote)
	categroyList := new([]RCategoryNote)
	switch categroy {
	case "0":
		_, err := o.QueryTable("j_note").OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(noteList)
		if err != nil {
			fmt.Println(err)
		}
		break
	case "1": // 推荐
		redisManager, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			fmt.Println("Connect to redis error", err)
			return noteList
		}
		defer redisManager.Close()

		var lastTime int64
		keyExit, err := redis.Bool(redisManager.Do("EXISTS", phoneSign))
		if err == nil && keyExit {
			value, err := redis.String(redisManager.Do("GET", phoneSign))
			if err == nil {
				lastTime, _ = strconv.ParseInt(value, 10, 64)
			}
		} else {
			lastTime = 0
		}
		_, err = o.QueryTable("j_note").Filter("creat_time__gte", lastTime).OrderBy("creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(noteList)
		if err != nil {
			fmt.Println(err)
		} else {
			for index, item := range *noteList {
				if index == len(*noteList)-1 {
					_, err = redisManager.Do("SET", phoneSign, item.CreatTime)
				}
			}
		}
		break
	case "2":
		_, err := o.QueryTable("r_category_note").Filter("cid", categroy).OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(categroyList)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, item := range *categroyList {
				*noteList = append(*noteList, *item.Note)
			}
		}
		break
	case "3":
		_, err := o.QueryTable("r_category_note").Filter("cid", categroy).OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(categroyList)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, item := range *categroyList {
				*noteList = append(*noteList, *item.Note)
			}
		}
		break
	case "4":
		_, err := o.QueryTable("r_category_note").Filter("cid", categroy).OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(categroyList)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, item := range *categroyList {
				*noteList = append(*noteList, *item.Note)
			}
		}
		break
	case "5":
		_, err := o.QueryTable("r_category_note").Filter("cid", categroy).OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(categroyList)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, item := range *categroyList {
				*noteList = append(*noteList, *item.Note)
			}
		}
		break
	case "6": // 加载全部
		_, err := o.QueryTable("j_note").OrderBy("-creat_time").RelatedSel().Limit(PageSize, (page-1)*PageSize).All(noteList)
		if err != nil {
			fmt.Println(err)
		}
		break
	}
	return noteList
}

//获取JNote详情
func GetJNoteDetails(noteId string, userId string) *JNote {
	o := orm.NewOrm()
	note := new(JNote)
	err1 := o.QueryTable("j_note").RelatedSel().Filter("note_id", noteId).One(note)
	if err1 != nil {
		fmt.Println(err1)
		return nil
	}

	//判断当前用户是否收藏此Note
	noteStar := new(JNoteStar)
	err2 := o.QueryTable("j_note_star").Filter("user_no", userId).Filter("note_id", noteId).One(noteStar)
	if err2 == nil {
		note.FavoriteId = 1
	}
	return note
}

//提交结果
func PostJNoteResult(userNo string, noteId string, status int, score int) (bool, *JNoteScore) {
	noteScore := new(JNoteScore)
	noteScore.NoteId = noteId
	noteScore.Score = score
	noteScore.Status = status
	noteScore.UserNo = userNo

	o := orm.NewOrm()
	_, err := o.Insert(noteScore)
	if status == 2 { // 刷新最佳成绩
		note := GetJNoteDetails(noteId, userNo)
		note.BestResults = score
		o.Update(note)
	}

	if err == nil {
		return true, noteScore
	} else {
		return false, noteScore
	}
}

//添加收藏
func AddStarJNote(userNo string, noteId string) (bool, *JNoteStar) {
	noteStar := new(JNoteStar)
	noteStar.UserNo = userNo
	noteStar.CreatTime = int64(time.Now().UnixNano() / 1e6)
	noteStar.Note = GetJNoteDetails(noteId, userNo)

	o := orm.NewOrm()
	_, err := o.Insert(noteStar)
	if err == nil {
		return true, noteStar
	} else {
		return false, noteStar
	}
}

//取消收藏
func DeleteStarJNote(userNo string, noteId string) bool {
	noteStar := new(JNoteStar)
	o := orm.NewOrm()
	err1 := o.QueryTable("j_note_star").Filter("user_no", userNo).Filter("note_id", noteId).One(noteStar)

	if err1 == nil {
		_, err2 := o.Delete(noteStar)
		if err2 == nil {
			return true
		}
	}
	return false
}

//获取收藏列表
func GetStarNoteList(userId string) *[]JNoteStar {
	o := orm.NewOrm()
	noteList := new([]JNoteStar)
	_, err := o.QueryTable("j_note_star").Filter("user_no", userId).OrderBy("-creat_time").RelatedSel().All(noteList)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return noteList
}

//获取用户发布列表
func GetUserReleaseNoteList(userId string) *[]JNote {
	o := orm.NewOrm()
	noteList := new([]JNote)
	_, err := o.QueryTable("j_note").Filter("releaser_id", userId).OrderBy("-creat_time").RelatedSel().All(noteList)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return noteList
}
