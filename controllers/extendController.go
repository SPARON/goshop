package controllers

import (
	"bytes"
	"io"
	"regexp"

	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/robertkrimen/otto"
	"github.com/sparon/goshop/services"
	"github.com/sparon/goshop/utils"
	"github.com/sparon/goshop/utils/urlpath"
)

//templateJS artTemplate模版引擎
var templateJS = `!function(){function a(a){return a.replace(t,"").replace(u,",").replace(v,"").replace(w,"").replace(x,"").split(y)}function b(a){return"'"+a.replace(/('|\\)/g,"\\$1").replace(/\r/g,"\\r").replace(/\n/g,"\\n")+"'"}function c(c,d){function e(a){return m+=a.split(/\n/).length-1,k&&(a=a.replace(/\s+/g," ").replace(/<!--[\w\W]*?-->/g,"")),a&&(a=s[1]+b(a)+s[2]+"\n"),a}function f(b){var c=m;if(j?b=j(b,d):g&&(b=b.replace(/\n/g,function(){return m++,"$line="+m+";"})),0===b.indexOf("=")){var e=l&&!/^=[=#]/.test(b);if(b=b.replace(/^=[=#]?|[\s;]*$/g,""),e){var f=b.replace(/\s*\([^\)]+\)/,"");n[f]||/^(include|print)$/.test(f)||(b="$escape("+b+")")}else b="$string("+b+")";b=s[1]+b+s[2]}return g&&(b="$line="+c+";"+b),r(a(b),function(a){if(a&&!p[a]){var b;b="print"===a?u:"include"===a?v:n[a]?"$utils."+a:o[a]?"$helpers."+a:"$data."+a,w+=a+"="+b+",",p[a]=!0}}),b+"\n"}var g=d.debug,h=d.openTag,i=d.closeTag,j=d.parser,k=d.compress,l=d.escape,m=1,p={$data:1,$filename:1,$utils:1,$helpers:1,$out:1,$line:1},q="".trim,s=q?["$out='';","$out+=",";","$out"]:["$out=[];","$out.push(",");","$out.join('')"],t=q?"$out+=text;return $out;":"$out.push(text);",u="function(){var text=''.concat.apply('',arguments);"+t+"}",v="function(filename,data){data=data||$data;var text=$utils.$include(filename,data,$filename);"+t+"}",w="'use strict';var $utils=this,$helpers=$utils.$helpers,"+(g?"$line=0,":""),x=s[0],y="return new String("+s[3]+");";r(c.split(h),function(a){a=a.split(i);var b=a[0],c=a[1];1===a.length?x+=e(b):(x+=f(b),c&&(x+=e(c)))});var z=w+x+y;g&&(z="try{"+z+"}catch(e){throw {filename:$filename,name:'Render Error',message:e.message,line:$line,source:"+b(c)+".split(/\\n/)[$line-1].replace(/^\\s+/,'')};}");try{var A=new Function("$data","$filename",z);return A.prototype=n,A}catch(B){throw B.temp="function anonymous($data,$filename) {"+z+"}",B}}var d=function(a,b){return"string"==typeof b?q(b,{filename:a}):g(a,b)};d.version="3.0.0",d.config=function(a,b){e[a]=b};var e=d.defaults={openTag:"<%",closeTag:"%>",escape:!0,cache:!0,compress:!1,parser:null},f=d.cache={};d.render=function(a,b){return q(a,b)};var g=d.renderFile=function(a,b){var c=d.get(a)||p({filename:a,name:"Render Error",message:"Template not found"});return b?c(b):c};d.get=function(a){var b;if(f[a])b=f[a];else if("object"==typeof document){var c=document.getElementById(a);if(c){var d=(c.value||c.innerHTML).replace(/^\s*|\s*$/g,"");b=q(d,{filename:a})}}return b};var h=function(a,b){return"string"!=typeof a&&(b=typeof a,"number"===b?a+="":a="function"===b?h(a.call(a)):""),a},i={"<":"&#60;",">":"&#62;",'"':"&#34;","'":"&#39;","&":"&#38;"},j=function(a){return i[a]},k=function(a){return h(a).replace(/&( ?![\w#]+;)|[<>"']/g,j)},l=Array.isArray||function(a){return"[object Array]"==={}.toString.call(a)},m=function(a,b){var c,d;if(l(a))for(c=0,d=a.length;d>c;c++)b.call(a,a[c],c,a);else for(c in a)b.call(a,a[c],c)},n=d.utils={$helpers:{},$include:g,$string:h,$escape:k,$each:m};d.helper=function(a,b){o[a]=b};var o=d.helpers=n.$helpers;d.onerror=function(a){var b="Template Error\n\n";for(var c in a)b+="<"+c+">\n"+a[c]+"\n\n";"object"==typeof console&&console.error(b)};var p=function(a){return d.onerror(a),function(){return"{Template Error}"}},q=d.compile=function(a,b){function d(c){try{return new i(c,h)+""}catch(d){return b.debug?p(d)():(b.debug=!0,q(a,b)(c))}}b=b||{};for(var g in e)void 0===b[g]&&(b[g]=e[g]);var h=b.filename;try{var i=c(a,b)}catch(j){return j.filename=h||"anonymous",j.name="Syntax Error",p(j)}return d.prototype=i.prototype,d.toString=function(){return i.toString()},h&&b.cache&&(f[h]=d),d},r=n.$each,s="break,case,catch,continue,debugger,default,delete,do,else,false,finally,for,function,if,in,instanceof,new,null,return,switch,this,throw,true,try,typeof,var,void,while,with,abstract,boolean,byte,char,class,const,double,enum,export,extends,final,float,goto,implements,import,int,interface,long,native,package,private,protected,public,short,static,super,synchronized,throws,transient,volatile,arguments,let,yield,undefined",t=/\/\*[\w\W]*?\*\/|\/\/[^\n]*\n|\/\/[^\n]*$|"(?:[^"\\]|\\[\w\W])*"|'(?:[^'\\]|\\[\w\W])*'|\s*\.\s*[$\w\.]+/g,u=/[^\w$]+/g,v=new RegExp(["\\b"+s.replace(/,/g,"\\b|\\b")+"\\b"].join("|"),"g"),w=/^\d[^,]*|,\d[^,]*/g,x=/^,+|,+$/g,y=/^$|,+/;e.openTag="{{",e.closeTag="}}";var z=function(a,b){var c=b.split(":"),d=c.shift(),e=c.join(":")||"";return e&&(e=", "+e),"$helpers."+d+"("+a+e+")"};e.parser=function(a){a=a.replace(/^\s/,"");var b=a.split(" "),c=b.shift(),e=b.join(" ");switch(c){case"if":a="if("+e+"){";break;case"else":b="if"===b.shift()?" if("+b.join(" ")+")":"",a="}else"+b+"{";break;case"/if":a="}";break;case"each":var f=b[0]||"$data",g=b[1]||"as",h=b[2]||"$value",i=b[3]||"$index",j=h+","+i;"as"!==g&&(f="[]"),a="$each("+f+",function("+j+"){";break;case"/each":a="});";break;case"echo":a="print("+e+");";break;case"print":case"include":a=c+"("+b.join(",")+");";break;default:if(/^\s*\|\s*[\w\$]/.test(e)){var k=!0;0===a.indexOf("#")&&(a=a.substr(1),k=!1);for(var l=0,m=a.split("|"),n=m.length,o=m[l++];n>l;l++)o=z(o,m[l]);a=(k?"=":"=#")+o}else a=d.helpers[c]?"=#"+c+"("+b.join(",")+");":"="+a}return a},"function"==typeof define?define(function(){return d}):"undefined"!=typeof exports?module.exports=d:this.template=d}();template.helper('forindex',function(index,format){value = "";for(i=0;i<index;i++){value+=format}return value;});var render = template.compile(source);render(data)`

//扩展beego控制器..

//LcfglyController 控制器扩展
type LcfglyController struct {
	beego.Controller
	OncePutOUT bool                   //一次性输出
	Pjax       bool                   //Pjax
	Data       map[string]interface{} //数据对象
	LoadTop    bool                   //加载顶部内容
}

// Init generates default values of controller operations.
func (c *LcfglyController) Init(ctx *context.Context, controllerName, actionName string, app interface{}) {
	c.Controller.Init(ctx, controllerName, actionName, app)
	c.Data = make(map[string]interface{})           //数据保存对象更换为map[string]interface{}
	c.Pjax = c.Ctx.Input.Header("x-pjax") == "true" //是否是PJAX加载
}

//Finish 结束处理
func (c *LcfglyController) Finish() {
	c.Controller.Data["json"] = c.Data
}

// Render sends the response with rendered template bytes as text/html type.
func (c *LcfglyController) Render() error {
	if !c.EnableRender {
		return nil
	}
	return c.renderTp()
}

// RenderBytes returns the bytes of rendered template string. Do not send out response.
func (c *LcfglyController) RenderBytes() ([]byte, error) { //使用JS模版引擎artTemplate解析模版
	c.Data["url"] = urlpath.PathAll()
	c.Data["gotoURL"] = urlpath.GotoURL(c.Ctx) //页面跳转地址
	c.Data["request"] = c.Ctx.Request.Form     //请求参数集合
	//content := bytes.NewBufferString("")
	//content.WriteString(templateJS)
	bs, err := c.Controller.RenderBytes()
	if err != nil {
		return nil, err
	}
	/*params, err := json.Marshal(c.Data)
	if err != nil {
		return nil, err
	}*/
	//beego.Debug("模版参数=======")
	//beego.Debug(string(params))
	//beego.Debug("End")
	//content.WriteString(`console.log(JSON.stringify(data));`)
	//content.WriteString(`template.helper('forindex',function(index,format){value = "";for(i=0;i<index;i++){value+=format}return value;});`)
	//content.WriteString(`var render = template.compile(source);render(data)`)
	//content.Write(params)
	//content.WriteString(")")
	JSEngine := otto.New()
	JSEngine.Set("source", string(bs))
	JSEngine.Set("data", c.Data)
	str, err := JSEngine.Run(templateJS)
	if err != nil {
		return nil, err
	}
	return []byte(str.String()), nil
}

//RenderTp 渲染模版,默认会先自动输出conf/app.conf中配置的defaultTemplate文件
func (c *LcfglyController) renderTp() error {
	//页面要加载顶部信息,但是如果是pjax加载,那么必须要加载整体布局.informations时才进行头部加载
	if c.LoadTop && (!c.Pjax || c.Ctx.Input.Header("X-PJAX-Container") == ".informations") {
		c.TplName = "public/top.html," + c.TplName //加载顶部信息
	}
	templateName := strings.Split(c.TplName, ",")
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8") //设置header
	//c.Ctx.Output.Header("Transfer-Encoding", "chunked")
	Once := bytes.Buffer{} //要输出的数据保存对象
	if c.Pjax {
		c.OncePutOUT = true
		Once.WriteString("<title>")
		Once.WriteString(c.Data["title"].(string))
		Once.WriteString("</title><body>")
	}
	regex, err := regexp.Compile("^\\[#\\[(.*?)\\]#\\]([\\s\\S]*)")
	if len(templateName) > 0 {
		//输出数据
		writer := func(c *LcfglyController, OncePutOut io.Writer, data []byte) error {
			if c.OncePutOUT {
				_, err := OncePutOut.Write(data)
				return err
			}
			_, err := c.Ctx.ResponseWriter.Write(data)
			return err
		}
		for index, name := range templateName {
			data := make(map[string]interface{})
			c.TplName = name
			rb, err := c.RenderBytes() //获取到解析出来的模版数据
			if err != nil {
				return err
			}
			if regex.Match(rb) { //匹配到子模版
				if index == 0 && !c.Pjax {
					c.TplName = utils.Conf.DefaultTemplate
					body, err := c.RenderBytes()
					if err != nil {
						return err
					}
					err = writer(c, &Once, body) //输出body数据
					if err != nil {
						return err
					}
				}
				rByte := regex.FindSubmatch(rb) //获取到正则匹配的数据
				data["id"] = string(rByte[1])
				reg, _ := regexp.Compile("\"")
				rByte[2] = reg.ReplaceAll(rByte[2], []byte("\\\""))
				reg, err = regexp.Compile("</script>")
				data["html"] = string(reg.ReplaceAll(rByte[2], []byte("<\\/script>")))
				rByte[2], _ = json.Marshal(data)
				writer(c, &Once, []byte("<script>bigpipe.onPageletArrive("))
				writer(c, &Once, rByte[2])
				writer(c, &Once, []byte(");</script>"))
			} else {
				writer(c, &Once, rb)
			}
		}
	}
	if Once.Len() > 0 { //一次性输出
		s := Once.String()
		reg, err := regexp.Compile("<\\/script><script>") //pjax多个script只会执行一个...BUG
		if err == nil && reg.MatchString(s) {
			s = reg.ReplaceAllString(s, "")
			Once.Truncate(0)    //清空buffer
			Once.WriteString(s) //重新写入数据
		}
	}
	Once.WriteString("</body>")
	if !c.Pjax {
		Once.WriteString("</html>")
	}
	_, err = c.Ctx.ResponseWriter.ResponseWriter.Write(Once.Bytes())
	return err
}

//Error404 404错误
func (c *LcfglyController) Error404() {
	beego.Debug("302:", urlpath.Host("main")+c.Ctx.Request.RequestURI)
	c.Redirect(urlpath.Host("main")+c.Ctx.Request.RequestURI, 302)
	c.EnableRender = false
}

//JSONSuccess 返回成功
//@Param data interface data数据
//@Param msg string 消息内容
//@Param redirect string 跳转地址
func (c *LcfglyController) JSONSuccess(data interface{}, msg string, redirect ...string) {
	c.Data["code"] = 0
	c.Data["data"] = data
	c.Data["msg"] = msg
	if len(redirect) > 0 {
		c.Data["redirect"] = redirect[0]
	}
	c.ServeJSON()
}

//ServeJSON 输出JSON数据
func (c *LcfglyController) ServeJSON() {
	c.Controller.Data["json"] = c.Data
	c.Controller.ServeJSON()
}

//JSONFail 返回错误
//@Param msg string 错误内容
//@Param other... interface 跳转地址,错误代码
//@Param code int 错误代码..默认500
func (c *LcfglyController) JSONFail(msg string, other ...interface{}) {
	if len(other) > 1 {
		c.Data["redirect"] = other[0]
		c.Data["code"] = other[1]
	} else if len(other) > 0 {
		c.Data["redirect"] = other[0]
		c.Data["code"] = 500
	} else {
		c.Data["code"] = 500
	}
	c.Data["msg"] = strings.TrimSpace(strings.TrimPrefix(msg, "com.lcfgly.hw123.exception.ApiException:"))
	c.ServeJSON()
}

//JSONFailError 错误异常
//@Param err error 错误
//@Param redirect string 跳转地址
func (c *LcfglyController) JSONFailError(err error, redirect ...string) {
	exp := services.ErrorToException(err)
	c.JSONFailException(exp, redirect...)
}

//JSONFailException 错误异常
//@Param err error 错误
//@Param redirect string 跳转地址
func (c *LcfglyController) JSONFailException(exp services.Exception, redirect ...string) {
	if len(redirect) > 0 {
		c.Data["redirect"] = redirect[0]
	}
	c.Data["code"] = exp.Code()
	c.Data["msg"] = exp.Error()
	c.ServeJSON()
}

//ErrorPage 错误页面提示
func (c *LcfglyController) ErrorPage(msg string, redirect ...string) {
	c.Data["msg"] = msg
	if redirect != nil && len(redirect) > 0 {
		c.Data["redirect"] = redirect[0]
	}
	c.Data["title"] = "错误"
	c.LoadTop = true
	c.TplName = "error.html" //显示错误页面
}
