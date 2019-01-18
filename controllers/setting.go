package controllers

import (
	"ApiJServer/models"
	"github.com/astaxie/beego"
)

type SettingController struct {
	beego.Controller
}

// @Title AddFeedBack
// @Description 反馈与建议
// @Param	userNo		query 	string	true		"用户Id"
// @Param	content		query 	string	true		"反馈内容"
// @Success 200 {string} login success
// @router /addFeedBack [post]
func (u *SettingController) AddFeedBack() {
	userNo := u.GetString("userNo")
	content := u.GetString("content")

	err := models.AddFeedBack(userNo, content)
	if err == nil {
		u.Data["json"] = models.GetJsonResult("")
	} else {
		u.Data["json"] = models.GetErrorResult("403", "")
	}
	u.ServeJSON()
}

// @Title GetVersionCode
// @Description 获取当前版本号
// @Success 200 {string} success
// @router /getVersionCode [get]
func (u *SettingController) GetVersionCode() {
	ver := models.GetVersionCode()

	if ver != nil {
		u.Data["json"] = models.GetJsonResult(ver)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "")
	}
	u.ServeJSON()
}
