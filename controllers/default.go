package controllers

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"gitee.com/sparon/goshop/services"
	"gitee.com/sparon/goshop/utils"
	"gitee.com/sparon/goshop/utils/urlpath"
)

//MainController 主控制器
type MainController struct {
	LcfglyController
}

//Index 应用首页,跳转到主域名
func (c *MainController) Index() {
	c.Redirect(urlpath.Host("main"), 301) //跳转到主域名
}

//OAuth 用户登入检测
func (c *MainController) OAuth() {
	login := services.IsLogin(c.Ctx) //获取用户登入信息
	c.JSONSuccess(login, "请求成功")
}

//VCode 验证码
func (c *MainController) VCode() {
	resp, err := http.Get(urlpath.Host("api") + "/vcode.htm")
	if err != nil {
		beego.Debug("HTTP PROXY代理异常", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	defer func() {
		resp.Body.Close()
	}()
	cookies := resp.Cookies()
	for _, cc := range cookies {
		if cc.Name == "vcode" {
			beego.Debug("生成验证码:", cc.Value)
			c.SetSession("loginvcode", cc.Value)
		}
	}
	bf := bytes.NewBufferString("")
	bf.ReadFrom(resp.Body)
	c.Ctx.Output.Body(bf.Bytes())
	c.EnableRender = false
}

//Crypt 加密解密测试
func (c *MainController) Crypt() {
	c.EnableRender = false
	response := c.Ctx.ResponseWriter.ResponseWriter
	src := c.GetString("src")
	if src == "" {
		response.Write([]byte("请输入参数src:要加密的字符串"))
		return
	}
	key := c.GetString("key", string(utils.Conf.SecretKey))
	res, err := utils.AesEncrypt([]byte(src), []byte(key))
	if err != nil {
		response.Write([]byte("加密失败:" + err.Error()))
		return
	}
	response.Write([]byte("加密字符串:" + src + "\r\n"))
	response.Write([]byte("加密结果:"))
	response.Write(res)
}

//ProxyHandler 代理服务
func (c *MainController) ProxyHandler() {
	c.EnableRender = false
	r := c.Ctx.Request
	w := c.Ctx.ResponseWriter.ResponseWriter
	regx, err := regexp.Compile("\\.api$")
	if err == nil && regx.MatchString(r.URL.Path) {
		remote, err := url.Parse(urlpath.Host("api"))
		if err != nil {
			c.JSONFail("服务器错误")
			return
		}
		bf := bytes.NewBufferString(remote.RawQuery)
		if remote.RawQuery != "" {
			bf.WriteString("&")
		}
		bf.WriteString("member_id=")
		loginID := services.IsLogin(c.Ctx)
		remote.RawQuery = strconv.FormatUint(loginID, 10)
		loginmember, _ := utils.AesEncrypt([]byte(remote.RawQuery))
		bf.Write(loginmember)
		remote.RawQuery = bf.String()
		r.URL.Path = regx.ReplaceAllString(r.URL.Path, ".htm")
		beego.Debug("代理请求地址:", r.URL.Path)
		r.RequestURI = r.URL.Path

		r.Body = &proxyBody{ //重写body请求数据
			buffer: bytes.NewBufferString(r.Form.Encode()),
		}
		r.ContentLength = int64(len(r.Form.Encode()))
		proxy := httputil.NewSingleHostReverseProxy(remote)
		r.Host = remote.Host
		proxy.ServeHTTP(w, r)
	} else {
		beego.Debug("Proxy Pass Fail :", r.URL.Path)
		c.Error404()
	}
}

//Test 测试页面
func (c *MainController) Test() {
	c.EnableRender = false
}

//proxyBody 获取代理对象body数据
//因为beego在解析路由时就把body中的内容已经解析出来,body为空
//这会导致post请求时出现请求头中包含数据长度而实际请求内容为空的错误
type proxyBody struct {
	buffer *bytes.Buffer
}

func (proxy *proxyBody) Read(p []byte) (n int, err error) {
	return proxy.buffer.Read(p)
}

func (proxy *proxyBody) Close() error {
	return nil
}
