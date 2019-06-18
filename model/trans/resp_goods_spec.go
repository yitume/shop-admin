package trans

import "git.yitum.com/saas/shop-admin/model/mysql"

type RespGoodsSpecList struct {
	Id     int                    `json:"id"`
	Name   string                 `json:"name"`
	Values []mysql.GoodsSpecValue `json:"values"`
}
