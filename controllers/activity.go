package controllers

import (
	"errors"

	"strconv"

	"gitee.com/sparon/goshop/services"
	"gitee.com/sparon/goshop/utils"
	"gitee.com/sparon/goshop/utils/urlpath"
)

//ActivityController 活动控制器
type ActivityController struct {
	LcfglyController
}

//Yuyue 预约
func (c *ActivityController) Yuyue() (resErr error) {
	defer errorPage(&c.LcfglyController, &resErr)
	productID, err := c.GetInt64("product_id", 0)
	if err != nil {
		resErr = errors.New("货品不存在")
		return
	}
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 { //用户未登录，跳转到登入页面
		return nil
	}
	activity, err := services.ActivityQiangou(int64(memberID), productID)
	c.Data["data"] = activity
	c.Data["Consignees"] = services.ConsigneesList(memberID) //会员收获地址列表
	c.TplName = "activity/yuyue.html"
	return nil
}

//Qiangou 抢购
func (c *ActivityController) Qiangou() (resErr error) {
	memberID := services.IsLogin(c.Ctx, true)
	if memberID == 0 { //用户未登录，跳转到登入页面
		return nil
	}
	productID, err := c.GetInt64("product_id", 0)
	if err != nil {
		resErr = errors.New("参数数据不足")
		return
	}
	ip := c.Ctx.Input.IP()
	orderID, err := services.RushBuy(int64(memberID), productID, ip, "pc")
	if err != nil {
		return err
	}
	order := strconv.FormatUint(uint64(orderID), 10)
	crpyt, _ := utils.AesEncrypt([]byte(order))
	c.JSONSuccess(order, "订单创建成功", urlpath.Path("createPayBill")+"?s="+string(crpyt)) //创建支付单据
	return nil
}
