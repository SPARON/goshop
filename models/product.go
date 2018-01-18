package models

import (
	"reflect"

	"github.com/hprose/hprose-go"
)

//ProductBean 商品数据处理对象
var ProductBean *ProductModel

//ProductModel 商品模型
type ProductModel struct {
	//获取商品信息
	//@Param ProductID uint64 商品ID
	ProductInfo func(ProductID uint64) (*Product, error)
	//获取多个产品的商品信息
	//@Param ProductIDS []uint64 商品ID集合
	ProductInfos func(ProductIDS []uint64) ([]*Product, error) `name:"ProductsInfo"`
	//先关商品连接
	//@Param goodsID uint64 商品ID
	//@Param catID int 商品栏目ID
	GoodsLink func(goodsID uint64, catID int) []map[string]interface{} `name:"goodsLink"`
	//评论
	//@Param goodsID uint64 商品ID
	Discuss func(goodsID uint64, page int32, pagelength int32) *PageObject `name:"discuss"`
	//商品销售记录
	//@Param goodsID 商品ID
	//@Param page PageObject 分页信息
	GoodsSellLogs func(goodsID uint64, page *PageObject) (*PageObject, error) `name:"goodsSellLogs"`
	//货品销售记录
	//@Param ProductID 货品ID
	//@Param page PageObject 分页信息
	ProductSellLogs func(ProductID uint64, page *PageObject) (*PageObject, error) `name:"goodsSellLogs"`
}

func initProduct() {
	hproseClient.UseService(&ProductBean)
	hprose.ClassManager.Register(reflect.TypeOf(Product{}), "Product", "v")
	hprose.ClassManager.Register(reflect.TypeOf(ProductSpec{}), "ProductSpec", "v")
	hprose.ClassManager.Register(reflect.TypeOf(GoodsSpec{}), "GoodsSpec", "v")
	hprose.ClassManager.Register(reflect.TypeOf(GSpec{}), "GSpec", "v")
	hprose.ClassManager.Register(reflect.TypeOf(Brand{}), "Brand", "v")
	hprose.ClassManager.Register(reflect.TypeOf(GoodsCat{}), "GoodsCat", "v")
	hprose.ClassManager.Register(reflect.TypeOf(GoodsImage{}), "GoodsImage", "v")
	hprose.ClassManager.Register(reflect.TypeOf(GoodsProp{}), "GoodsProp", "v")
	hprose.ClassManager.Register(reflect.TypeOf(Goods{}), "Goods", "v")
	hprose.ClassManager.Register(reflect.TypeOf(PageObject{}), "PageObject", "v")
}

//Product 货品对象
type Product struct {
	ProductID  uint64  `v:"product_id"`  //货品ID
	GoodsID    uint64  `v:"goods_id"`    //商品ID
	Name       string  `v:"name"`        //商品名称
	SpecInfo   string  `v:"spec_info"`   //货品型号
	Bn         string  `v:"bn"`          //货品编号
	Store      int64   `v:"store"`       //库存数量
	Freez      int64   `v:"freez"`       //锁定库存数量
	Price      float64 `v:"price"`       //单价
	MktPrice   float64 `v:"mktprice"`    //市场价
	Marktable  string  `v:"marktable"`   //是否上架("true","false")
	Weight     float32 `v:"weight"`      //单体重量
	CatID      int64   `v:"cat_id"`      //分类ID
	TypeID     int64   `v:"type_id"`     //类型ID
	Image      string  `v:"image"`       //图片编号
	SImage     string  `v:"s_image"`     //小图地址
	MImage     string  `v:"m_image"`     //中图地址
	LImage     string  `v:"l_image"`     //大图地址
	ImageURL   string  `v:"image_url"`   //原始图地址
	SupplierID uint64  `v:"supplier_id"` //供应商ID
	Goods      *Goods  `v:"goods"`       //商品
}

//GSpec 商品规格
type GSpec struct {
	Specs []*GoodsSpec `v:"specs"`
}

//ProductSpec 商品规格
type ProductSpec struct {
	PrivateSpecValueID int64  `v:"private_spec_value_id"`
	SpecValue          string `v:"spec_value"`
	SpecValueID        int64  `v:"spec_value_id"`
	SpecGoodsImage     string `v:"spec_goods_image"`
	ProductID          int64  `v:"product"`
}

//GoodsSpec 商品规格信息
type GoodsSpec struct {
	Spec     []map[string]interface{} `v:"spec"`
	GoodSpec int64                    `v:"goodspec"`
	Name     string                   `v:"name"`
}

//Brand 品牌信息
type Brand struct {
	ID   int32  `v:"brand_id"`   //ID
	Name string `v:"brand_name"` //名称
	URL  string `v:"brand_url"`  //地址
	Desc string `v:"brand_desc"` //介绍
	Logo string `v:"brand_logo"` //图标
}

//GoodsCat 商品分类
type GoodsCat struct {
	ID    int32     `v:"cat_id"`   //ID
	Path  string    `v:"cat_path"` //路径
	Name  string    `v:"cat_name"` //名称
	Chind *GoodsCat `v:"child"`    //子类
}

//GoodsImage 商品图片
type GoodsImage struct {
	ID     string `v:"image_id"` //图片ID
	LImg   string `v:"l_image"`  //大图
	MImg   string `v:"m_image"`  //中图
	SImg   string `v:"s_image"`  //小图
	URL    string `v:"url"`      //原图
	Weight int    `v:"weight"`   //宽
	Height int    `v:"height"`   //高
}

//GoodsProp 商品属性
type GoodsProp struct {
	GoodsP    int    `v:"goodsP"`     //序号
	PropID    int    `v:"props_id"`   //属性
	TypeID    int    `v:"type_id"`    //类型
	Name      string `v:"name"`       //名称
	Alias     string `v:"alias"`      //别名
	PropValue string `v:"prop_value"` //属性值
}

//Goods 商品
type Goods struct {
	GoodsID       uint64                 `v:"goods_id"`         //商品ID
	Bn            string                 `v:"bn"`               //商品编号
	Name          string                 `v:"name"`             //商品名称
	Price         float64                `v:"price"`            //销售价
	TypeID        int                    `v:"type_id"`          //类型
	Cat           *GoodsCat              `v:"cat"`              //分类
	Brand         *Brand                 `v:"brand"`            //品牌
	Cost          float64                `v:"cost"`             //成本价
	MktPrice      float64                `v:"mktprice"`         //市场价
	Weight        float64                `v:"weight"`           //重量
	Unit          string                 `v:"unit"`             //单位
	Brief         string                 `v:"brief"`            //商品简介
	ImgDefID      string                 `v:"image_default_id"` //默认图片
	ThumbnailPic  string                 `v:"thumbnail_pic"`    //缩略图
	SmallPic      string                 `v:"small_pic"`        //小图
	BigPic        string                 `v:"big_pic"`          //大图
	Intro         string                 `v:"intro"`            //详细介绍
	SpecDesc      string                 `v:"spec_desc"`        //货品规格序列化
	Params        string                 `v:"params"`           //商品规格序列换
	CommentsCount int                    `v:"comments_count"`   //评论次数
	ViewWCount    int32                  `v:"view_w_count"`     //周浏览次数
	ViewCount     int32                  `v:"view_count"`       //浏览次数
	CountStat     string                 `v:"count_stat"`       //统计数据序列化
	BuyCount      int                    `v:"buy_count"`        //购买次数
	BuyWCount     int                    `v:"buy_w_count"`      //周购买次数
	Spec          *GSpec                 `v:"spec"`             //货品规格信息
	Prop          []*GoodsProp           `v:"prop"`             //参数信息
	Images        map[string]*GoodsImage `v:"images"`           //图片数据
}

//PageObject 分页信息
type PageObject struct {
	Row       int64                    `v:"rows"`
	Pages     int64                    `v:"pages"`
	NowPage   int64                    `v:"nowpage"`
	PageCount int64                    `v:"pagecount"`
	Data      []map[string]interface{} `v:"data"`
}
