package colly

import (
	"time"
)

//启动定时器
func InitTimer() {
	go StartReleaseTimer()
}

//启动分发器定时器
func StartReleaseTimer() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			ReleaseJNoteByFactory(1)
		}
	}
}

//启动数据爬虫定时器
func StartSpiderTimer() {

}

//启动分析系统定时器
func StartAnalysisTimer() {

}

//启动日志排查定时器
func StartLogerTimer() {

}
