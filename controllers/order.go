package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"

	"gitee.com/sparon/goshop/models"
	"gitee.com/sparon/goshop/services"
	"gitee.com/sparon/goshop/utils"
	"gitee.com/sparon/goshop/utils/urlpath"
)

//OrderController 订单控制器
type OrderController struct {
	LcfglyController
}

//errorPage 显示错误页面
var errorPage = func(c *LcfglyController, resErr *error) {
	if (*resErr) != nil {
		c.ErrorPage((*resErr).Error())
	}
}

//errorCreatePayOrder 支付单据创建失败错误
var errorCreatePayOrder = errors.New("支付单据创建失败")

//errorNoWaitPayOrder 无法获取支付订单信息错误
var errorNoWaitPayOrder = errors.New("无法获取待支付订单信息")

//GetCart 购物车列表
func (c *OrderController) GetCart() {
	if c.Ctx.GetCookie(utils.Conf.CookieNoPostOrder) != "" { //如果存在未提交的订单信息,则直接提交订单,并且用户已经登入
		c.Redirect(urlpath.Path("order"), 302) //跳转到订单页面
		return
	}
	//显示购物车
	c.Data["title"] = "购物车"
	//获取购物车中的产品
	memberID := services.IsLogin(c.Ctx)
	if memberID == 0 { //如果用户为登入,,获取当前未登入的购物车数据
		memberID, _ = strconv.ParseUint(c.Ctx.GetCookie(utils.Conf.CookieCartMember), 10, 64)
	}
	carts := services.CheckOutCart(memberID)
	for k, v := range carts {
		c.Data[k] = v
	}
	c.LoadTop = true
	c.TplName = "orders/cart.html"
}

//GetMinCartView 迷你购物车
func (c *OrderController) GetMinCartView() {
	memberID := services.IsLogin(c.Ctx)
	if memberID == 0 { //如果用户为登入,,获取当前未登入的购物车数据
		memberID, _ = strconv.ParseUint(c.Ctx.GetCookie(utils.Conf.CookieCartMember), 10, 64)
	}
	carts := services.CheckOutCart(memberID)
	for k, v := range carts {
		c.Data[k] = v
	}
	c.TplName = "orders/cartView.html"
}

//AddCart 添加商品到购物车
func (c *OrderController) AddCart() {
	productID, _ := c.GetInt64("goods[product_id]", 0)
	nums, _ := c.GetInt64("goods[num]", 0)
	if productID == 0 || nums == 0 {
		c.Data["error"] = true
	}
	memberID := services.IsLogin(c.Ctx)
	if memberID == 0 {
		memberID, _ = strconv.ParseUint(c.Ctx.GetCookie(utils.Conf.CookieCartMember), 10, 64) //未登入购物车编号
	}
	cartList, err := services.AddCart(memberID, uint64(productID), nums)
	if err != nil {
		c.Data["error"] = true
	}
	for k, v := range cartList {
		c.Data[k] = v
	}
	c.TplName = "orders/addCart.html"
}

//PostOrder 生成订单信息
func (c *OrderController) PostOrder() {
	memberID := services.IsLogin(c.Ctx, false)
	if memberID == 0 { //为登入跳转到登入信息
		c.Ctx.SetCookie(utils.Conf.CookieNoPostOrder, c.Ctx.Request.Form.Encode(), utils.Conf.CookieExpireTime, "/", utils.Conf.CookieDomain)
		c.JSONFail("请先登录", urlpath.Host("sso")+urlpath.Path("login")+"?redirectURL="+urlpath.Host("order")+urlpath.Path("order")) //未登录跳转到登录页面
		return
	}
	orderInfo, err := services.ExecOrder(c.Ctx.Request.Form, memberID)
	if err != nil {
		c.JSONFail("订单创建失败:" + err.Error())
	} else {
		c.JSONSuccess(nil, "订单检出成功", urlpath.Path("order")+"?orderID="+orderInfo.Ident)
	}
}

//GetOrder 订单信息
func (c *OrderController) GetOrder() {
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 { //为登入跳转到登入信息
		return
	}
	orderIden := c.GetString("orderID")
	if orderIden == "" { //如果不存在传入的订单
		cookieNoPostOrder := c.Ctx.GetCookie(utils.Conf.CookieNoPostOrder)
		if cookieNoPostOrder != "" && memberID != 0 { //如果存在未提交的订单信息,则直接提交订单,并且用户已经登入
			c.Ctx.SetCookie(utils.Conf.CookieNoPostOrder, "", -1, "/", utils.Conf.CookieDomain) //移除未提交的订单信息
			urlValues, err := url.ParseQuery(cookieNoPostOrder)                                 //解析未提交的参数信息
			if err == nil {                                                                     //未提交的订单信息解析成功
				orderInfo, err := services.ExecOrder(urlValues, memberID) //提交商品信息
				if err == nil {                                           //订单处理成功的
					orderIden = orderInfo.Ident
				}
			}
		}
	}
	if orderIden == "" { //不存在订单信息,跳转到购物车页面
		c.ErrorPage("订单信息不存在", urlpath.Path("cart"))
		return
	}
	orderInfo, err := services.CartOrder(orderIden)
	if err != nil || orderInfo == nil || orderInfo.MemberID != memberID { //订单数据不存在
		c.ErrorPage("订单信息不存在", urlpath.Path("cart"))
		return
	}
	c.Data["order"] = orderInfo
	c.Data["title"] = "订单确认"
	c.LoadTop = true
	c.TplName = "orders/checkOrder.html" //订单提交页面
}

//PostCreateOrder 提交订单
func (c *OrderController) PostCreateOrder() {
	ident := c.GetString("ident")
	addressID := c.GetString("address")
	source := c.GetString("source", "pc")
	memberID, err := c.GetInt64("member_id", 0)
	if ident == "" || addressID == "" || err != nil || memberID == 0 {
		beego.Debug("订单信息不足:ident=", ident, " address=", addressID)
		c.JSONFail("订单信息不足")
		return
	}
	memo := "memo"
	shipping := "shipping"
	args := make(map[string]interface{})
	for k, v := range c.Ctx.Request.Form {
		if strings.HasPrefix(k, memo) || strings.HasPrefix(k, shipping) {
			startPrefix := memo
			if strings.HasPrefix(k, shipping) {
				startPrefix = shipping
			}
			kres := strings.SplitAfterN(k, "[", 2)
			if len(kres) > 1 {
				key := beego.Substr(kres[1], 0, len(kres)-1)
				if vv, ok := args[key]; ok {
					vv.(map[string]interface{})[startPrefix] = v[0]
				} else {
					args[key] = map[string]interface{}{startPrefix: v[0]}
				}
			}
		}
	}
	if len(args) < 1 {
		c.JSONFail("订单信息不足")
		return
	}
	args["address_id"] = addressID
	args["ip"] = c.Ctx.Input.IP()
	args["source"] = source
	result, err := services.CreateOrder(args, ident, uint64(memberID))
	if err != nil {
		beego.Info("订单创建失败:", err)
		c.JSONFail("订单创建失败")
		return
	}
	res := bytes.NewBufferString("")
	for _, s := range result.OrderIDS {
		res.WriteString(strconv.FormatUint(s, 10))
		res.WriteString(",")
	}
	if len(result.OrderIDS) < 1 {
		c.JSONFail("订单创建失败")
		return
	}
	res.Truncate(res.Len() - 1) //去除最后一个','
	crpyt, _ := utils.AesEncrypt(res.Bytes())
	res.Truncate(0)
	res.WriteString(urlpath.Path("createPayBill"))
	res.WriteString("?s=")
	res.Write(crpyt)
	c.JSONSuccess(result, "订单创建成功", res.String()) //订单生成成功
}

//CratePayBill 生成支付单据
func (c *OrderController) CratePayBill() (resErr error) {
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 { //未登入
		return
	}
	defer errorPage(&c.LcfglyController, &resErr)
	query := c.GetString("s")
	if query == "" {
		return errorCreatePayOrder

	}
	queryBytes, err := utils.AesDecrypt([]byte(query))
	if err != nil {
		return err
	}
	orderIDS := strings.Split(string(queryBytes), ",")
	ods := make([]interface{}, len(orderIDS))
	for index, v := range orderIDS {
		ods[index] = v
	}
	result, err := services.CreatePayBill(memberID, ods...)
	if err != nil {
		return err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return errorCreatePayOrder
	}
	b, err = utils.AesEncrypt(b)
	if err != nil {
		return errorCreatePayOrder
	}
	c.EnableRender = false
	c.Redirect(urlpath.Host("order")+urlpath.Path("waitpay")+"?s="+string(b), 302) //跳转到待支付页面
	return nil
}

//GetWaitPay 待支付页面
func (c *OrderController) GetWaitPay() (resErr error) {
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 {
		return
	}
	defer errorPage(&c.LcfglyController, &resErr)
	query := c.GetString("s")
	if query == "" {
		return errorNoWaitPayOrder
	}
	queryBytes, err := utils.AesDecrypt([]byte(query))
	if err != nil {
		return errorNoWaitPayOrder
	}
	err = json.Unmarshal(queryBytes, &c.Data)
	if err != nil {
		return errorNoWaitPayOrder
	}
	c.Data["memberPoint"] = services.GetMemberPoint(memberID) //会员可用积分
	c.Data["MaxPoint"] = c.Data["amount"].(float64) * 0.05
	c.Data["title"] = "单据结算收银台"
	c.LoadTop = true
	c.TplName = "orders/waitPay.html"
	return nil
}

//PostPay 支付单据
func (c *OrderController) PostPay() (resErr error) {
	defer errorPage(&c.LcfglyController, &resErr)
	lostParam := errors.New("参数异常")
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 {
		return
	}
	ident, err := c.GetInt64("bill", 0)
	if err != nil {
		return lostParam
	}
	payStr := c.GetString("pay")
	if payStr == "" {
		return lostParam
	}
	point, _ := c.GetInt("point", 0) //获取用户提交的消费积分
	pay := models.StringToPaymentType(payStr)
	result, err := services.Pay(int(pay), c.Ctx.Input.IP(), uint64(ident), memberID, point)
	if err != nil {
		return err
	}
	c.EnableRender = false
	c.Ctx.ResponseWriter.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Ctx.ResponseWriter.ResponseWriter.Write([]byte(result)) //直接输出用户支付数据
	return nil
}

//GetPayResult 支付结果
func (c *OrderController) GetPayResult() {
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 {
		return
	}
	query, err := c.GetInt64("s")
	if err != nil {
		c.ErrorPage("无法获取支付单据信息")
		return
	}
	c.Data, err = services.GetPayBill(uint64(query))
	if err != nil {
		c.ErrorPage(err.Error())
		return
	}
	c.LoadTop = true
	c.TplName = "orders/payResult.html"
}
