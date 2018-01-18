package models

import (
	"reflect"
	"strings"
	"time"

	"github.com/hprose/hprose-go"
)

//OrderBean 订单数据操作
var OrderBean *OrderModel

//OrderModel 订单模型
type OrderModel struct {
	//订单数据处理
	//@Param products map[uint64]uint64 产品信息
	//@Param memberID uint64 会员ID
	ExecOrder func(products map[uint64]uint64, memberID uint64) (*OrderInfo, error)
	//检出购物车
	//@Param memberID uint64 会员ID
	CheckOutCart func(memberID uint64) map[uint64]int64 `name:"checkOut"`
	//添加产品到购物车
	//@Param memberID uint64 用户ID
	//@Param productID uint64 货品ID
	//@Param nums int64 数量
	AddCart func(memberID uint64, productID uint64, nums int64) (bool, error) `name:"add"`
	//获取订单信息
	//@Param OrderIdent string 订单号
	CartOrder func(OrderIdent string) (*OrderInfo, error)
	//订单信息增加快递数据
	//@Param order OrderInfo 订单信息
	ExpressList func(order *OrderInfo) error `byref:"true"`
	//创建订单
	//@Param args map[string]interface{} 订单参数
	//@Param ident string //订单标识
	//@Param memberID uint64 // 用户ID
	CreateOrder func(args map[string]interface{}, ident string, memberID uint64) (*CreateOrderResult, error)
	//创建支付订单
	//@Param memberID uint64 用户ID
	//@Param orderids ...interface{} 订单号集合
	CreatePayBill func(memberID uint64, orderids []interface{}) ([]interface{}, error)
	//付款
	//@Param pay PaymentType 付款方式
	//@Param ip string 客户端IP
	//@Param billID uint64 支付单据
	//@Param memberID uint64 用户ID
	//@Param point int 用户消费积分
	Pay func(pay PaymentType, ip string, billID uint64, memberID uint64, point int) (string, error)
	//获取支付单据信息
	//@Param bill uint64 单据ID
	GetPayBill func(bill uint64) (PayBill, error)
}

func initOrder() {
	hproseClient.UseService(&OrderBean)
	hprose.ClassManager.Register(reflect.TypeOf(OrderProduct{}), "OrderProduct", "v")
	hprose.ClassManager.Register(reflect.TypeOf(ExpressLocal{}), "ExpressLocal", "v")
	hprose.ClassManager.Register(reflect.TypeOf(Express{}), "Express", "v")
	hprose.ClassManager.Register(reflect.TypeOf(ConditionDiscount{}), "ConditionDiscount", "v")
	hprose.ClassManager.Register(reflect.TypeOf(Supplier{}), "Supplier", "v")
	hprose.ClassManager.Register(reflect.TypeOf(OrderSupplier{}), "OrderSupplier", "v")
	hprose.ClassManager.Register(reflect.TypeOf(OrderInfo{}), "OrderInfo", "v")
	hprose.ClassManager.Register(reflect.TypeOf(CreateOrderResult{}), "CreateOrderResult", "v")
}

//OrderProduct 订单中商品信息
type OrderProduct struct {
	Product  *Product `v:"product"`  //货品
	Nums     uint64   `v:"nums"`     //购买数量
	GaveNums uint64   `v:"gaveNums"` //赠送数量
}

//ExpressLocal 快递各个地区信息
type ExpressLocal struct {
	Region     uint64  `v:"region"`     //地区ID
	Express    uint64  `v:"express"`    //快递ID
	SupplierID uint64  `v:"supplier"`   //供应商ID
	FirstCost  float64 `v:"firstCost"`  //首重价格
	ExtendcOST float64 `v:"extendCost"` //续重价格
}

//Express 快递信息
type Express struct {
	ID          uint64          `v:"id"`           //ID
	Name        string          `v:"name"`         //快递名称
	HasCod      bool            `v:"hasCod"`       //支付货到付款
	FirstUnit   int64           `v:"firstunit"`    //首重
	ExtendUnit  int64           `v:"extendunit"`   //续重
	Protect     bool            `v:"protect"`      //保价
	ProtectRate float32         `v:"protect_rate"` //保价费率
	MinPrice    float64         `v:"minprice"`     //最低保价金额
	Setting     bool            `v:"setting"`      //是否统计设置收费
	DefAreaFee  bool            `v:"defAreaFee"`   //地区配置不存在时是否按照默认数据计算费用
	FirstPrice  float64         `v:"firstPrice"`   //首重价格
	ExtendPrice float64         `v:"extendPrice"`  //续重价格
	CorpID      int64           `v:"corp_id"`      //物流公司ID
	OrderNum    int16           `v:"ordernum"`     //排序
	Disabled    bool            `v:"disabled"`     //是否禁用
	Supplier    uint64          `v:"supplier"`     //供应商ID
	Status      bool            `v:"status"`       //是否开启
	LocalInfo   []*ExpressLocal `v:"localInfo"`    //快递地区信息
}

//ConditionDiscount 优惠活动
type ConditionDiscount struct {
	ConditionDiscountID int64     `v:"condition_discount_id"` //ID
	DiscountID          int64     `v:"discount_id"`           //优惠ID
	ConditionID         int64     `v:"condition_id"`          //条件ID
	Begin               time.Time `v:"begin"`                 //开始时间
	End                 time.Time `v:"end"`                   //结束时间
	Sort                int16     `v:"sort"`                  //排序
	CdName              string    `v:"cdname"`                //优惠活动名称
	CdType              bool      `v:"cdtype"`                //是否唯一
	Createtime          time.Time `v:"createtime"`            //创建时间
	Disabled            bool      `v:"disabled"`              //是否禁用
	SupplierID          uint64    `v:"supplier_id"`           //供应商ID
	Description         string    `v:"description"`           //详细描述
}

//Supplier 供应商信息
type Supplier struct {
	ID         uint64 `v:"Id"`         //ID
	Name       string `v:"name"`       //供应商名称
	Bn         string `v:"bn"`         //供应商编码
	LinkMan    string `v:"linkMan"`    //供应商联系人
	Phone      string `v:"phone"`      //供应商联系人电话
	Disabled   bool   `v:"disabled"`   //是否有效
	Address    string `v:"address"`    //供应商地址
	CreateTime string `v:"createtime"` //创建时间
}

//OrderSupplier 订单中供应商信息
type OrderSupplier struct {
	Supplier           *Supplier            `v:"supplier"`           //供应商
	Memo               string               `v:"memo"`               //订单留言
	Discount           float64              `v:"discount"`           //优惠金额
	CostFreight        float64              `v:"cost_freight"`       //快递费用
	FreeCost           bool                 `v:"free_cost"`          //是否免邮
	Products           []*OrderProduct      `v:"products"`           //商品列表
	Express            []*Express           `v:"expresses"`          //快递列表
	ConditionDiscounts []*ConditionDiscount `v:"conditionDiscounts"` //优惠列表
	Weight             float32              `v:"weight"`             //订单重量
	Cost               float64              `v:"cost"`               //总金额
	SubTotal           float64              `v:"subtotal"`           //订单总金额
}

//OrderInfo 订单信息
type OrderInfo struct {
	MemberID    uint64           `v:"member_id"`    //会员ID
	ConsigneeID uint64           `v:"consignee_id"` //收获地址ID
	Ident       string           `v:"ident"`        //唯一标识
	Consignees  []*ConsigneeInfo `v:"consignees"`   //收获地址列表
	Suppliers   []*OrderSupplier `v:"suppliers"`    //供应商列表
	Weight      float32          `v:"weight"`       //订单重量
	SubTotal    float64          `v:"subtotal"`     //商品总金额
	Cost        float64          `v:"cost"`         //订单总金额
	CostFreight float64          `v:"cost_freight"` //快递费用
}

//CreateOrderResult 订单创建结果信息
type CreateOrderResult struct {
	OrderIDS      []uint64                 `v:"order_ids"`     //订单号集合
	APIExceptions []map[string]interface{} `v:"apiExceptions"` //订单异常
}

//PaymentType 支付方式
type PaymentType int

const (
	AliPay   PaymentType = iota //支付宝支付
	WxPay                       //微信支付
	ChinaPay                    //银联支付
	Deposit                     //预存款
)

//StringToPaymentType 获取支付方式
func StringToPaymentType(name string) PaymentType {
	switch strings.ToLower(name) {
	case "wxpay":
		return WxPay
	case "chinapay":
		return ChinaPay
	case "alipay":
		return AliPay
	default:
		return Deposit
	}
}

//PayBill 单据支付结果
type PayBill map[string]interface{}
