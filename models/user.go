package models

import (
	"ApiJServer/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	"math/rand"
	"strconv"
	"time"
)

type User struct {
	UserNo    string `orm:"pk"`
	UserName  string `orm:"size(64)"`
	UserToken string `orm:"size(64)"`
	UserPhone string `orm:"size(16)"`
	NameHead  string `orm:"size(128)"`
	NameCity  string `orm:"size(32)"`
	CreatTime int64
}

func init() {
	orm.RegisterModel(new(User))
}

//添加新用户
func AddUser(phone string) (result bool, u *User) {
	user := new(User)
	user.UserNo = util.GetRandomString(24)
	user.UserName = phone
	user.UserToken = util.GetRandomString(16)
	user.UserPhone = phone
	user.NameHead = ""
	user.NameCity = ""
	user.CreatTime = int64(time.Now().UnixNano() / 1e6)

	o := orm.NewOrm()
	_, err := o.Insert(user)
	if err == nil {
		return true, user
	} else {
		return false, user
	}
}

//获取用户信息
func GetUser(phoneNumber string) (u *User) {
	if util.Validate(phoneNumber) {
		return nil
	}
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("user_phone", phoneNumber).One(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

//获取用户信息
func GetUserById(Id string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("user_no", Id).One(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

//修改用户信息
func UpdateUserInfo(uid string, ctype int, content string) bool {
	o := orm.NewOrm()
	user := GetUserById(uid)
	if user != nil {
		switch ctype {
		case 0:
			user.UserName = content
			break
		case 1:
			user.NameCity = content
			break
		case 2:
			user.NameHead = content
			break
		}
		if _, err := o.Update(user); err == nil {
			return true
		}
		return false
	} else {
		return false
	}
}

//登录
func Login(phoneNumber string) (u *User) {
	user := GetUser(phoneNumber)
	if user != nil {
		return UpdateToken(user)
	} else {
		result, newUser := AddUser(phoneNumber)
		if result {
			return newUser
		}
		return nil
	}
}

//刷新token
func UpdateToken(u *User) (result *User) {
	o := orm.NewOrm()
	u.UserToken = util.GetRandomString(16)
	if num, err := o.Update(u); err == nil {
		fmt.Println(num)
	}
	return u
}

func DeleteUser(uid string) {

}

func InitUser() {

}

//生成随机发布人
func CreatRandReleaser(num int) {
	var basePhoneNumber = 4004160000
	for i := 0; i < num; i++ {
		user := GetUser(strconv.Itoa(basePhoneNumber))
		if user != nil {
			UpdateToken(user)
		} else {
			user := new(User)
			user.UserNo = util.GetRandomString(24)
			user.UserName = ""
			user.UserToken = util.GetRandomString(16)
			user.UserPhone = strconv.Itoa(basePhoneNumber)
			user.NameHead = ""
			user.NameCity = ""
			user.CreatTime = int64(time.Now().UnixNano() / 1e6)

			o := orm.NewOrm()
			_, _ = o.Insert(user)
		}
		basePhoneNumber++
	}
}

//获取随机发布人
func GetRandReleaser() *User {
	o := orm.NewOrm()
	user := new(User)

	start := rand.Intn(1000)
	err := o.QueryTable("j_note_factory").Filter("user_phone__contains", "400416").Limit(1, start).RelatedSel().One(user)
	if err == nil {
		return user
	} else {
		return nil
	}
}
