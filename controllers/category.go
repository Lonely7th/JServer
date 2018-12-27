package controllers

import (
	"ApiJServer/models"
	"github.com/astaxie/beego"
)

type CategroyController struct {
	beego.Controller
}

// @Title GetCategroy
// @Description 获取分类列表
// @Success 200 {string} success
// @router /getCategroy [get]
func (c *CategroyController) GetCategroy() {
	result := models.GetJCategroyList()
	if result != nil {
		c.Data["json"] = models.GetJsonResult(result)
	}else{
		c.Data["json"] = models.GetErrorResult("403","失败")
	}
	c.ServeJSON()
}