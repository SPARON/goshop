package models

import (
	"fmt"

	"github.com/hprose/hprose-go"
)

var hproseClient *hprose.TcpClient

func init() {
	hproseClient = hprose.NewTcpClient("tcp://localhost:4321")
	//hproseClient.AddFilter(&filter{})
	hproseClient.DebugEnabled = true
	initLoginModel() //登入模型
	initMember()     //会员模型
	initProduct()    //商品模型
	initOrder()      //订单模型
	initActivity()   //活动模型
}

type filter struct {
}

func (c *filter) InputFilter(data []byte, context hprose.Context) []byte {
	fmt.Println("接受到数据：", len(data), "个字节")
	fmt.Println(string(data[len(data)-2:]))
	return data
}

func (c *filter) OutputFilter(data []byte, context hprose.Context) []byte {
	return data
}
