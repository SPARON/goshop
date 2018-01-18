package utils

import (
	"github.com/astaxie/beego"
)

//Config 基本配置信息
type Config struct {
	HproseClient      string //Hprose服务器地址配置
	CookieDomain      string //Cookie域名
	CookieKey         string //Cookie加密密钥
	CookieExpireTime  int64  //Cookie过期时间
	CookieLoginID     string //登入ID信息cookie名称
	CookieLogin       string //登入信息cookie名称
	CookieNoPostOrder string //未提交的订单信息
	CookieCartMember  string //购物车标识
	SecretKey         []byte //加密密钥
	RedisHost         string //Redis连接地址
	RedisDb           string //Redis数据库名称
	RedisMaxidle      int    //最大的空闲连接数
	RedisMaxActive    int    //最大的激活连接数,表示同时最多有N个连接
	DefaultTemplate   string //默认页面框架地址
}

//Conf 配置信息
var Conf *Config

//初始化Config信息
//该初始化函数在redisUtils-init方法中调用
func initConfig() {
	Conf = &Config{}
	Conf.HproseClient = beego.AppConfig.String("hproseClient")
	Conf.CookieDomain = beego.AppConfig.DefaultString("cookie.domain", ".hw123.com")
	Conf.CookieKey = beego.AppConfig.String("cookie.key")
	Conf.CookieExpireTime = beego.AppConfig.DefaultInt64("cookie.expireTime", 86400)
	Conf.CookieLoginID = beego.AppConfig.String("cookie.loginId")
	Conf.CookieLogin = beego.AppConfig.String("cookie.login")
	Conf.CookieNoPostOrder = beego.AppConfig.String("cookie.noPostOrder")
	Conf.SecretKey = []byte(beego.AppConfig.String("secretkey"))
	Conf.RedisHost = beego.AppConfig.String("redis.host")
	Conf.RedisDb = beego.AppConfig.String("redis.db")
	Conf.RedisMaxidle = beego.AppConfig.DefaultInt("redis.maxidle", 1)
	Conf.RedisMaxActive = beego.AppConfig.DefaultInt("redis.maxactive", 10)
	Conf.DefaultTemplate = beego.AppConfig.String("defaultTemplate")
	Conf.CookieCartMember = beego.AppConfig.String("cookie.cartMember")
}
