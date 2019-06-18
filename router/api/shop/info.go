package shop

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/service"
)

// Add Add a menu
func Info(c *gin.Context) {
	resp := mysql.Shop{}
	model.Db.Find(&resp)

	api.JSONOK(c, gin.H{
		"info": resp,
	})
}

func Update(c *gin.Context) {
	req := trans.ReqShopUpdate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.Shop.Update(c, model.Db, 1, mysql.Ups{
		"logo":           req.Logo,
		"name":           req.Name,
		"contact_number": req.ContactNumber,
		"description":    req.Description,
	}); err != nil {
		api.JSON(c, api.MsgErr, "update is error")
		return
	}
	api.JSONOK(c)
}

func SetGoodsCategoryStyle(c *gin.Context) {
	req := trans.ReqShopList{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.Shop.Update(c, model.Db, 1, mysql.Ups{
		"goods_category_style": req.GoodsCategoryStyle,
	}); err != nil {
		api.JSON(c, api.MsgErr, "update is error")
		return
	}
	api.JSONOK(c)
}

func SetOrderExpires(c *gin.Context) {
	req := trans.ReqShopUpdate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	if err := service.Shop.Update(c, model.Db, 1, mysql.Ups{
		"order_auto_close_expires":         req.OrderAutoCloseExpires,
		"order_auto_confirm_expires":       req.OrderAutoConfirmExpires,
		"order_auto_close_refound_expires": req.OrderAutoCloseRefoundExpires,
	}); err != nil {
		api.JSON(c, api.MsgErr, "update is error")
		return
	}
	api.JSONOK(c)
}
