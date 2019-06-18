package coupon

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func List(c *gin.Context) {
	req := trans.ReqCouponList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	// total, list := service.Coupon.ListPage(c, req.ReqPage, auth.Default(c).Id)
	total, list := service.Coupon.ListPage(c, mysql.Conds{}, req.ReqPage)
	api.JSONList(c, list, total)
}
