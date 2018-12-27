package main

import (
	_ "ApiJServer/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init(){
	InitDataBase()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"//http://127.0.0.1:8080/swagger/
	}

	beego.SetLogger("file", `{"filename":"logs/test.log"}`)
	beego.Run()

}

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
