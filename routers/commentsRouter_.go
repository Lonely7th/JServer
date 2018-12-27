package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["ApiJServer/controllers:CategroyController"] = append(beego.GlobalControllerRouter["ApiJServer/controllers:CategroyController"],
        beego.ControllerComments{
            Method: "GetCategroy",
            Router: `/getCategroy`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ApiJServer/controllers:UserController"] = append(beego.GlobalControllerRouter["ApiJServer/controllers:UserController"],
        beego.ControllerComments{
            Method: "ChangeHead",
            Router: `/changeHead`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ApiJServer/controllers:UserController"] = append(beego.GlobalControllerRouter["ApiJServer/controllers:UserController"],
        beego.ControllerComments{
            Method: "ChangeInfo",
            Router: `/changeInfo`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ApiJServer/controllers:UserController"] = append(beego.GlobalControllerRouter["ApiJServer/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetUserInfo",
            Router: `/getUserInfo`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ApiJServer/controllers:UserController"] = append(beego.GlobalControllerRouter["ApiJServer/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["ApiJServer/controllers:UserController"] = append(beego.GlobalControllerRouter["ApiJServer/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
