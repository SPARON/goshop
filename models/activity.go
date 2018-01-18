package models

import (
	"reflect"
	"time"

	"github.com/hprose/hprose-go"
)

//Activity 活动对象
var Activity *ActivityModel

//ActivityModel 活动
type ActivityModel struct {

	//产品预约
	//@Param memberID int64 会员ID
	//@Param productID int64 预约货品ID
	//@Param nums int64 购买数量
	//@Param localID int64 收货人信息ID
	//@Param ShippingID int64 快递ID
	Yuyue func(memberID int64, productID int64, nums int64, localID int64, ShippingID int64) (bool, error)

	//抢购
	//@Param memberID int64 会员ID
	//@Param productID int64 货品ID
	//@Param ip string 客户IP
	//@Param source string 订单来源
	//@return 订单号
	RushBuy func(memberID int64, productID int64, ip string, source string) (int64, error)

	//货品相关活动信息
	//@Param productID int64 会员ID
	//@Param tp string 活动类型(qiangou＝抢购,xiangou=限购)
	//@return 活动信息
	Activity func(productID int64, tp string) (*ActivityBean, error)

	//用户商品预约信息
	//@Param memberID int64 会员ID
	//@Param product int64 商品ID
	//@return 预约信息
	YuyueInfo func(memberID int64, product int64) map[string]interface{}

	//抢购信息
	//@Param memberID int64 会员ID
	//@Param productID int64 货品ID
	//@return 活动信息
	ActivityQiangou func(memberID int64, productID int64) (*ActivityBean, error)
}

//ActivityBean 活动数据对象
type ActivityBean struct {
	ActivityID int64                  `v:"activityID"` //活动ID
	SupplierID int64                  `v:"supplierID"` //供应商ID
	Name       string                 `v:"name"`       //活动名称
	StartTime  *time.Time             `v:"startTime"`  //活动开始时间
	EndTime    *time.Time             `v:"endTime"`    //活动结束时间
	URL        string                 `v:"url"`        //活动详细地址
	NumsType   int16                  `v:"numstype"`   //活动类型(0=商品数量累积,1=单一)
	Nums       int64                  `v:"nums"`       //商品数量限制
	Type       string                 `v:"type"`       //活动类型标识(qiangou＝抢购,xiangou=限购)
	Extend     map[string]interface{} `v:"extend"`     //附加数据信息
	GoodsType  int16                  `v:"goodsType"`  //作用对象(0=商品,1＝货品)
	Product    *Product               `v:"product"`    //货品对象
}

func initActivity() {
	hprose.ClassManager.Register(reflect.TypeOf(ActivityBean{}), "ActivityBean", "v")
	hproseClient.UseService(&Activity)
}
