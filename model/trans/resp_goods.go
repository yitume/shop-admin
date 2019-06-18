package trans

import (
	"git.yitum.com/saas/shop-admin/model/mysql"
)

type Goods struct {
	Id                 int                        `json:"id"`                   // 商品公共表id
	Title              string                     `json:"title"`                // 商品名称
	Images             mysql.GoodsImagesJson      `json:"images"`               // 商品图片 默认第一个为封面图片
	CategoryIds        mysql.GoodsCategoryIdsJson `json:"category_ids"`         // 商品分类
	BaseSaleNum        int                        `json:"base_sale_num"`        // 基础销量
	Body               mysql.GoodsBodyJson        `json:"body"`                 // 商品内容
	IsOnSale           int                        `json:"is_on_sale"`           // 是否需上架出售 0 否 1 是
	ImageSpecId        int                        `json:"image_spec_id"`        // 使用图片的规格id
	ImageSpecImages    []string                   `json:"image_spec_images"`    // 规格图片集合，废弃
	SkuList            []mysql.GoodsSku           `json:"sku_list"`             // sku商品集合，数组
	CreateTime         int64                      `json:"create_time"`          // 创建时间
	Price              float64                    `json:"price"`                // 商品价格
	UpdateTime         int64                      `json:"update_time"`          // 修改时间
	EvaluationGoodStar int                        `json:"evaluation_good_star"` // 好评星级
	EvaluationCount    int                        `json:"evaluation_count"`     // 评价数
	Stock              int                        `json:"stock"`                // goods表库存之和
	SaleNum            int                        `json:"sale_num"`             // 销售量
	GroupSaleNum       int                        `json:"group_sale_num"`       // 拼团销量
	SaleTime           int64                      `json:"sale_time"`            // 开售时间
	DeleteTime         int64                      `json:"delete_time"`          // 软删除时间
	SpecList           mysql.GoodsSpecListJson    `json:"spec_list"`            //
	Img                string                     `json:"img"`                  // 封面图
	PayType            int                        `json:"pay_type"`             // 计算方式：1 按件数 2 按重量
	FreightFee         float64                    `json:"freight_fee"`          // 运费
	FreightId          int                        `json:"freight_id"`           // 运费模板id
}
