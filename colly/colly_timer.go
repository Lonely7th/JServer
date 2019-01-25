package colly

import (
	"ApiJServer/controllers"
	"time"
)

//启动定时器
func InitTimer() {
	go StartReleaseTimer()
	go StartAnalysisTimer()
}

//启动分发器定时器
func StartReleaseTimer() {
	tickerSpider := time.NewTicker(100 * time.Second)

	for {
		select {
		case <-tickerSpider.C:
			ReleaseJNoteByFactory(1)
			break
		}
	}
}

//启动数据爬虫定时器
func StartSpiderTimer() {

}

//启动分析系统定时器
func StartAnalysisTimer() {
	tickerAnalysis := time.NewTicker(600 * time.Second)

	for {
		select {
		case <-tickerAnalysis.C:
			controllers.NoteSuccessAnalysisSystem()
			break
		}
	}
}

//启动日志排查定时器
func StartLogerTimer() {

}
