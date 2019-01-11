package colly

import (
	"crypto/md5"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

const ZolBaseUrl = "http://sj.zol.com.cn"

func GetJNoteByZol() {
	c := colly.NewCollector()
	detailLink := c.Clone()
	detailLink2 := c.Clone()
	detailLink3 := c.Clone()
	detailController := c.Clone()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.Headers)
	})

	//获取标签列表
	c.OnHTML("dl.filter-item.first.clearfix > dd > a", func(element *colly.HTMLElement) {
		label := element.Attr("href")

		for i := 1; i < 30; i++ {
			_ = detailLink.Visit(ZolBaseUrl + label + strconv.Itoa(i) + ".html")
		}
	})

	//获取当前页的一级页面列表
	detailLink.OnHTML("ul.pic-list2 > li", func(element *colly.HTMLElement) {
		link := element.ChildAttr("a", "href")
		//fmt.Println("link = " + BaseUrl + link)
		_ = detailLink2.Visit(ZolBaseUrl + link)
	})

	//获取当前页的二级页面列表
	detailLink2.OnHTML("div.wrapper.photo-set.mt15", func(element *colly.HTMLElement) {
		links := element.ChildAttrs("ul > li > a", "href")
		for _, item := range links {
			//fmt.Println("link2 = " + BaseUrl + item)
			_ = detailLink3.Visit(ZolBaseUrl + item)
		}
	})

	//获取当前页的三级页面列表
	detailLink3.OnHTML("body", func(element *colly.HTMLElement) {
		title := element.DOM.Find("div.wrapper.photo-tit.clearfix > h1 > a").Text()
		link, _ := element.DOM.Find("dl.model.wallpaper-down.clearfix > dd > a").Attr("href")
		//fmt.Println("title = " + title)
		//fmt.Println("link = " + BaseUrl + link)
		_ = detailController.Visit(ZolBaseUrl + link + "?" + title)
	})

	//解析页面内容
	detailController.OnHTML("body", func(element *colly.HTMLElement) {
		title := strings.Split(element.Request.URL.String(), "?")[1]
		link := element.ChildAttr("img", "src")
		//fmt.Println("title = ", title)
		//fmt.Println("link = " + link)
		//生成签名
		data := []byte(title + link)
		has := md5.Sum(data)
		md5str1 := fmt.Sprintf("%x", has)
		//插入数据库
		AddJNote2Factory(link, md5str1, title)
	})

	c.OnScraped(func(_ *colly.Response) {

	})

	var err = c.Visit(ZolBaseUrl + "/bizhi/p2/")
	if err != nil {
		log.Fatal(err)
	}
}
