package controllers

import (
	"ApiJServer/models"
	"ApiJServer/util"
	"github.com/astaxie/beego"
)

type NoteController struct {
	beego.Controller
}

// @Title AddJNote
// @Description 添加JNote
// @Param	res		query 	file	true		"图片资源"
// @Param	jtype		query 	int	true		"限制类型"
// @Param	label1		query 	string	true		"标签1"
// @Param	label2		query 	string	true		"标签2"
// @Param	label3		query 	string	true		"标签3"
// @Param	limitNum		query 	int	true		"限制数值"
// @Param	content		query 	string	true		"内容简介"
// @Param	releaser		query 	string	true		"发布者Id"
// @Param	hideUser		query 	bool	true		"匿名发布"
// @Param	cropFormat		query 	string	true		"裁剪格式"
// @Param	labelTitle1		query 	string	true		"标签1标题"
// @Param	labelTitle2		query 	string	true		"标签2标题"
// @Param	labelTitle3		query 	string	true		"标签3标题"
// @Success 200 {string} success
// @router /addJNote [post]
func (u *NoteController) AddJNote() {
	content := u.GetString("content")
	releaser := u.GetString("releaser")
	resData, _, _ := u.GetFile("res")
	jtype, _ := u.GetInt("jtype")
	limitNum, _ := u.GetInt("limitNum")
	hideUser, _ := u.GetBool("hideUser")
	cropFormat := u.GetString("cropFormat")
	label1, _ := u.GetInt("label1")
	labelTitle1 := u.GetString("labelTitle1")
	label2, _ := u.GetInt("label2")
	labelTitle2 := u.GetString("labelTitle2")
	label3, _ := u.GetInt("label3")
	labelTitle3 := u.GetString("labelTitle3")

	filePath := releaser + util.GetCurrentTime() + ".jpg"
	result, note := models.AddJNote(content, releaser, filePath, jtype, limitNum, hideUser, cropFormat, label1, label2, label3, labelTitle1, labelTitle2, labelTitle3)
	if result == true {
		//保存图片
		picPath := PicDir + filePath // 图片保存的路径
		defer resData.Close()
		_ = u.SaveToFile("res", picPath)

		note.ResPath = ImagePath + filePath
		u.Data["json"] = models.GetJsonResult(note)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "保存失败")
	}
	u.ServeJSON()
}

// @Title GetJNoteList
// @Description 获取JNote列表
// @Param	categroy		query 	string	true		"用户Id"
// @Success 200 {string} success
// @router /getJNoteList [get]
func (u *NoteController) GetJNoteList() {
	category := u.GetString("categroy")
	page, _ := u.GetInt("page")
	list := models.GetJNoteList(category, page)
	for _, item := range *list {
		item.ResPath = ImagePath + item.ResPath
		item.Releaser.UserToken = ""
		item.Releaser.NameHead = ImagePath + item.Releaser.NameHead
	}
	if list != nil {
		u.Data["json"] = models.GetJsonResult(list)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title GetJNoteDetails
// @Description 获取JNote详情
// @Param	note_id		query 	string	true		"JNote Id"
// @Success 200 {string} success
// @router /getJNoteDetails [get]
func (u *NoteController) GetJNoteDetails() {
	noteId := u.GetString("note_id")
	userId := u.GetString("user_id")
	note := models.GetJNoteDetails(noteId, userId)
	if note != nil {
		note.ResPath = ImagePath + note.ResPath
		note.Releaser.UserToken = ""
		note.Releaser.NameHead = ImagePath + note.Releaser.NameHead
		u.Data["json"] = models.GetJsonResult(note)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title PostJNoteResult
// @Description 提交结果
// @Param	user_no		query 	string	true		"用户Id"
// @Param	note_id		query 	string	true		"NoteId"
// @Param	status		query 	int	true		"完成状态"
// @Param	score		query 	int	true		"完成成绩"
// @Success 200 {string} success
// @router /postJNoteResult [post]
func (u *NoteController) PostJNoteResult() {
	userNo := u.GetString("user_no")
	noteId := u.GetString("note_id")
	status, _ := u.GetInt("status")
	score, _ := u.GetInt("score")
	result, noteScore := models.PostJNoteResult(userNo, noteId, status, score)
	if result {
		u.Data["json"] = models.GetJsonResult(noteScore)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title StarJNote
// @Description 收藏JNote
// @Param	type		query 	int	true		"操作类型 0.添加收藏 1.删除收藏"
// @Param	user_no		query 	string	true		"用户Id"
// @Param	note_id		query 	string	true		"NoteId"
// @Param	star_id		query 	string	true		"收藏Id"
// @Success 200 {string} success
// @router /starJNote [post]
func (u *NoteController) StarJNote() {
	starType, _ := u.GetInt("type")
	userNo := u.GetString("user_no")
	noteId := u.GetString("note_id")
	if starType == 0 { //添加收藏
		result, starJNote := models.AddStarJNote(userNo, noteId)
		if result {
			u.Data["json"] = models.GetJsonResult(starJNote)
		} else {
			u.Data["json"] = models.GetErrorResult("403", "失败")
		}
	} else if starType == 1 { //删除收藏
		if models.DeleteStarJNote(userNo, noteId) {
			u.Data["json"] = models.GetJsonResult("")
		} else {
			u.Data["json"] = models.GetErrorResult("403", "失败")
		}
	}
	u.ServeJSON()
}

// @Title GetStarNoteList
// @Description 获取用户收藏列表
// @Param	user_no		query 	string	true		"用户Id"
// @Success 200 {string} success
// @router /getStarNoteList [get]
func (u *NoteController) GetStarNoteList() {
	user_id := u.GetString("user_no")
	list := models.GetStarNoteList(user_id)
	if list != nil {
		u.Data["json"] = models.GetJsonResult(list)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title GetLabelList
// @Description 获取标签列表
// @Success 200 {string} success
// @router /getLabelList [get]
func (u *NoteController) GetLabelList() {
	list := models.GetLabelList()
	if list != nil {
		u.Data["json"] = models.GetJsonResult(list)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}

// @Title GetReleaseList
// @Description 获取用户发布列表
// @Success 200 {string} success
// @router /getReleaseList [get]
func (u *NoteController) GetReleaseList() {
	user_id := u.GetString("user_no")
	list := models.GetUserReleaseNoteList(user_id)
	if list != nil {
		u.Data["json"] = models.GetJsonResult(list)
	} else {
		u.Data["json"] = models.GetErrorResult("403", "失败")
	}
	u.ServeJSON()
}
