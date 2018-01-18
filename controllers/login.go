package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/sparon/goshop/models"
	"github.com/sparon/goshop/services"
	"github.com/sparon/goshop/utils"
	"github.com/sparon/goshop/utils/urlpath"
)

//LoginController 登入控制器
type LoginController struct {
	LcfglyController
}

//GetLogin 登入页面
func (c *LoginController) GetLogin() {
	login := services.IsLogin(c.Ctx)
	if login != 0 { //用户已经登入,跳转到来源页面
		c.Redirect(urlpath.GotoURL(c.Ctx), 302)
		return
	}
	c.Data["title"] = "用户登入"
	c.Data["username"] = c.Ctx.GetCookie("ssoUsername") //获取COOKIE中的用户名
	c.TplName = "login.html"
}

//PostLogin 登入验证
func (c *LoginController) PostLogin() {
	vcode := c.GetString("vcode")
	if vcode == "" {
		c.JSONFail("请填写验证码")
		return
	}
	username := c.GetString("username")
	if username == "" {
		c.JSONFail("请填写用户名")
		return
	}
	password := c.GetString("password")
	if password == "" {
		c.JSONFail("请填写密码")
		return
	}
	remember := c.GetString("is_remember") == "on"
	if remember { //记住用户名
		c.Ctx.SetCookie("ssoUsername", username, utils.Conf.CookieExpireTime, utils.Conf.CookieDomain)
	}
	if strings.ToLower(vcode) != strings.ToLower(c.GetSession("loginvcode").(string)) { //验证验证码是否正确
		c.JSONFail("验证码错误")
		beego.Debug("正确验证码:", strings.ToLower(c.GetSession("loginvcode").(string)))
		return
	}
	memberID, err := models.Login.LoginUser(username, password)
	if err != nil {
		c.JSONFailError(err)
		return
	}
	conn := utils.RedisClient.Get() //获取redis客户端
	defer func() {
		conn.Close()
	}()
	id, _ := utils.IdWorker.NextString()
	_, err = conn.Do("SET", utils.RedisLoginKey(id), memberID, "EX", utils.Conf.CookieExpireTime)
	if err != nil {
		beego.Debug("Redis数据存储失败", err)
		enbytes, _ := utils.AesEncrypt([]byte(string("D" + id)))
		c.Ctx.SetCookie(utils.Conf.CookieLoginID, string(enbytes), utils.Conf.CookieExpireTime, "/", utils.Conf.CookieDomain)
		enbytes, _ = utils.AesEncrypt([]byte(strconv.FormatUint(uint64(memberID), 10)), []byte(id))
		c.Ctx.SetCookie(utils.Conf.CookieLogin, string(enbytes), utils.Conf.CookieExpireTime, "/", utils.Conf.CookieDomain)
	} else {
		enbytes, _ := utils.AesEncrypt([]byte(id))
		c.Ctx.SetCookie(utils.Conf.CookieLoginID, string(enbytes), utils.Conf.CookieExpireTime, "/", utils.Conf.CookieDomain)
	}
	//数据保存完毕,返回成功信息
	c.JSONSuccess(memberID, "用户登入成功", urlpath.GotoURL(c.Ctx))
}

//LoginOut 退出登入
func (c *LoginController) LoginOut() {
	id := c.Ctx.GetCookie(utils.Conf.CookieLoginID)
	if id != "" { //删除redis中的数据
		id, err := utils.AesDecrypt([]byte(id))
		if err == nil && !strings.HasPrefix(string(id), "D") { //数据存放在redis中
			conn := utils.RedisClient.Get()
			defer func() {
				conn.Close()
			}()
			conn.Do("DEL", utils.RedisLoginKey(string(id)))
		}
	}
	c.Ctx.SetCookie(utils.Conf.CookieLoginID, "", -1, "/", utils.Conf.CookieDomain)
	c.Ctx.SetCookie(utils.Conf.CookieLogin, "", -1, "/", utils.Conf.CookieDomain)
	isJSON, err := c.GetBool("json", false)
	if err == nil && isJSON {
		c.JSONSuccess(nil, "退出成功")
	} else {
		c.Redirect(urlpath.GotoURL(c.Ctx), 302) //跳转到来源页面
	}
}

//GetRegister 注册页面
func (c *LoginController) GetRegister() {
	c.Data["title"] = "用户注册"
	c.TplName = "register.html"
}

//PostRegister 注册信息
func (c *LoginController) PostRegister() {
	username := c.GetString("username")        //用户名
	password := c.GetString("password")        //密码
	confPassword := c.GetString("psw_confirm") //二次密码
	vcode := c.GetString("vcodesms")           //验证码
	refer := c.GetString("refer_id")           //推广员ID
	source := c.GetString("source", "pc")      //来源平台
	ip := c.Ctx.Input.IP()                     //IP地址
	referURL := c.Ctx.Input.Referer()          //来源地址
	if vcode == "" {
		c.JSONFail("请填写验证码")
		return
	}
	if username == "" || password == "" {
		c.JSONFail("用户名密码不能为空")
		return
	}
	if password != confPassword {
		c.JSONFail("两次输入密码不一致")
		return
	}
	if strings.ToLower(vcode) != strings.ToLower(c.GetSession("smscode").(string)) {
		c.JSONFail("验证码错误")
		return
	}
	if !models.Login.CheckUser(username) {
		c.JSONFail("帐号已存在")
		return
	}
	memberID, exp := services.Register(username, password, refer, referURL, source, ip)
	if exp != nil || memberID == 0 {
		c.JSONFailException(exp)
		return
	}
	c.JSONSuccess(memberID, "会员注册成功")
}

//PostSms 发送短信验证码
func (c *LoginController) PostSms() {
	phone := c.GetString("phone")
	loginvcode := c.GetString("vcode")
	if phone == "" || !utils.ValidateString(phone, utils.Phone) {
		c.JSONFail("请填写手机号")
		return
	}
	if c.GetSession("loginvcode") != strings.ToLower(loginvcode) {
		c.JSONFail("验证码错误")
		return
	}
	if c.GetSession("lastendtime").(int64)-time.Now().Unix() < 120 {
		c.JSONFail("发送请求过于频繁,请稍后再试")
		return
	}
	if services.SendSmsVcode(phone) {
		c.SetSession("lastendtime", time.Now().Unix()) //记录发送时间
		c.JSONSuccess(nil, "验证码发送成功请注意查收")
		return
	}
	c.JSONFail("验证码发送失败")
}

//PostEmail 密码找回发送邮件
func (c *LoginController) PostEmail() {
	email := strings.TrimSpace(c.GetString("email"))
	key := strings.TrimSpace(c.GetString("key"))
	c.Data["account"] = strings.TrimSpace(c.GetString("account"))
	if key == "" || c.Data["account"].(string) == "" {
		c.JSONFail("缺少必要数据,请刷新重试")
		return
	}
	if email == "" {
		c.JSONFail("邮箱地址不能为空")
		return
	}
	data := make(map[string]interface{})
	data["createtime"] = time.Now().Unix()
	data["key"] = key
	data["account"] = c.Data["account"]
	str, err := json.Marshal(data)
	if err != nil {
		c.JSONFail("数据处理错误")
		return
	}
	str, err = utils.AesEncrypt(str)
	if err != nil {
		c.JSONFail("数据处理错误")
		return
	}
	c.Data["url"] = "http://" + c.Ctx.Request.Host + urlpath.Path("lostPwd") + "/step3?key=" + string(str)
	c.Data["ip"] = c.Ctx.Input.IP()
	c.TplName = "lostPwd/email.html"
	str, err = c.RenderBytes()
	if err != nil {
		c.JSONFail("数据处理错误")
		return
	}
	if services.SendEmail(email, "HW123商城密码找回", string(str)) {
		c.JSONSuccess(nil, "邮件发送成功,请注意查收")
	} else {
		c.JSONFail("邮件发送失败")
	}
}

//GetLostPwd 密码找回页面
func (c *LoginController) GetLostPwd() {
	c.Data["title"] = "密码找回"
	c.Data["step"] = 1
	c.TplName = "lostPwd/lostPwd.html,lostPwd/step1.html"
}

//PostLostPwdStep2 找回密码第二步
func (c *LoginController) PostLostPwdStep2() {
	account := c.GetString("account")
	memberLoginInfo := services.MemberLoginAccount(account)
	if len(memberLoginInfo) < 1 {
		c.Data["msg"] = "由于您并未验证手机或者邮箱，无法自助找回密码，请联系网站客服！"
		c.TplName = "lostPwd/error.html"
		return
	}
	c.Data["account"] = account
	memberID := uint64(0)
	notMsg := true
	for _, v := range memberLoginInfo {
		if memberID == 0 {
			memberID = uint64(v["member_id"].(int))
		}
		if strings.ToLower(v["login_type"].(string)) != "local" {
			c.Data[strings.ToLower(v["login_type"].(string))] = v["login_account"]
			notMsg = false
		}
	}
	if notMsg { //未绑定任何邮箱和手机,返回错误
		c.Data["msg"] = "由于您并未验证手机或者邮箱，无法自助找回密码，请联系网站客服！"
		c.TplName = "lostPwd/error.html"
		return
	}
	secret, err := services.ChangePwdCreateKey(memberID)
	if err != nil { //密钥创建失败,返回错误
		c.Data["msg"] = "由于您并未验证手机或者邮箱，无法自助找回密码，请联系网站客服！"
		c.TplName = "lostPwd/error.html"
		return
	}
	c.Data["secret"] = secret
	c.TplName = "lostPwd/step2.html"
}

//PostLostPwdStep3 密码找回第三步[POST]
func (c *LoginController) PostLostPwdStep3() {
	c.Data["key"] = c.GetString("key")
	c.Data["account"] = c.GetString("account")
	vcode := c.GetString("vcode")
	//两个小时验证码过期
	if time.Now().Unix()-c.GetSession("lastendtime").(int64) > 2*60*60 {
		c.JSONFail("验证码已过期")
		return
	}
	if vcode != strings.ToLower(c.GetSession("smscode").(string)) {
		c.JSONFail("验证码错误")
		return
	}
	c.TplName = "lostPwd/step3.html"
	str, err := c.RenderString()
	if err != nil {
		c.JSONFail("数据处理异常")
	} else {
		c.JSONSuccess(str, "密码修改成功")
	}
}

//GetLostPwdStep3 密码找回第三步[GET] 邮箱验证通过邮箱中的连接地址跳转过来
func (c *LoginController) GetLostPwdStep3() {
	key := strings.TrimSpace(c.GetString("key"))
	c.Data["title"] = "密码找回"
	if key == "" {
		c.Data["msg"] = "非法操作！"
		c.TplName = "lostPwd/error.html"
		return
	}
	data, err := utils.AesDecrypt([]byte(key))
	if err != nil {
		c.Data["msg"] = "非法操作！"
		c.TplName = "lostPwd/error.html"
		return
	}
	err = json.Unmarshal(data, &c.Data)
	if err != nil {
		c.Data["msg"] = "非法操作！"
		c.TplName = "lostPwd/error.html"
		return
	}
	it, err := strconv.ParseInt(c.Data["createtime"].(string), 10, 64)
	if err != nil {
		c.Data["msg"] = "非法操作！"
		c.TplName = "lostPwd/error.html"
		return
	}
	if time.Now().Unix()-it > 2*60*60 {
		c.Data["msg"] = "本链接已过期"
		c.TplName = "lostPwd/error.html"
		return
	}
	c.TplName = "lostPwd/step3.html"
}

//PostLostPwdStep4 密码找回第四步
func (c *LoginController) PostLostPwdStep4() {
	key := strings.TrimSpace(c.GetString("key"))
	account := c.GetString("account")
	pwd := c.GetString("pwd")
	pwdC := c.GetString("pwd_c")
	if pwd != pwdC {
		c.JSONFail("两次密码不一致")
		return
	}
	pamMember, err := services.LoginPam(account)
	if err != nil {
		c.JSONFail("用户不存在")
		return
	}
	pwdlog, err := services.MemberPwdLog(key, pamMember.MemberID)
	if err != nil {
		c.JSONFail(err.Error())
		return
	}
	if services.UpdatePwd(pwd, pamMember.MemberID) {
		services.UpdateMemberPwdLogUsed(pwdlog["pwdlog_id"].(int))
		c.TplName = "lostPwd/step4.html"
		str, _ := c.RenderString()
		c.JSONSuccess(str, "密码修改成功")
	} else {
		c.JSONFail("密码修改失败")
	}
}
