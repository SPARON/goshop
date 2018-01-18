package services

import (
	"gitee.com/sparon/goshop/models"
)

//MemberLoginAccount 获取指定账号的所有可用的登入账号,以及账号类型
//@Param loginName string 帐号名
//@Return 会员登入信息
func MemberLoginAccount(loginName string) []map[string]interface{} {
	result, err := models.MemberBean.MemberLoginAccount(loginName)
	if err != nil {
		return []map[string]interface{}{}
	}
	return result
}

//UpdatePwd 修改密码
//@Param password string 新密码
//@Param memberID uint64 用户ID
func UpdatePwd(password string, memberID uint64) bool {
	err := models.MemberBean.UpdatePwd(password, memberID)
	if err != nil {
		return false
	}
	return true
}

//LoginPam 获取账号的登入信息
//@Param account 用户帐号
func LoginPam(account string) (*models.PamMember, error) {
	return models.MemberBean.LoginPam(account)
}

//ConsigneesList 会员收货列表
//@Param memberID uint64 会员ID
func ConsigneesList(memberID uint64) []*models.ConsigneeInfo {
	return models.MemberBean.ConsigneesList(memberID)
}

//AddConsignee 添加收货人
//@Param cons ConsigneeInfo 收货人信息
func AddConsignee(cons *models.ConsigneeInfo) (uint64, error) {
	return models.MemberBean.AddConsignee(cons)
}

//UpdateConsignee 更新收货人信息
//@Param cons ConsigneeInfo 收货人信息
func UpdateConsignee(cons *models.ConsigneeInfo) (bool, error) {
	return models.MemberBean.UpdateConsignee(cons)
}

//DeleteConsignee 删除收货人信息
//@Param cons ConsigneeInfo 收货人信息
func DeleteConsignee(cons *models.ConsigneeInfo) (bool, error) {
	return models.MemberBean.DeleteConsignee(cons)
}

//Consignee 获取会员收货地址信息
//@Param memberID uint64 会员ID
//@Param consignessID uint64 地区ID
func Consignee(memberID uint64, consignessID uint64) (*models.ConsigneeInfo, error) {
	return models.MemberBean.Consignee(memberID, consignessID)
}

//DefaultConsignee 获取会员默认收货地址信息
//@Param memberID uint64 会员ID
func DefaultConsignee(memberID uint64) (*models.ConsigneeInfo, error) {
	return models.MemberBean.DefaultConsignee(memberID)
}

//GetMemberPoint 用户可用的积分
//@Param memberID uint64 会员ID
func GetMemberPoint(memberID uint64) int {
	return models.MemberBean.GetMemberPoint(memberID)
}
