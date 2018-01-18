package controllers

import (
	"errors"
	"strconv"

	"github.com/astaxie/beego"

	"gitee.com/sparon/goshop/models"
	"gitee.com/sparon/goshop/services"
)

//ProductController 会员控制器
type ProductController struct {
	LcfglyController
}

//GetProductStore 获取商品可售库存
func (c *ProductController) GetProductStore() {
	productID, err := c.GetInt64("product_id", 0)
	if err != nil || productID == 0 {
		c.JSONFail("缺少参数product_id(货品ID)")
		return
	}
	product, err := services.ProductInfo(uint64(productID))
	if err != nil {
		c.JSONFail(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["store"] = product.Store
	data["title"] = product.Store - product.Freez
	c.JSONSuccess(data, "数据查询成功")
}

//GetProductPrice 获取商品价格
func (c *ProductController) GetProductPrice() {
	productID, err := c.GetInt64("product_id", 0)
	if err != nil || productID == 0 {
		c.JSONFail("缺少参数product_id(货品ID)")
		return
	}
	product, err := services.ProductInfo(uint64(productID))
	if err != nil {
		c.JSONFail(err.Error())
		return
	}
	data := make(map[string]interface{})
	data["price"] = product.Price
	data["mktprice"] = product.MktPrice
	data["minprice"] = nil
	data["saveprice"] = "(节省￥" + strconv.FormatFloat((product.MktPrice-product.Price), 'f', 3, 64) + ")"
	c.JSONSuccess(data, "数据查询成功")
}

//ProductAjax 商品基本信息Ajax
func (c *ProductController) ProductAjax() (resErr error) {
	defer errorPage(&c.LcfglyController, &resErr)
	pid, err := strconv.ParseUint(c.Ctx.Input.Param(":product"), 10, 64)
	if err != nil {
		return errors.New("参数不足")
	}
	product, err := services.ProductInfo(pid)
	if err != nil {
		return errors.New("货品不存在")
	}
	//spec, err := services.GoodsSpec(product.GoodsID)
	//if err != nil {
	//return err
	//}
	//c.Data["spec"] = spec
	c.EnableRender = false
	c.JSONSuccess(product, "数据查询成功")
	//c.Data["product"] = product
	//c.TplName = "product/product_basic.html"
	return nil
}

//GetProduct 商品页
func (c *ProductController) GetProduct() (resErr error) {
	defer errorPage(&c.LcfglyController, &resErr)
	pid, err := strconv.ParseUint(c.Ctx.Input.Param(":product"), 10, 64)
	if err != nil {
		return errors.New("参数不足")
	}
	product, err := services.ProductInfo(pid)
	if err != nil {
		return errors.New("货品不存在")
	}
	c.LoadTop = true
	c.Data["product"] = product
	c.Data["title"] = product.Name
	c.Data["cat"], c.Data["catPath"] = catPath(product.Goods.Cat)
	c.TplName = "product/product.html,product/product_basic.html"
	return nil
}

func catPath(c *models.GoodsCat) (*models.GoodsCat, string) {
	var cat *models.GoodsCat
	if c != nil {
		res := `<span><a href="/gallery-` + string(c.ID) + `.html">` + string(c.Name) + `</a></span>
            <span>&gt;</span> `
		if c.Chind != nil {
			ct, r := catPath(c.Chind)
			res += r
			cat = ct
		} else {
			cat = c
		}
		return cat, res
	}
	return nil, ""
}

//GetGoodsLink 商品先关连接
func (c *ProductController) GetGoodsLink() {
	goodsID, _ := c.GetInt64("goods", 0)
	catID, _ := c.GetInt16("cat", 0)
	c.Data["data"] = services.GoodsLink(uint64(goodsID), int(catID))
	c.TplName = "product/goodsLink.html"
}

//GoodsDiscuss 商品评论
func (c *ProductController) GoodsDiscuss() {
	goodsID, _ := c.GetInt64("goods", 0)
	page, _ := c.GetInt32("page", 1)
	pagelength, _ := c.GetInt32("length", 20)
	c.Data["data"] = services.Discuss(uint64(goodsID), page, pagelength)
	beego.Error(c.Data)
	c.TplName = "product/DiscussInit.html"
}

//SellLogList 销售记录
func (c *ProductController) SellLogList() {
	productID, _ := c.GetInt64(":product", 0)
	page, _ := c.GetInt64(":page", 1)
	isGoods := c.GetString("tp", "goods")
	var pageRes *models.PageObject
	if isGoods == "product" {
		pageRes = services.ProductSellLogs(uint64(productID), page, 50)
	} else {
		pageRes = services.GoodsSellLogs(uint64(productID), page, 50)
	}
	c.Data["data"] = pageRes.Data
	c.Data["total"] = pageRes.Row
	c.TplName = "product/SellLogs.html"
}
