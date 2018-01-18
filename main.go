package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/sparon/goshop/controllers"
	_ "github.com/sparon/goshop/routers"
)

func main() {
	beego.SetStaticPath("/image", "static/image")
	beego.SetStaticPath("/images", "static/images")
	beego.SetStaticPath("/statics", "static/statics")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/fonts", "static/fonts")
	beego.SetStaticPath("/js", "static/js")
	beego.ErrorController(&controllers.LcfglyController{})
	if beego.BConfig.RunMode == "proc" {
		beego.SetLevel(beego.LevelError) //设置日志等级
	} else {
		//beego.SetLevel(beego.LevelError)
	}
	beego.Run()
}
