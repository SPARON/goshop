package routers

import (
	"github.com/astaxie/beego"
	"gitee.com/sparon/goshop/controllers"
	"gitee.com/sparon/goshop/utils/urlpath"
)

func init() {
	mainCtl := &controllers.MainController{}
	//首页
	beego.Router("/", mainCtl, "get:Index")
	//用户登入检测接口
	beego.Router(urlpath.Path("oauth"), mainCtl, "*:OAuth")
	//获取图片验证码
	beego.Router(urlpath.Path("vcodeimg"), mainCtl, "get:VCode")
	//代理服务器
	beego.Router("/*.api", mainCtl, "*:ProxyHandler")
	//加解密测试
	beego.Router("/crypt", mainCtl, "get:Crypt")
	//测试
	beego.Router("/test", mainCtl, "*:Test")
}
