package colly

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/gocolly/colly"
	"log"
	"strconv"
)

type UserHead struct {
	Id   int
	Path string `orm:"size(256)"`
}

func init() {
	orm.RegisterModel(new(UserHead))
}

const WYGXBaseUrl = "https://www.woyaogexing.com/touxiang"

func GetUserInfoByWYGX() {
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.Headers)
	})

	//获取数据
	c.OnHTML("div.pMain > div.txList", func(element *colly.HTMLElement) {
		links := element.ChildAttrs("a > img", "src")
		fmt.Println("link = ", links)

		//保存数据
		for _, item := range links {
			o := orm.NewOrm()
			head := UserHead{Path: item[2:]}
			// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
			if created, id, err := o.ReadOrCreate(&head, "path"); err == nil {
				if created {
					fmt.Println("New Head:", id)
				} else {
					fmt.Println("Get Head:", id)
				}
			}
		}
	})

	for i := 2; i < 100; i++ {
		url := WYGXBaseUrl + "/index_" + strconv.Itoa(i) + ".html"
		fmt.Println(url)
		var err = c.Visit(url)
		if err != nil {
			log.Fatal(err)
		}
	}
}
