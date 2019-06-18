package goods

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/mdw/auth"

	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

func SpecvalueCreate(c *gin.Context) {
	req := trans.ReqGoodsspecvalueCreate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.GoodsSpecValue.Create(c, model.Db, &mysql.GoodsSpecValue{
		0, req.SpecId, req.Name, 0, "", "", 0, auth.Default(c).Id,
	}); err != nil {
		api.JSON(c, api.MsgErr, "create error")
		return
	}
	api.JSONOK(c)

}
