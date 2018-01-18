package services

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

//ExceptionMap 错误代码集合.以免统一错误多次创建
var exceptionMap = map[int]Exception{}

//Exception 异常对象
type Exception interface {
	error
	Code() int //错误代码
}

type errorException struct {
	Msg string //错误消息
	Cod int    //错误代码
}

//错误代码
func (e *errorException) Error() string {
	return e.Msg
}

//Code 错误代号
func (e *errorException) Code() int {
	return e.Cod
}

//ErrorToException error类型转换成Exception
func ErrorToException(err error) Exception {
	reg, er := regexp.Compile("\\[(.*?)\\]")
	if er == nil && reg.MatchString(err.Error()) { //匹配到数据
		codestring := reg.FindStringSubmatch(err.Error())[1]
		code, er := strconv.Atoi(codestring)
		if er == nil {
			if v, ok := exceptionMap[code]; ok { //该错误已存在
				return v
			}
			exp := &errorException{}
			exp.Cod = code
			exp.Msg = beego.Substr(err.Error(), strings.Index(err.Error(), "]")+1, len(err.Error()))
			return exp
		}
	}
	return exception500(err)
}

//Exception500 未知异常
func exception500(err error) Exception {
	exp := &errorException{}
	exp.Cod = 500
	exp.Msg = err.Error()
	return exp
}
