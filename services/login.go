package services

import (
	"time"

	"github.com/astaxie/beego/context"
	"github.com/sparon/goshop/models"
	"github.com/sparon/goshop/utils"
)

//CheckUser 检测用户是否可用
//@Param username string 用户名
func CheckUser(username string) bool {
	return models.Login.CheckUser(username)
}

//Register 注册用户
//@Param username string 用户名
//@Param password string 密码
//@Param other...string 其他参数
//  other[0] refer_id int64 推广员ID 默认:0
//  other[1] refer_url string 来源连接地址 默认:""
//  other[2] source string 来源平台 默认:"pc"
//  other[3] ip string IP地址 默认:""
//@Return 新用户ID
func Register(username, password string, other ...string) (int64, Exception) {
	member := models.Member{}
	pamMember := models.PamMember{}
	if len(other) > 3 {
		member.ReferID = other[0]
		member.ReferURL = other[1]
		member.Source = models.StringConvertMemberSource(other[2])
		member.RegIP = other[3]
	} else if len(other) > 2 {
		member.ReferID = other[0]
		member.ReferURL = other[1]
		member.Source = models.StringConvertMemberSource(other[2])
	} else if len(other) > 1 {
		member.ReferID = other[0]
		member.ReferURL = other[1]
		member.Source = models.StringConvertMemberSource("pc")
	} else if len(other) > 0 {
		member.ReferID = other[0]
		member.ReferURL = ""
		member.Source = models.StringConvertMemberSource("pc")
	} else {
		member.ReferID = "0"
		member.ReferURL = ""
		member.Source = models.StringConvertMemberSource("pc")
	}
	//验证用户帐号类型
	if utils.ValidateString(username, utils.Email) {
		pamMember.LoginType = models.Email
	} else if utils.ValidateString(username, utils.Phone) {
		pamMember.LoginType = models.Mobile
	} else {
		pamMember.LoginType = models.Local
	}
	pamMember.Account = username
	pamMember.CreateTime = time.Now().Unix()                                                                //创建时间
	pamMember.Password = models.Login.EcstoreEncodeMemberPassword(password, username, pamMember.CreateTime) ///密码加密方式
	pamMember.PasswordAccount = username
	result, err := models.Login.AddMember(member, pamMember)
	if err != nil {
		return result, ErrorToException(err)
	}
	return result, nil
}

//SendSmsVcode 发送短信验证码
//@Param phone string 手机号
func SendSmsVcode(phone string) bool {
	if !utils.ValidateString(phone, utils.Phone) {
		return false
	}
	vcode := models.Sms.GenerateVerifyCode(4)
	ok, err := models.Sms.SendVerifyCode(phone, vcode)
	if err == nil && ok {
		return true
	}
	return false
}

//SendEmail 发送邮件
//@Param email string 邮箱地址
//@Param title string 邮件标题
//@Param content string 邮件内容
func SendEmail(email, title, content string) bool {
	b, err := models.Sms.SendEmail(email, title, content)
	if err != nil {
		return false
	}
	return b
}

//ChangePwdCreateKey 创建修改密码的临时密钥
func ChangePwdCreateKey(memberID uint64) (string, error) {
	return models.Login.ChangePwdCreateKey(memberID, 2*60*60)
}

//MemberPwdLog 获取可用的临时密钥信息
//@param secret 密钥
//@param memberID 会员ID
func MemberPwdLog(secret string, memberID uint64) (map[string]interface{}, error) {
	return models.Login.MemberPwdLog(secret, memberID)
}

//UpdateMemberPwdLogUsed 修改临时密钥成已使用
//@Param pwdlogID
func UpdateMemberPwdLogUsed(pwdlogID int) {
	models.Login.UpdateMemberPwdLogUsed(pwdlogID)
}

//IsLogin 检测用户是否登入
//@Param ctx context.Context
//@Param redirectLogin bool 为登入时是否跳转到登入页面
//@Return 用户登入ID
func IsLogin(c *context.Context, redirectLogin ...bool) (memberID uint64) {
	return models.IsLogin(c.Request, c.ResponseWriter.ResponseWriter, redirectLogin...)
}
