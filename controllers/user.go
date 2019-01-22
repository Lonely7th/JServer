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
	picPath := util.PicDir + filePath // 图片保存的路径

	if err != nil {
		u.Data["json"] = models.GetErrorResult("404", "加载文件失败")
	} else {
		defer src.Close()
		_ = u.SaveToFile("res", picPath)

		if models.UpdateUserInfo(userNo, 2, filePath) {
			u.Data["json"] = models.GetJsonResult(util.ImagePath + filePath)
		} else {
			u.Data["json"] = models.GetErrorResult("403", "保存数据失败")
		}

	}
	u.ServeJSON()
}

// @Title GetUserInfo
// @Description 查询用户信息
// @Param	user_no		query 	string	true		"用户Id"
// @Param	follow_no		query 	string	true		"关注者Id"
// @Success 200 {string} success
// @router /getUserInfo [get]
func (u *UserController) GetUserInfo() {
	userNo := u.GetString("user_no")
	followNo := u.GetString("follow_no")
	user := models.GetUserById(userNo, followNo)
	user.NameHead = util.ImagePath + user.NameHead
	if user != nil {
		u.Data["json"] = models.GetJsonResult(user)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title AddFollower
// @Description 添加关注
// @Param	type		query 	int	true		"操作类型 0.添加关注 1.删除关注"
// @Param	user_no		query 	string	true		"关注人Id"
// @Param	follow_id		query 	string	true		"被关注人Id"
// @Success 200 {string} success
// @router /addFollower [post]
func (u *UserController) AddFollower() {
	starType, _ := u.GetInt("type")
	userNo := u.GetString("user_no")
	followId := u.GetString("follow_id")
	if starType == 0 { //添加关注
		result, userFollow := models.AddFollower(userNo, followId)
		if result {
			u.Data["json"] = models.GetJsonResult(userFollow)
		} else {
			u.Data["json"] = models.GetErrorResult("403", "失败")
		}
	} else if starType == 1 { //删除关注
		if models.DeleteFollower(userNo, followId) {
			u.Data["json"] = models.GetJsonResult("")
		} else {
			u.Data["json"] = models.GetErrorResult("403", "失败")
		}
	}
	u.ServeJSON()
}

// @Title GetFollowList
// @Description 获取用户关注列表
// @Param	user_no		query 	string	true		"用户Id"
// @Success 200 {string} success
// @router /getFollowList [get]
func (u *UserController) GetFollowList() {
	user_id := u.GetString("user_no")
	list := models.GetFollowList(user_id)
	if list != nil {
		u.Data["json"] = models.GetJsonResult(list)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}
