package area

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := trans.ReqAreaList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	var conds = make(mysql.Conds)
	if req.Pid != 0 {
		conds["pid"] = req.Pid
	}
	switch req.Level {
	case 1:
		conds["level"] = []int{1, 2}
	case 2:
		conds["level"] = []int{1, 2, 3}
	default:
		conds["level"] = []int{1}
	}
	req.ReqPage.PageSize = 10000
	total, list := service.Area.ListPage(c, conds, req.ReqPage)
	api.JSONList(c, list, total)
}
