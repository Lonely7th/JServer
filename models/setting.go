package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type VersionCode struct {
	Id            int
	Code          string `orm:"size(64)"`
	UpdateContent string `orm:"size(64)"`
	ApkPath       string `orm:"size(64)"`
	Type          int    // 1.Android
}

type FeedBack struct {
	Id        int
	UserNo    string `orm:"size(64)"`
	Content   string `orm:"size(64)"`
	CreatTime int64  // 反馈时间
}

func init() {
	orm.RegisterModel(new(VersionCode), new(FeedBack))
}

//添加反馈
func AddFeedBack(userNo string, content string) error {
	back := new(FeedBack)
	back.UserNo = userNo
	back.Content = content
	back.CreatTime = int64(time.Now().UnixNano() / 1e6)

	o := orm.NewOrm()
	_, err := o.Insert(back)
	return err
}

//获取当前版本号
func GetVersionCode() *VersionCode {
	o := orm.NewOrm()
	version := new(VersionCode)
	err := o.QueryTable("version_code").Filter("type", 1).One(version)
	if err != nil {
		return nil
	}
	return version
}
