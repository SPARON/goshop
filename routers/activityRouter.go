package routers

import (
	"github.com/astaxie/beego"
	"github.com/sparon/goshop/controllers"
	"github.com/sparon/goshop/utils/urlpath"
)

func init() {
	ActivityCtl := &controllers.ActivityController{}
	//预约页面
	beego.Router(urlpath.Path("yuyue"), ActivityCtl, "get:Yuyue")
	//抢购
	beego.Router(urlpath.Path("qiangou"), ActivityCtl, "post:Qiangou")
}
