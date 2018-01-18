package controllers

import (
	"strings"

	"gitee.com/sparon/goshop/models"
	"gitee.com/sparon/goshop/services"
)

//MemberController 会员控制器
type MemberController struct {
	LcfglyController
}

//GetAddAddress 添加用户收获地址
func (c *MemberController) GetAddAddress() {
	c.Data["cons"] = &models.ConsigneeInfo{}
	c.TplName = "members/add_address.html"
}

//PostAddAddress 添加用户收获地址
func (c *MemberController) PostAddAddress() {
	login := services.IsLogin(c.Ctx)
	if login == 0 { //用户为登入
		c.JSONFail("请先登入")
		return
	}
	addr := c.GetString("addr")
	area := c.GetString("area")
	mobile := c.GetString("mobile")
	name := c.GetString("name")
	addressID, _ := c.GetInt64("address_id", 0)
	tel := c.GetString("tel")
	zip := c.GetString("zip")
	var cons *models.ConsigneeInfo
	var err error
	if addressID > 0 {
		cons, err = services.Consignee(login, uint64(addressID))
		if err != nil {
			cons = &models.ConsigneeInfo{}
			addressID = 0
		}
	} else {
		cons = &models.ConsigneeInfo{}
	}
	cons.Address = addr
	cons.Area = area
	cons.Phone = mobile
	cons.MemberID = login
	cons.Name = name
	cons.Tel = tel
	cons.Zip = zip
	if addressID > 0 {
		result, err := services.UpdateConsignee(cons)
		if err != nil || !result { //错误
			c.JSONFail(err.Error())
			return
		}
	} else {
		cons.ID, err = services.AddConsignee(cons)
		if err != nil { //错误
			c.JSONFail("数据保存失败")
			return
		}
	}
	c.Data["address_id"] = cons.ID
	c.Data["cons"] = cons
	c.Data["result"] = true
	c.Data["area_id"] = strings.Split(cons.Area, ":")[2]
	c.Data["area_name"] = strings.Split(cons.Area, ":")[1]
	c.Data["area_name_array"] = strings.Split(strings.Split(cons.Area, ":")[1], "/")
	c.TplName = "members/add_address.html"
}

//GetEditAddress 编辑用户地址
func (c *MemberController) GetEditAddress() {
	addressID, _ := c.GetInt64("address_id", 0)
	login := services.IsLogin(c.Ctx)
	if login == 0 { //用户已经登入,跳转到来源页面
		c.JSONFail("请先登入")
		return
	}
	var cons *models.ConsigneeInfo
	var err error
	if addressID == 0 {
		cons, err = services.DefaultConsignee(login)
	} else {
		cons, err = services.Consignee(login, uint64(addressID))
	}
	if err == nil {
		c.Data["cons"] = cons
	} else {
		c.Data["cons"] = &models.ConsigneeInfo{}
	}
	c.TplName = "members/add_address.html"
}
