package controllers

import (
	"ApiJServer/models"
	"ApiJServer/util"
	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

var PicDir = "D:/go_server_src/"
var ImagePath = "http://172.31.71.35:8089/"

// @Title Login
// @Description Logs user into the system
// @Param	phoneNumber		query 	string	true		"The phoneNumber for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [post]
func (u *UserController) Login() {
	phoneNumber := u.GetString("phoneNumber")

	user := models.Login(phoneNumber)
	u.Data["json"] = models.GetJsonResult(user)
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [post]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

// @Title ChangeInfo
// @Description 修改用户信息
// @Param	user_no		query 	string	true		"用户Id"
// @Param	content		query 	string	true		"修改后的内容"
// @Param	ctype		query 	int	true		"待修改的类型(0：用户昵称 1：所在城市)"
// @Success 200 {string} success
// @router /changeInfo [post]
func (u *UserController) ChangeInfo() {
	ctype, _ := u.GetInt("ctype")
	userNo := u.GetString("user_no")
	content := u.GetString("content")
	if models.UpdateUserInfo(userNo, ctype, content) {
		u.Data["json"] = models.GetJsonResult("")
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title ChangeHead
// @Description 修改用户头像
// @Param	user_no		query 	string	true		"用户Id"
// @Param	res		query 	string	true		"图片内容base64"
// @Success 200 {string} success
// @router /changeHead [post]
func (u *UserController) ChangeHead() {
	src, _, err := u.GetFile("res")
	userNo := u.GetString("user_no")
	filePath := userNo + util.GetCurrentTime() + ".jpg"
	picPath := PicDir + filePath // 图片保存的路径

	if err != nil {
		u.Data["json"] = models.GetErrorResult("404", "加载文件失败")
	} else {
		defer src.Close()
		_ = u.SaveToFile("res", picPath)

		if models.UpdateUserInfo(userNo, 2, filePath) {
			u.Data["json"] = models.GetJsonResult(ImagePath + filePath)
		} else {
			u.Data["json"] = models.GetErrorResult("403", "保存数据失败")
		}

	}
	u.ServeJSON()
}

// @Title GetUserInfo
// @Description 查询用户信息
// @Param	user_no		query 	string	true		"用户Id"
// @Success 200 {string} success
// @router /getUserInfo [get]
func (u *UserController) GetUserInfo() {
	userNo := u.GetString("user_no")
	user := models.GetUserById(userNo)
	user.NameHead = ImagePath + user.NameHead
	if user != nil {
		u.Data["json"] = models.GetJsonResult(user)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}
