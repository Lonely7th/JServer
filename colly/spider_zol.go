package colly

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

func GetJNoteByZol() {
	c := colly.NewCollector()

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299"

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("%+v\r\n%+v\r\n", *r, *(r.Headers))
	})

	c.OnHTML("td>table.list-div-table>tbody>tr", func(e *colly.HTMLElement) {
		fmt.Println(e)
		//hyShy := ZhjhHyShyl{
		//
		//	Hydm: e.ChildText("td:first-child"),
		//
		//	Hymc: e.ChildText("td:nth-child(2)"),
		//}
		//
		//zxsj, err := strconv.ParseFloat(e.ChildText("td:nth-child(3)"), 64)
		//
		//if err == nil {
		//
		//	hyShy.Zxsj = &zxsj
		//
		//}
		//
		//gpjs, err := strconv.ParseInt(e.ChildText("td:nth-child(4)"), 10, 32)
		//
		//if err == nil {
		//
		//	hyShy.Gpjs = int(gpjs)
		//
		//}
		//
		//ksjs, err := strconv.ParseInt(e.ChildText("td:nth-child(5)"), 10, 32)
		//
		//if err == nil {
		//
		//	hyShy.Ksjs = int(ksjs)
		//
		//}
		//
		//jygy, err := strconv.ParseFloat(e.ChildText("td:nth-child(6)"), 64)
		//
		//if err == nil {
		//
		//	hyShy.Jygy = &jygy
		//
		//}
		//
		//jsgy, err := strconv.ParseFloat(e.ChildText("td:nth-child(7)"), 64)
		//
		//if err == nil {
		//
		//	hyShy.Jsgy = &jsgy
		//
		//}
	})

	c.OnScraped(func(_ *colly.Response) {
		bData, _ := json.MarshalIndent("", "", "\t")
		fmt.Println(string(bData))
	})

	var err = c.Visit("http://sj.zol.com.cn/bizhi/p2/")
	if err != nil {
		log.Fatal(err)
	}
}
