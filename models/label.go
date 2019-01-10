package models

import (
	"fmt"
	"github.com/Luxurioust/excelize"
	"github.com/astaxie/beego/orm"
	"os"
)

type JLabel struct {
	Id    int
	Title string `orm:"size(64)"`
}

func init() {
	orm.RegisterModel(new(JLabel))
}

//获取标签列表
func GetLabelList() *[]JLabel {
	o := orm.NewOrm()
	labelList := new([]JLabel)
	_, err := o.QueryTable("j_label").RelatedSel().All(labelList)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return labelList
}

//初始化标签列表
func InitLabel() {
	fmt.Println("Init Label...")
	//读取数据
	xlsx, err := excelize.OpenFile("./conf/label.xlsx")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get all the rows in a sheet.
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		for _, colCell := range row {
			//插入数据
			fmt.Println(colCell, "\t")
			o := orm.NewOrm()
			label := JLabel{Title: colCell}
			// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
			if created, id, err := o.ReadOrCreate(&label, "title"); err == nil {
				if created {
					fmt.Println("New Label:", id)
				} else {
					fmt.Println("Get Label:", id)
				}
			}
		}
	}
}
