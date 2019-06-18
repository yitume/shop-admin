package freight

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	// TODO auth.Default(c).Id open_id
	total, list := service.Freight.ListPage(c, mysql.Conds{}, trans.ReqPage{Sort: "update_time desc"})
	api.JSONList(c, list, total)
}
