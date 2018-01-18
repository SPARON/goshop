package utils

//字符串操作扩展函数

import (
	"regexp"
)

//ValidateMode 验证模式
type ValidateMode int

const (
	Phone  ValidateMode = iota //手机号
	Email                      //邮箱
	Number                     //数字
	Tel                        //电话
)

//ValidateString 验证字符串
func ValidateString(str string, v ValidateMode) bool {
	var comp string
	switch v {
	case Number:
		comp = "^\\s+$"
	case Email:
		comp = "/^[a-z\\d][a-z\\d_.]*@[\\w-]+(?:\\.[a-z]{2,})+$/i"
	case Phone:
		comp = "^0?1[34578]\\d{9}$"
	case Tel:
		comp = "/^(0\\d{2,3}-?)?[2-9]\\d{5,7}(-\\d{1,5})?$/"
	}
	reg, err := regexp.Compile(comp)
	if err != nil {
		return false
	}
	return reg.MatchString(str)
}
