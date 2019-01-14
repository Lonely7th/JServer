package main

import (
	"ApiJServer/colly"
	"ApiJServer/models"
	_ "ApiJServer/routers"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

const StrMd5Sign = "9b063dfaef3f9deaf4413ffb8f26d247"

func init() {
	//初始化数据库
	InitDataBase()
	//初始化拦截器
	InitSignFilter()
	//初始化日志
	InitLoger()
	//初始化标签系统(需要更新标签列表时开启)
	//models.InitLabel()
	//初始化定时器
	colly.InitTimer()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger" //http://127.0.0.1:8080/swagger/
	}

	beego.Run()

}

// 初始化数据库
func InitDataBase() {
	dbhost := beego.AppConfig.String("mysqlurls")
	dbport := beego.AppConfig.String("mysqlport")
	dbuser := beego.AppConfig.String("mysqluser")
	dbpassword := beego.AppConfig.String("mysqlpass")
	db := beego.AppConfig.String("mysqldb")

	//注册mysql Driver
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	//构造conn连接
	conn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + db + "?charset=utf8"
	//注册数据库连接
	orm.RegisterDataBase("default", "mysql", conn)

	_ = orm.RunSyncdb("default", false, true)
	orm.Debug = true
}

// 初始化控制台日志
func InitLoger() {
	// 创建一个日志记录器，参数为缓冲区的大小
	log := logs.NewLogger(10000)
	// 设置日志记录方式：控制台记录
	log.SetLogger("console", "")
	// 设置日志写入缓冲区的等级：Debug级别（最低级别，所以所有log都会输入到缓冲区）
	log.SetLevel(logs.LevelDebug)
	// 输出log时能显示输出文件名和行号（非必须）
	log.EnableFuncCallDepth(true)
}

// 初始化拦截器
func InitSignFilter() {
	beego.InsertFilter("/*", beego.BeforeRouter, FilterSign)
}

var FilterSign = func(ctx *context.Context) {
	timeStamp := ctx.Input.Query("timeStamp")
	token := ctx.Input.Query("token")
	sign := ctx.Input.Query("sign")

	//判断签名是否正确
	data := []byte(timeStamp + StrMd5Sign + token)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	if md5str1 != sign {
		data, _ := json.Marshal(models.GetErrorResult("503", "签名错误"))
		ctx.Output.Body([]byte(data))
		return
	}

	//判断时间戳是否正确
	curTime := time.Now()
	reqTime, _ := strconv.ParseInt(timeStamp, 10, 64)
	if curTime.Sub(time.Unix(reqTime/1e3, 0)).Seconds() > 20 {
		data, _ := json.Marshal(models.GetErrorResult("408", "请求超时"))
		ctx.Output.Body([]byte(data))
		return
	}

	//判断token是否正确
}
