package models

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"github.com/hprose/hprose-go"
	"gitee.com/sparon/goshop/utils"
	"gitee.com/sparon/goshop/utils/urlpath"
)

//Login 操作数据
var Login *LoginModel

//LoginModel 登入接口
type LoginModel struct {
	//LoginUser 用户账号密码检测
	//@Param username string 用户名
	//@Param password string 密码
	//@return 用户ID
	LoginUser func(username string, password string) (int64, error)
	//AddMember 添加用户
	//@Param member Member 会员信息
	//@Param pamMember PamMember 会员登入信息
	//@return 新用户ID
	AddMember func(member Member, pamMember PamMember) (int64, error)
	//CheckUser 检测用户是否可用
	//@Param username string 用户名
	CheckUser func(username string) bool
	//ecstore加密用户密码
	//@Param password string 密码
	//@Param password_account string 加密字符串
	//@Param createtime int64 创建时间
	//@Return 加密后的字符串
	EcstoreEncodeMemberPassword func(password, password_account string, createtime int64) string
	//创建修改密码的临时密钥
	//@Param memberID uint64 会员ID
	//@Param time int64 过期时间戳
	//@Return 临时密钥
	ChangePwdCreateKey func(memberID uint64, time int64) (string, error) `name:"createMemberPwdLog"`

	/**
	 * 获取可用的临时密钥信息
	 * @param secret 密钥
	 * @param memberID 会员ID
	 * @return
	 */
	MemberPwdLog func(secret string, memberID uint64) (map[string]interface{}, error)

	/**
	 * 修改临时密钥成已使用
	 * @param pwdlogID
	 */
	UpdateMemberPwdLogUsed func(pwdlogID int)
}

var Sms *SmsModel

//SmsModel 短信接口
type SmsModel struct {

	//SendVerifyCode 发送验证码短信
	//@Param phonenum string 手机号
	//@Param vcode string 验证码
	SendVerifyCode func(phonenum, vcode string) (bool, error)

	//Send 发送短信
	//@Param phonenum string 手机号
	//@Param content string 短信内容
	Send func(phonenum, content string) (bool, error)

	//SendEmail 发送邮件
	//@Param email string 邮箱地址
	//@Param title string 邮件标题
	//@Param content string 邮件内容
	SendEmail func(email, title, content string) (bool, error) `name:"Send"`

	//SendEmails 群发邮件
	//@Param emails []string 邮件地址集合
	//@Param title string 邮件标题
	//@Param content string 邮件内容
	SendEmails func(emails []string, title, context string) (bool, error) `name:"Send"`

	//GenerateVerifyCode 生成验证码
	//@Param verifySize int 验证码长度
	//@Return 验证码
	GenerateVerifyCode func(verifySize int) string
}

//初始化登入模型
func initLoginModel() {
	hproseClient.UseService(&Login) //获取模型实例
	hproseClient.UseService(&Sms)   //邮件模型实例
	hprose.ClassManager.Register(reflect.TypeOf(PamMember{}), "PamMember", "v")
	hprose.ClassManager.Register(reflect.TypeOf(Member{}), "Member", "v")
}

//MemberSource 用户来源枚举
type MemberSource int

const (
	PC     MemberSource = iota //PC
	WeiXin                     //微信
	API                        //API
	Wap                        //移动端
)

func (m MemberSource) String() string {
	switch m {
	case WeiXin:
		return "weixin"
	case API:
		return "api"
	case Wap:
		return "wap"
	}
	return "pc"
}

//StringConvertMemberSource string转换成 MemberSource
func StringConvertMemberSource(str string) MemberSource {
	switch strings.ToLower(str) {
	case "weixin":
		return WeiXin
	case "api":
		return API
	case "wap":
		return Wap
	}
	return PC
}

//IntConvertMemberSource int转换MemberSource
func IntConvertMemberSource(in int) MemberSource {
	switch MemberSource(in) {
	case WeiXin:
		return WeiXin
	case API:
		return API
	case Wap:
		return Wap
	}
	return PC
}

//Member 用户对象
type Member struct {
	MemberID   uint64       `v:"member_id"`    //会员ID
	MemberLvID uint64       `v:"member_lv_id"` //会员等级
	Point      int64        `v:"point"`        //积分
	Mobile     string       `v:"mobile"`       //手机号
	ReferID    string       `v:"refer_id"`     //推广员ID
	ReferURL   string       `v:"refer_url"`    //推广来源
	RegIP      string       `v:"reg_ip"`       //注册IP
	RegTime    int64        `v:"regtime"`      //注解时间
	CUR        string       `v:"cur"`          //默认货币
	Email      string       `v:"email"`        //电子邮箱
	Source     MemberSource `v:"source"`       //注册平台
}

//LoginTypeEnum 登入类型枚举
type LoginTypeEnum int

const (
	Local  LoginTypeEnum = iota //用户
	Email                       //邮箱
	Mobile                      //手机
)

func (m LoginTypeEnum) String() string {
	switch m {
	case Email:
		return "email"
	case Mobile:
		return "mobile"
	}
	return "local"
}

//StringConvertLoginTypeEnum string转换成LoginTypeEnum
func StringConvertLoginTypeEnum(str string) LoginTypeEnum {
	switch strings.ToLower(str) {
	case "mobile":
		return Mobile
	case "email":
		return Email
	}
	return Local
}

//IntConvertLoginTypeEnum int转换成LoginTypeEnum
func IntConvertLoginTypeEnum(in int) LoginTypeEnum {
	switch LoginTypeEnum(in) {
	case Mobile:
		return Mobile
	case Email:
		return Email
	}
	return Local
}

//LoginTypeEnum 登入类型枚举
type BOOLEnum int

const (
	True  BOOLEnum = iota //true
	False                 //false
)

func (m BOOLEnum) String() string {
	if m == True {
		return "true"
	}
	return "false"
}

//StringConvertBOOLEnum string转换成BOOLEnum
func StringConvertBOOLEnum(str string) BOOLEnum {
	if strings.ToLower(str) == "true" {
		return True
	}
	return False
}

//IntConvertBOOLEnum int转换成BOOLEnum
func IntConvertBOOLEnum(in int) BOOLEnum {
	if BOOLEnum(in) == True {
		return True
	}
	return False
}

//PamMember 会员登入信息
type PamMember struct {
	MemberID        uint64        `v:"member_id"`        //会员ID
	LoginType       LoginTypeEnum `v:"login_type"`       //登入类型
	Account         string        `v:"account"`          //登入名
	Password        string        `v:"password"`         //密码
	PasswordAccount string        `v:"password_account"` //加密字符串
	CreateTime      int64         `v:"createtime"`       //创建时间
	Disabled        BOOLEnum      `v:"disabled"`         //是否禁用
}

//IsLogin 检测用户是否登入
//@Param ctx context.Context
//@Param redirectLogin bool 为登入时是否跳转到登入页面
//@Return 用户登入ID
func IsLogin(r *http.Request, w http.ResponseWriter, redirectLogin ...bool) (memberID uint64) {
	loginIDCookie, err := r.Cookie(utils.Conf.CookieLoginID) //登入数据ID
	loginID := ""
	if err == nil && loginIDCookie != nil {
		loginID = loginIDCookie.Value
	}
	loginCookie, err := r.Cookie(utils.Conf.CookieLogin) //登入加密信息
	login := ""
	if err == nil && loginCookie != nil {
		login = loginCookie.Value
	}
	defer func() {
		//未登入时判断是否进行跳转
		if memberID == 0 && redirectLogin != nil && len(redirectLogin) > 0 && redirectLogin[0] {
			http.Redirect(w, r, urlpath.Host("sso")+urlpath.Path("login"), 302)
		}
	}()
	var data uint64 //用户登入信息json数据
	if loginID != "" {
		idbyte, err := utils.AesDecrypt([]byte(loginID))
		if err != nil { //信息解密失败
			beego.Debug("用户登入信息LoginID解密失败", err)
			return 0
		}
		id := string(idbyte)
		if strings.HasPrefix(id, "D") { //标识数据存放在cookie信息中
			dd, err := utils.AesDecrypt([]byte(login), []byte(beego.Substr(id, 1, len(id)-1)))
			if err != nil {
				beego.Debug("用户数据信息Login解析失败", err)
				return 0
			}
			data, err = strconv.ParseUint(string(dd), 10, 64)
			if err != nil {
				beego.Debug("用户数据信息Login解析失败", err)
				return 0
			}
		} else { //否则从redis中获取登入信息
			conn := utils.RedisClient.Get() //从连接池中获取一个连接
			defer func() {
				conn.Close() //关闭redis连接
			}()
			data, err = redis.Uint64(conn.Do("GET", utils.RedisLoginKey(id)))
			if err != nil {
				beego.Debug("Redis数据获取失败", err)
				return 0
			}
			conn.Do("EXPIRE", utils.RedisLoginKey(id), utils.Conf.CookieExpireTime) //更新过期时间
		}
	}
	return data
}
