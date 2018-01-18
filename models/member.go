package models

import (
	"reflect"

	"github.com/hprose/hprose-go"
)

//MemberBean 会员模型对象
var MemberBean *MemberModel

//MemberModel 会员模型
type MemberModel struct {

	//MemberLoginAccount 获取指定账号的所有可用的登入账号,以及账号类型
	//@Param loginName string 帐号名
	//@Return 会员登入信息
	MemberLoginAccount func(loginName string) ([]map[string]interface{}, error)

	//修改密码
	//@Param password string 新密码
	//@Param memberID uint64 用户ID
	UpdatePwd func(password string, memberID uint64) error

	//获取账号的登入信息
	//@Param account 用户帐号
	LoginPam func(account string) (*PamMember, error)

	//获取会员收货地址信息
	//@Param memberID uint64 会员ID
	//@Param consignessID uint64 地区ID
	Consignee func(memberID uint64, consignessID uint64) (*ConsigneeInfo, error)

	//获取会员默认收货地址信息
	//@Param memberID uint64 会员ID
	DefaultConsignee func(memberID uint64) (*ConsigneeInfo, error)

	//会员收货列表
	//@Param memberID uint64 会员ID
	ConsigneesList func(memberID uint64) []*ConsigneeInfo

	//添加收货人
	//@Param cons ConsigneeInfo 收货人信息
	AddConsignee func(cons *ConsigneeInfo) (uint64, error)

	//更新收货人信息
	//@Param cons ConsigneeInfo 收货人信息
	UpdateConsignee func(cons *ConsigneeInfo) (bool, error)

	//删除收货人信息
	//@Param cons ConsigneeInfo 收货人信息
	DeleteConsignee func(cons *ConsigneeInfo) (bool, error)

	//用户可用的积分
	//@Param memberID uint64 会员ID
	GetMemberPoint func(memberID uint64) int `name:"getPoint"`
}

func initMember() {
	hproseClient.UseService(&MemberBean)
	hprose.ClassManager.Register(reflect.TypeOf(ConsigneeInfo{}), "ConsigneeInfo", "v")
}

//ConsigneeInfo 收货人信息
type ConsigneeInfo struct {
	ID       uint64   `v:"id"`        //收货地址编号
	MemberID uint64   `v:"member_id"` //会员ID
	Name     string   `v:"name"`      //收货人姓名
	Phone    string   `v:"phone"`     //联系方式
	LocalID  uint64   `v:"local_id"`  //地区编号
	Address  string   `v:"address"`   //详细地址
	Zip      string   `v:"zip"`       //邮编
	Tel      string   `v:"tel"`       //电话
	Area     string   `v:"area"`      //地区
	DefAddr  bool     `v:"def_addr"`  //默认收获地址
	Areas    []string `v:"areas"`     //地区分割后的数据
}
