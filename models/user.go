package models

import (
	"ApiJServer/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id int
	UserNo string `orm:"size(64)"`
	UserName string `orm:"size(64)"`
	UserToken string `orm:"size(64)"`
	UserPhone string `orm:"size(16)"`
	NameHead string `orm:"size(128)"`
	NameCity string `orm:"size(32)"`
	CreatTime int
}

func init() {
	orm.RegisterModel(new(User))
}

func AddUser(phone string) (result bool, u *User) {
	user := new(User)
	user.UserNo = util.GetRandomString(24)
	user.UserName = phone
	user.UserToken = util.GetRandomString(16)
	user.UserPhone = phone
	user.NameHead = ""
	user.NameCity = ""
	user.CreatTime = int(time.Now().Unix())

	o := orm.NewOrm()
	_, err := o.Insert(user)
	if err == nil {
		return true, user
	}else{
		return false, user
	}
}

func GetUser(phoneNumber string) (u *User) {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("user_phone", phoneNumber).One(user)
	//err := o.QueryTable("user").Filter("id",key).All(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

func GetUserById(Id string) (u *User) {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("user_no", Id).One(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

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
	}else{
		return false
	}
}

func Login(phoneNumber string) (u *User) {
	user := GetUser(phoneNumber)
	if user != nil {
		return UpdateToken(user)
	}else{
		result, newUser := AddUser(phoneNumber)
		if result{
			return newUser
		}
		return nil
	}
}

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
