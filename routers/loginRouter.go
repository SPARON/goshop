package routers

import (
	"github.com/astaxie/beego"
	"github.com/sparon/goshop/controllers"
	"github.com/sparon/goshop/utils/urlpath"
)

func init() {
	loginCtl := &controllers.LoginController{}
	//登入页面
	beego.Router(urlpath.Path("login"), loginCtl, "get:GetLogin")
	//登入用户密码验证
	beego.Router(urlpath.Path("login"), loginCtl, "post:PostLogin")
	//用户注册
	beego.Router(urlpath.Path("register"), loginCtl, "get:GetRegister")
	//发送短信验证码
	beego.Router(urlpath.Path("sms"), loginCtl, "post:PostSms")
	//密码找回
	beego.Router(urlpath.Path("lostpwd"), loginCtl, "get:GetLostPwd")
	//密码找回第二步
	beego.Router(urlpath.Path("lostpwd")+"/step2", loginCtl, "post:PostLostPwdStep2")
	//密码找回第三步[POST]
	beego.Router(urlpath.Path("lostpwd")+"/step3", loginCtl, "post:PostLostPwdStep3")
	//密码找回第三步[GET]
	beego.Router(urlpath.Path("lostpwd")+"/step3", loginCtl, "get:GetLostPwdStep3")
	//密码找回邮件发送
	beego.Router(urlpath.Path("lostpwd")+"/email", loginCtl, "post:PostEmail")
	//密码找回第四步
	beego.Router(urlpath.Path("lostpwd")+"/step4", loginCtl, "post:PostLostPwdStep4")
	//退出登入
	beego.Router(urlpath.Path("loginout"), loginCtl, "*:LoginOut")
}
