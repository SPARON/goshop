package services

import (
	"gitee.com/sparon/goshop/models"
)

//ActivityQiangou 抢购信息
//@Param memberID int64 会员ID
//@Param productID int64 货品ID
//@return 活动信息
func ActivityQiangou(memberID int64, productID int64) (*models.ActivityBean, error) {
	return models.Activity.ActivityQiangou(memberID, productID)
}

//RushBuy 抢购
//@Param memberID int64 会员ID
//@Param productID int64 货品ID
//@Param ip string 客户IP
//@Param source string 订单来源
//@return 订单号
func RushBuy(memberID int64, productID int64, ip string, source string) (int64, error) {
	return RushBuy(memberID, productID, ip, source)
}
