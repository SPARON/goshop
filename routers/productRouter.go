package routers

import (
	"github.com/astaxie/beego"
	"gitee.com/sparon/goshop/controllers"
	"gitee.com/sparon/goshop/utils/urlpath"
)

func init() {
	ProductCtl := &controllers.ProductController{}
	//Ajax获取商品基本信息
	beego.Router(urlpath.Path("product")+"/basic/ajax/:product", ProductCtl, "*:ProductAjax")
	//商品基本信息
	beego.Router(urlpath.Path("product")+"/:product", ProductCtl, "get:GetProduct")
	//获取商品可售库存
	beego.Router(urlpath.Path("product_store"), ProductCtl, "get:GetProductStore")
	//获取商品价格
	beego.Router(urlpath.Path("product_price"), ProductCtl, "get:GetProductPrice")
	//先关商品
	beego.Router(urlpath.Path("goodsLink"), ProductCtl, "get:GetGoodsLink")
	//商品评论
	beego.Router(urlpath.Path("goodsDiscuss"), ProductCtl, "*:GoodsDiscuss")
	//商品评论
	beego.Router(urlpath.Path("goodsDiscussInit"), ProductCtl, "*:GoodsDiscuss")
	//销售记录
	beego.Router(urlpath.Path("/goodsSellLoglist"), ProductCtl, "*:SellLogList")
}
