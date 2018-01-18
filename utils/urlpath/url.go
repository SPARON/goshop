package urlpath

import (
	"net/url"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/context"
)

//URL 路由地址信息
var URL config.Configer

func init() {
	var err error
	URL, err = config.NewConfig("ini", "conf/url.conf")
	if err != nil {
		beego.Error("配置文件加载失败")
		os.Exit(1)
	}
}

//Path 根据名称获取URL路由地址
func Path(name string) string {
	if URL != nil {
		return URL.String(name)
	}
	beego.Warning("没有配置[" + name + "]的URL信息")
	return ""
}

//Host 根据名称获取域名信息
func Host(name string) string {
	return URL.String("hosts." + name)
}

//PathAll 获取所有地址
func PathAll() map[string]interface{} {
	result, err := URL.GetSection("default")
	if err != nil {
		beego.Error("获取所有地址错误", err)
	}
	res := make(map[string]interface{})
	for k, v := range result { //循环处理'.'分级
		splitPath(k, v, res)
	}
	return res
}

//splitPath 分割路径信息
func splitPath(key, value string, res map[string]interface{}) {
	ks := strings.SplitN(key, ".", 2)
	if len(ks) > 1 {
		if _, ok := res[ks[0]]; !ok { //如果KV存在
			res[ks[0]] = make(map[string]interface{})
		}
		splitPath(ks[1], value, res[ks[0]].(map[string]interface{}))
	} else {
		res[ks[0]] = value
	}
}

//GotoURL 获取跳转地址
func GotoURL(c *context.Context) string {
	u := c.Input.Query("redirectURL")
	if u == "" { //获取到页面来源地址
		u = c.Request.Referer()
	}
	if u != "" { //如果地址存在
		l, err := url.Parse(u)
		if err == nil && l.Host != c.Request.URL.Host && l.Path != c.Request.URL.Path {
			//地址是正确的URL地址,并且不是当前请求地址..
			return u
		}
	}
	return Host("main")
}
