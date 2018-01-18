package services

import (
	"github.com/sparon/goshop/models"
)

//ProductInfo 获取商品信息
//@Param ProductID uint64 货品ID
func ProductInfo(ProductID uint64) (*models.Product, error) {
	return models.ProductBean.ProductInfo(ProductID)
}

//ProductInfos 获取多个产品的商品信息
//@Param ProductIDS []uint64 商品ID集合
func ProductInfos(ProductIDS []uint64) ([]*models.Product, error) {
	return models.ProductBean.ProductInfos(ProductIDS)
}

//GoodsLink 先关商品连接
//@Param goodsID uint64 商品ID
//@Param catID int 商品栏目ID
func GoodsLink(goodsID uint64, catID ...int) []map[string]interface{} {
	if len(catID) > 0 {
		return models.ProductBean.GoodsLink(goodsID, catID[0])
	}
	return models.ProductBean.GoodsLink(goodsID, 0)
}

//Discuss 评论
//@Param goodsID uint64 商品ID
//@Param args ...int32 0=>当前页,1=>每页条数
func Discuss(goodsID uint64, args ...int32) *models.PageObject {
	page := int32(1)
	pagelength := int32(20)
	if len(args) > 1 {
		pagelength = args[1]
		page = args[0]
	} else if len(args) > 0 {
		page = args[0]
	}
	return models.ProductBean.Discuss(goodsID, page, pagelength)
}

//GoodsSellLogs 商品销售记录
//@Param goodsID 商品ID
//@Param args ...int64 0=>当前页,1=>每页条数
func GoodsSellLogs(goodsID uint64, args ...int64) *models.PageObject {
	page := int64(1)
	pagelength := int64(20)
	if len(args) > 1 {
		pagelength = args[1]
		page = args[0]
	} else if len(args) > 0 {
		page = args[0]
	}
	p := &models.PageObject{
		NowPage:   page,
		PageCount: pagelength,
	}
	p, err := models.ProductBean.GoodsSellLogs(goodsID, p)
	if err != nil { //查询出错时返回空数据
		p = &models.PageObject{
			NowPage:   page,
			PageCount: pagelength,
		}
		p.Data = make([]map[string]interface{}, 0)
	}
	return p
}

//ProductSellLogs 货品销售记录
//@Param ProductID 货品ID
//@Param args ...int32 0=>当前页,1=>每页条数
func ProductSellLogs(ProductID uint64, args ...int64) *models.PageObject {
	page := int64(1)
	pagelength := int64(20)
	if len(args) > 1 {
		pagelength = args[1]
		page = args[0]
	} else if len(args) > 0 {
		page = args[0]
	}
	p := &models.PageObject{
		NowPage:   page,
		PageCount: pagelength,
	}
	p, err := models.ProductBean.ProductSellLogs(ProductID, p)
	if err != nil { //查询出错时返回空数据
		p = &models.PageObject{
			NowPage:   page,
			PageCount: pagelength,
		}
		p.Data = make([]map[string]interface{}, 0)
	}
	return p
}
