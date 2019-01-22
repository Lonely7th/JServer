package models

import (
	"ApiJServer/util"
	"bytes"
	"fmt"
	"github.com/Luxurioust/excelize"
	"github.com/astaxie/beego/orm"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
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
	IsFollow  bool // 是否已经关注
}

type UserFollow struct {
	Id        int
	UserNo    string `orm:"size(64)"` // 关注人
	Follower  *User  `orm:"rel(fk)"`  // 被关注人
	CreatTime int64  // 关注时间
}

func init() {
	orm.RegisterModel(new(User), new(UserFollow))
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
func getUserById(userNo string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("user_no", userNo).One(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return user
}

//获取用户信息
func GetUserById(userNo string, followNo string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable("user").Filter("user_no", userNo).One(user)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//判断当前用户是否已经关注
	userFollow := new(UserFollow)
	err2 := o.QueryTable("user_follow").Filter("user_no", followNo).Filter("follower", userNo).One(userFollow)
	if err2 == nil {
		user.IsFollow = true
	} else {
		user.IsFollow = false
	}
	return user
}

//修改用户信息
func UpdateUserInfo(uid string, ctype int, content string) bool {
	o := orm.NewOrm()
	user := getUserById(uid)
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

//生成随机发布人
func CreatRandReleaser(num int) {
	xlsx, err := excelize.OpenFile("./conf/user.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var basePhoneNumber = 4004160000
	for i := 1; i < num; i++ {
		user := GetUser(strconv.Itoa(basePhoneNumber))
		if user != nil {
			UpdateToken(user)
		} else {
			userNo := util.GetRandomString(24)
			//加载图片
			filePath := userNo + util.GetCurrentTime() + ".jpg" //图片路径
			fmt.Println(xlsx.GetCellValue("Sheet1", "B"+strconv.Itoa(i)))
			resp, _ := http.Get(xlsx.GetCellValue("Sheet1", "B"+strconv.Itoa(i)))
			body, _ := ioutil.ReadAll(resp.Body)
			outRes, _ := os.Create(util.PicDir + filePath)
			_, _ = io.Copy(outRes, bytes.NewReader(body))

			user := new(User)
			user.UserNo = userNo
			user.UserName = xlsx.GetCellValue("Sheet1", "A"+strconv.Itoa(i))
			user.UserToken = util.GetRandomString(16)
			user.UserPhone = strconv.Itoa(basePhoneNumber)
			user.NameHead = filePath
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
	rand.Seed(time.Now().UnixNano())
	o := orm.NewOrm()
	user := new(User)

	start := rand.Intn(1000)
	fmt.Println(start)
	err := o.QueryTable("user").Filter("user_phone__contains", "400416").Limit(1, start).RelatedSel().One(user)
	fmt.Println(user)
	if err == nil {
		return user
	} else {
		return nil
	}
}

//添加关注
func AddFollower(userNo string, follow string) (bool, *UserFollow) {
	userFollow := new(UserFollow)
	userFollow.UserNo = userNo
	userFollow.CreatTime = int64(time.Now().UnixNano() / 1e6)
	userFollow.Follower = getUserById(follow)

	o := orm.NewOrm()
	_, err := o.Insert(userFollow)
	if err == nil {
		return true, userFollow
	} else {
		return false, userFollow
	}
}

//取消关注
func DeleteFollower(userNo string, follow string) bool {
	userFollow := new(UserFollow)
	o := orm.NewOrm()
	err1 := o.QueryTable("user_follow").Filter("user_no", userNo).Filter("follower_id", follow).One(userFollow)

	if err1 == nil {
		_, err2 := o.Delete(userFollow)
		if err2 == nil {
			return true
		}
	}
	return false
}

//获取关注列表
func GetFollowList(userNo string) *[]UserFollow {
	o := orm.NewOrm()
	followList := new([]UserFollow)
	_, err := o.QueryTable("user_follow").Filter("user_no", userNo).OrderBy("-creat_time").RelatedSel().All(followList)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return followList
}
