package services

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sparon/goshop/models"
)

//ExecOrder 获取商品信息
//@Param params url.Values 请求参数
//@Param memberID uint64 用户ID
func ExecOrder(params url.Values, memberID uint64) (*models.OrderInfo, error) {
	products := make(map[uint64]uint64)
	for k, v := range params {
		if strings.HasPrefix(k, "modify_quantity[") { //解析product
			nums, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				continue
			}
			productid, err := strconv.ParseInt(beego.Substr(k, len("modify_quantity["), len(k)-len("modify_quantity[")-1), 10, 64)
			if err != nil {
				continue
			}
			products[uint64(productid)] = uint64(nums)
		}
	}
	if len(products) > 0 { //如果存在未提交的商品订单,提交订单
		orderinfo, err := models.OrderBean.ExecOrder(products, memberID)
		return orderinfo, err
	}
	return nil, errors.New("商品信息获取失败")
}

//CheckOutCart 检出购物车
//@Param memberID uint64 会员ID
func CheckOutCart(memberID uint64) map[string]interface{} {
	data := make(map[string]interface{})
	carts := models.OrderBean.CheckOutCart(memberID)
	if carts == nil || len(carts) < 1 {
		data["subtotal"] = 0 //金额0
	} else {
		productIDS := make([]uint64, len(carts))
		i := 0
		for k := range carts {
			productIDS[i] = k
			i = i + 1
		}
		products, err := ProductInfos(productIDS)
		if err != nil {
			data["subtotal"] = 0 //金额0
		} else {
			var result []map[string]interface{}
			subtotal := float64(0)
			for _, p := range products {
				if nums, ok := carts[p.ProductID]; ok {
					obj := make(map[string]interface{})
					obj["nums"] = nums
					obj["product"] = p
					obj["price"] = float64(nums) * p.Price
					subtotal = subtotal + obj["price"].(float64)
					result = append(result, obj)
				}
			}
			data["data"] = result
			data["subtotal"] = subtotal
		}
	}
	return data
}

//AddCart 添加产品到购物车
//@Param memberID uint64 用户ID
//@Param productID uint64 货品ID
//@Param nums int64 数量
func AddCart(memberID uint64, productID uint64, nums int64) (map[string]interface{}, error) {
	ok, err := models.OrderBean.AddCart(memberID, productID, nums)
	if ok {
		cartList := CheckOutCart(memberID) //获取购物车信息
		result := make(map[string]interface{})
		if v, ok := cartList["data"]; ok { //购物车中存在产品数据
			for _, vv := range v.([]map[string]interface{}) {
				p := vv["product"].(*models.Product)
				if p.ProductID == productID {
					result["product_price"] = p.Price
					break
				}
			}
			result["count"] = len(v.([]map[string]interface{}))
		} else {
			result["count"] = 0
		}
		if _, ok := result["product_price"]; ok {
			result["product_price"] = 0
		}
		result["subtotal"] = cartList["subtotal"]
		return result, nil
	}
	return nil, err
}

//CartOrder 获取订单信息
//@Param OrderIdent string 订单号
func CartOrder(OrderIdent string) (*models.OrderInfo, error) {
	res, err := models.OrderBean.CartOrder(OrderIdent)
	if err != nil || res == nil {
		return nil, errors.New("订单不存在")
	}
	err = models.OrderBean.ExpressList(res)
	for _, con := range res.Consignees {
		if con.Areas == nil {
			con.Areas = strings.Split(strings.Split(con.Area, ":")[1], ".")
		}
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

//CreateOrder 创建订单
//@Param args map[string]interface{} 订单参数
//@Param ident string //订单标识
//@Param memberID uint64 // 用户ID
func CreateOrder(args map[string]interface{}, ident string, memberID uint64) (*models.CreateOrderResult, error) {
	return models.OrderBean.CreateOrder(args, ident, memberID)
}

//CreatePayBill 创建支付订单
//@Param memberID uint64 用户ID
//@Param orderids ...interface{} 订单号集合
func CreatePayBill(memberID uint64, orderids ...interface{}) (map[string]interface{}, error) {
	data, err := models.OrderBean.CreatePayBill(memberID, orderids)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	result["bill"] = data[0]
	result["orders"] = data[1]
	result["amount"] = data[2]
	return result, nil
}

//Pay 付款
//@Param pay int 付款方式
//@Param ip string 客户端IP
//@Param billID uint64 支付单据
//@Param memberID uint64 用户ID
//@Param point int 消费积分
func Pay(pay int, ip string, billID uint64, memberID uint64, point int) (string, error) {
	payment, err := paymentType(pay)
	if err != nil {
		return "", err
	}
	return models.OrderBean.Pay(payment, ip, billID, memberID, point)
}

//paymentType 转换支付类型
func paymentType(tp int) (models.PaymentType, error) {
	switch models.PaymentType(tp) {
	case models.AliPay:
		return models.AliPay, nil
	case models.WxPay:
		return models.WxPay, nil
	case models.ChinaPay:
		return models.ChinaPay, nil
	case models.Deposit:
		return models.Deposit, nil
	default:
		return models.AliPay, errors.New("未知支付方式")
	}
}

//GetPayBill 获取支付单据信息
//@Param bill uint64 单据ID
func GetPayBill(bill uint64) (models.PayBill, error) {
	result, err := models.OrderBean.GetPayBill(bill)
	if err != nil {
		return nil, err
	}
	if result["type"].(int) == 0 {
		result["type"] = "退款"
	} else {
		result["type"] = "付款"
	}
	switch result["status"].(int) {
	case 3:
		result["status"] = "Succ"
	case 2:
		result["status"] = "close"
		result["errormsg"] = "单据已经关闭"
	default:
		result["status"] = "failed"
		result["errormsg"] = "支付失败"
	}
	return result, nil
}
