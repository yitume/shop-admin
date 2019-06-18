package goods

import (
	"time"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/router/mdw/auth"
	"git.yitum.com/saas/shop-admin/service"
)

func CategoryList(c *gin.Context) {
	// TODO auth.Default(c).Id open_id
	total, list := service.GoodsCategory.ListPage(c, mysql.Conds{}, trans.ReqPage{Sort: "update_time desc"})
	api.JSONList(c, list, total)
}

func CategoryCreate(c *gin.Context) {
	req := trans.ReqGoodscategoryCreate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	if err := service.GoodsCategory.Create(c, model.Db, &mysql.GoodsCategory{
		0, req.Name, req.Pid, req.Icon, 0, time.Now().Unix(), time.Now().Unix(), "", "", "", "", 0, "", 0, auth.Default(c).Id}); err != nil {
		api.JSON(c, api.MsgErr, "创建分类失败")
		return
	}
	api.JSONOK(c)
}

func CategoryUpdate(c *gin.Context) {
	req := trans.ReqGoodscategoryUpdate{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}

	// TODO auth.Default(c).Id open_id
	if err := service.GoodsCategory.Update(c, model.Db, req.Id, mysql.Ups{
		"name": req.Name,
		"pid":  req.Pid,
		"icon": req.Icon,
	}); err != nil {
		api.JSON(c, api.MsgErr, "更新分类失败")
		return
	}
	api.JSONOK(c)
}

func CategoryInfo(c *gin.Context) {
	req := trans.ReqGoodscategoryInfo{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	// TODO 完善ReqGoodscategoryInfo字段
	res, err := service.GoodsCategory.Info(c, req.Id)
	if err != nil {
		api.JSON(c, api.MsgErr, "获取分类失败")
		return
	}
	api.JSONOK(c, gin.H{
		"info": res,
	})
}

func CategoryDel(c *gin.Context) {
	req := trans.ReqGoodscategoryInfo{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request app list params is error")
		return
	}
	if err := service.GoodsCategory.Delete(c, model.Db, req.Id); err != nil {
		api.JSON(c, api.MsgErr, "删除分类失败")
		return
	}
	api.JSONOK(c)
}
