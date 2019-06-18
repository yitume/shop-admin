package trans

import (
	"git.yitum.com/saas/shop-admin/model/mysql"
)

type ReqGoodsCreateOrUpdate struct {
	ReqGoodsCreate
	Skus        []mysql.GoodsSku `json:"skus"`
	CategoryIds []string         `json:"category_ids" form:"category_ids[]" ` // 商品分类
}

type ReqGoodsOnSale struct {
	Ids []int `json:"ids"`
}

type ReqGoodscategoryCreate struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
	Pid  int    `json:"pid"`
}

type ReqGoodscategoryInfo struct {
	Id int `form:"id"`
}

type ReqGoodscategoryUpdate struct {
	Id   int    `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
	Pid  int    `json:"pid"`
}

type ReqGoodsspecCreate struct {
	Name string `json:"name"`
}

type ReqGoodsspecvalueCreate struct {
	SpecId int    `json:"spec_id"`
	Name   string `json:"name"`
}
