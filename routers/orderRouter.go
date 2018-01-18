package routers

import (
	"github.com/astaxie/beego"
	"gitee.com/sparon/goshop/controllers"
	"gitee.com/sparon/goshop/utils/urlpath"
)

func init() {
	OrderCtl := &controllers.OrderController{}
	//购物车列表
	beego.Router(urlpath.Path("cart"), OrderCtl, "get:GetCart")
	//迷你购物车
	beego.Router(urlpath.Path("cartview"), OrderCtl, "get:GetMinCartView")
	//添加商品到购物车
	beego.Router(urlpath.Path("addcart"), OrderCtl, "*:AddCart")
	//POST生成订单信息
	beego.Router(urlpath.Path("order"), OrderCtl, "post:PostOrder")
	//Get生成订单信息
	beego.Router(urlpath.Path("order"), OrderCtl, "get:GetOrder")
	//提交订单
	beego.Router(urlpath.Path("orderCreate"), OrderCtl, "post:PostCreateOrder")
	//生成支付单据
	beego.Router(urlpath.Path("createPayBill"), OrderCtl, "*:CratePayBill")
	//待支付页面
	beego.Router(urlpath.Path("waitpay"), OrderCtl, "get:GetWaitPay")
	//订单支付
	beego.Router(urlpath.Path("pay"), OrderCtl, "post:PostPay")
	//订单支付结果
	beego.Router(urlpath.Path("pay"), OrderCtl, "get:GetPayResult")
}
